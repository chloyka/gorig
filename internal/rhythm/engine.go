package rhythm

import (
	"sync"
	"time"

	"github.com/chloyka/gorig/internal/onset"
)

type Engine struct {
	mu sync.RWMutex

	tempo *TempoState

	totalSamples int64
	currentSlot  int64

	onsetEvents <-chan onset.Event

	pendingOnset  *onset.Event
	lastSlotFired int64

	quantizedChan chan QuantizedOnset

	onStateChange func(bpm float64, subdivision int)
}

type EngineConfig struct {
	SampleRate    float32
	InitialBPM    float64
	Subdivision   Subdivision
	OnsetEvents   <-chan onset.Event
	OnStateChange func(bpm float64, subdivision int)
}

func NewEngine(cfg EngineConfig) *Engine {

	bpm := cfg.InitialBPM
	if bpm < MinBPM || bpm > MaxBPM {
		bpm = DefaultBPM
	}

	sub := cfg.Subdivision
	if !sub.IsValid() {
		sub = Sub8
	}

	e := &Engine{
		tempo:         NewTempoState(bpm, sub, cfg.SampleRate),
		onsetEvents:   cfg.OnsetEvents,
		quantizedChan: make(chan QuantizedOnset, 16),
		onStateChange: cfg.OnStateChange,
		lastSlotFired: -1,
	}

	return e
}

func (e *Engine) ProcessBuffer(bufferSize int) *QuantizedOnset {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.drainOnsetEvents()

	oldSlot := e.currentSlot
	e.totalSamples += int64(bufferSize)

	if e.tempo.SamplesPerSlot > 0 {
		e.currentSlot = e.totalSamples / e.tempo.SamplesPerSlot
	}

	if e.currentSlot > oldSlot && e.pendingOnset != nil && e.currentSlot > e.lastSlotFired {
		result := &QuantizedOnset{
			OriginalEvent: *e.pendingOnset,
			SlotIndex:     int(e.currentSlot % int64(e.tempo.Subdivision)),
			BeatPosition:  e.getBeatPositionLocked(),
			WasQueued:     true,
		}

		select {
		case e.quantizedChan <- *result:
		default:
		}

		e.pendingOnset = nil
		e.lastSlotFired = e.currentSlot
		return result
	}

	return nil
}

func (e *Engine) drainOnsetEvents() {
	if e.onsetEvents == nil {
		return
	}

	if e.pendingOnset != nil {
		for {
			select {
			case _, ok := <-e.onsetEvents:
				if !ok {
					return
				}

			default:
				return
			}
		}
	}

	select {
	case event, ok := <-e.onsetEvents:
		if ok {
			e.pendingOnset = &event
		}
	default:
		return
	}

	for {
		select {
		case _, ok := <-e.onsetEvents:
			if !ok {
				return
			}
		default:
			return
		}
	}
}

func (e *Engine) GetBeatPhase() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.getBeatPositionLocked()
}

func (e *Engine) getBeatPositionLocked() float64 {
	if e.tempo.SamplesPerBeat == 0 {
		return 0
	}
	samplesIntoBeat := e.totalSamples % e.tempo.SamplesPerBeat
	return float64(samplesIntoBeat) / float64(e.tempo.SamplesPerBeat)
}

func (e *Engine) GetBeatCount() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.tempo.SamplesPerBeat == 0 {
		return 0
	}
	return e.totalSamples / e.tempo.SamplesPerBeat
}

func (e *Engine) GetCurrentSlot() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentSlot
}

func (e *Engine) GetSlotInBeat() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return int(e.currentSlot % int64(e.tempo.Subdivision))
}

func (e *Engine) QuantizedOnsets() <-chan QuantizedOnset {
	return e.quantizedChan
}

func (e *Engine) RegisterTap(now time.Time) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	updated := e.tempo.RegisterTap(now)
	if updated && e.onStateChange != nil {
		e.onStateChange(e.tempo.BPM, int(e.tempo.Subdivision))
	}
	return updated
}

func (e *Engine) SetBPM(bpm float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.tempo.SetBPM(bpm)
	if e.onStateChange != nil {
		e.onStateChange(e.tempo.BPM, int(e.tempo.Subdivision))
	}
}

func (e *Engine) AdjustBPM(delta float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.tempo.AdjustBPM(delta)
	if e.onStateChange != nil {
		e.onStateChange(e.tempo.BPM, int(e.tempo.Subdivision))
	}
}

func (e *Engine) GetBPM() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.tempo.BPM
}

func (e *Engine) SetSubdivision(sub Subdivision) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.tempo.SetSubdivision(sub)
	if e.onStateChange != nil {
		e.onStateChange(e.tempo.BPM, int(e.tempo.Subdivision))
	}
}

func (e *Engine) NextSubdivision() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.tempo.NextSubdivision()
	if e.onStateChange != nil {
		e.onStateChange(e.tempo.BPM, int(e.tempo.Subdivision))
	}
}

func (e *Engine) PrevSubdivision() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.tempo.PrevSubdivision()
	if e.onStateChange != nil {
		e.onStateChange(e.tempo.BPM, int(e.tempo.Subdivision))
	}
}

func (e *Engine) GetSubdivision() Subdivision {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.tempo.Subdivision
}

func (e *Engine) GetSamplesPerSlot() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.tempo.SamplesPerSlot
}

func (e *Engine) GetSamplesPerBeat() int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.tempo.SamplesPerBeat
}

func (e *Engine) Close() {
	close(e.quantizedChan)
}
