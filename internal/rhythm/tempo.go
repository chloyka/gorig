package rhythm

import (
	"math"
	"time"
)

const (
	MaxTapHistory = 4

	TapResetThreshold = 2 * time.Second

	MinBPM = 30.0

	MaxBPM = 300.0

	DefaultBPM = 120.0
)

type TempoState struct {
	BPM         float64
	Subdivision Subdivision
	SampleRate  float32

	SamplesPerBeat int64
	SamplesPerSlot int64

	tapTimes []time.Time
}

func NewTempoState(bpm float64, subdivision Subdivision, sampleRate float32) *TempoState {
	t := &TempoState{
		BPM:         clampBPM(bpm),
		Subdivision: subdivision,
		SampleRate:  sampleRate,
		tapTimes:    make([]time.Time, 0, MaxTapHistory),
	}
	t.recalculateDerivedValues()
	return t
}

func (t *TempoState) SetBPM(bpm float64) {
	t.BPM = clampBPM(bpm)
	t.recalculateDerivedValues()
}

func (t *TempoState) AdjustBPM(delta float64) {
	t.SetBPM(t.BPM + delta)
}

func (t *TempoState) SetSubdivision(sub Subdivision) {
	if sub.IsValid() {
		t.Subdivision = sub
		t.recalculateDerivedValues()
	}
}

func (t *TempoState) NextSubdivision() {
	t.SetSubdivision(t.Subdivision.Next())
}

func (t *TempoState) PrevSubdivision() {
	t.SetSubdivision(t.Subdivision.Prev())
}

func (t *TempoState) RegisterTap(now time.Time) bool {

	if len(t.tapTimes) > 0 {
		lastTap := t.tapTimes[len(t.tapTimes)-1]
		if now.Sub(lastTap) > TapResetThreshold {
			t.tapTimes = t.tapTimes[:0]
		}
	}

	t.tapTimes = append(t.tapTimes, now)

	if len(t.tapTimes) > MaxTapHistory {
		t.tapTimes = t.tapTimes[1:]
	}

	if len(t.tapTimes) < 2 {
		return false
	}

	var totalInterval time.Duration
	for i := 1; i < len(t.tapTimes); i++ {
		totalInterval += t.tapTimes[i].Sub(t.tapTimes[i-1])
	}
	avgInterval := totalInterval / time.Duration(len(t.tapTimes)-1)

	if avgInterval > 0 {
		bpm := 60.0 / avgInterval.Seconds()
		t.SetBPM(bpm)
		return true
	}

	return false
}

func (t *TempoState) ResetTapTempo() {
	t.tapTimes = t.tapTimes[:0]
}

func (t *TempoState) TapCount() int {
	return len(t.tapTimes)
}

func (t *TempoState) recalculateDerivedValues() {

	samplesPerBeat := float64(t.SampleRate) * 60.0 / t.BPM
	t.SamplesPerBeat = int64(samplesPerBeat)

	if t.Subdivision > 0 {
		t.SamplesPerSlot = t.SamplesPerBeat / int64(t.Subdivision)
	} else {
		t.SamplesPerSlot = t.SamplesPerBeat / int64(Sub8)
	}

	if t.SamplesPerSlot < 1 {
		t.SamplesPerSlot = 1
	}
}

func (t *TempoState) GetBeatDuration() time.Duration {
	return time.Duration(float64(time.Minute) / t.BPM)
}

func (t *TempoState) GetSlotDuration() time.Duration {
	return t.GetBeatDuration() / time.Duration(t.Subdivision)
}

func clampBPM(bpm float64) float64 {
	return math.Max(MinBPM, math.Min(MaxBPM, bpm))
}
