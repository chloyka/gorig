package onset

import (
	"math"
	"sync"
)

type Event struct {
	Energy    float32
	Timestamp int64
}

type Detector struct {
	mu sync.RWMutex

	threshold     float32
	minEnergy     float32
	attackCoeff   float32
	releaseCoeff  float32
	minIntervalMs float32
	sampleRate    float32
	samplesPerMs  float32

	energyFollower     float32
	baseline           float32
	samplesSinceLast   int64
	minIntervalSamples int64
	lastEnergy         float32
	currentEnergy      float32

	wasLow     bool
	peakEnergy float32

	onsetChan    chan Event
	totalSamples int64

	enabled bool
}

type Config struct {
	Threshold     float32
	MinEnergy     float32
	AttackMs      float32
	ReleaseMs     float32
	MinIntervalMs float32
	SampleRate    float32
	BufferSize    int
}

func DefaultConfig(sampleRate float32) Config {
	return Config{
		Threshold:     5.0,
		MinEnergy:     0.01,
		AttackMs:      2,
		ReleaseMs:     500,
		MinIntervalMs: 150,
		SampleRate:    sampleRate,
		BufferSize:    4,
	}
}

func NewDetector(cfg Config) *Detector {
	samplesPerMs := cfg.SampleRate / 1000.0

	attackCoeff := float32(1.0 - math.Exp(-1.0/float64(cfg.AttackMs*samplesPerMs)))
	releaseCoeff := float32(1.0 - math.Exp(-1.0/float64(cfg.ReleaseMs*samplesPerMs)))

	return &Detector{
		threshold:          cfg.Threshold,
		minEnergy:          cfg.MinEnergy,
		attackCoeff:        attackCoeff,
		releaseCoeff:       releaseCoeff,
		minIntervalMs:      cfg.MinIntervalMs,
		sampleRate:         cfg.SampleRate,
		samplesPerMs:       samplesPerMs,
		minIntervalSamples: int64(cfg.MinIntervalMs * samplesPerMs),
		onsetChan:          make(chan Event, cfg.BufferSize),
		enabled:            true,
		baseline:           0.01,
		wasLow:             true,
	}
}

func (d *Detector) Process(samples []float32) bool {
	if !d.enabled {
		return false
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	var sumSquares float32
	for _, s := range samples {
		sumSquares += s * s
	}
	bufferEnergy := float32(math.Sqrt(float64(sumSquares / float32(len(samples)))))

	if bufferEnergy < d.minEnergy {
		d.samplesSinceLast += int64(len(samples))
		d.totalSamples += int64(len(samples))

		d.wasLow = true
		d.peakEnergy = 0
		d.lastEnergy = 0
		d.currentEnergy = 0
		d.energyFollower = 0
		return false
	}

	if bufferEnergy > d.energyFollower {
		d.energyFollower += d.attackCoeff * (bufferEnergy - d.energyFollower)
	} else {
		d.energyFollower += d.releaseCoeff * (bufferEnergy - d.energyFollower)
	}

	if bufferEnergy < d.baseline*2 {
		d.baseline += d.releaseCoeff * 0.05 * (bufferEnergy - d.baseline)
	}

	if d.baseline < 0.0001 {
		d.baseline = 0.0001
	}

	d.lastEnergy = d.currentEnergy
	d.currentEnergy = bufferEnergy

	ratio := d.energyFollower / d.baseline

	if bufferEnergy > d.peakEnergy {
		d.peakEnergy = bufferEnergy
	}

	if !d.wasLow && d.peakEnergy > d.minEnergy*5 && bufferEnergy < d.peakEnergy*0.1 {
		d.wasLow = true
		d.peakEnergy = bufferEnergy
	}

	isOnset := ratio > d.threshold &&
		d.wasLow &&
		d.samplesSinceLast >= d.minIntervalSamples &&
		d.currentEnergy > d.lastEnergy*3.0 &&
		d.currentEnergy > d.minEnergy*3

	if isOnset {

		normalizedEnergy := float32(math.Min(float64(ratio/d.threshold/2), 1.0))

		event := Event{
			Energy:    normalizedEnergy,
			Timestamp: d.totalSamples,
		}

		select {
		case d.onsetChan <- event:
		default:

		}

		d.samplesSinceLast = 0
		d.wasLow = false
		d.peakEnergy = bufferEnergy
	} else {
		d.samplesSinceLast += int64(len(samples))
	}

	d.totalSamples += int64(len(samples))

	return isOnset
}

func (d *Detector) Events() <-chan Event {
	return d.onsetChan
}

func (d *Detector) CurrentEnergy() float32 {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.baseline < 0.0001 {
		return 0
	}
	ratio := d.energyFollower / d.baseline / d.threshold
	return float32(math.Min(float64(ratio), 1.0))
}

func (d *Detector) SetThreshold(threshold float32) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.threshold = threshold
}

func (d *Detector) SetEnabled(enabled bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = enabled
}

func (d *Detector) IsEnabled() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.enabled
}

func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.energyFollower = 0
	d.baseline = 0.001
	d.samplesSinceLast = d.minIntervalSamples
	d.lastEnergy = 0
	d.currentEnergy = 0
	d.wasLow = true
	d.peakEnergy = 0
}

func (d *Detector) Close() {
	close(d.onsetChan)
}
