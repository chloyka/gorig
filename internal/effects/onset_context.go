package effects

import (
	"sync"
)

type OnsetContext struct {
	HasOnset     bool
	Energy       float32
	BeatPosition float64
	SlotIndex    int
}

var (
	currentOnsetMu sync.RWMutex
	currentOnset   OnsetContext
)

func SetCurrentOnset(hasOnset bool, energy float32, beatPosition float64, slotIndex int) {
	currentOnsetMu.Lock()
	defer currentOnsetMu.Unlock()
	currentOnset = OnsetContext{
		HasOnset:     hasOnset,
		Energy:       energy,
		BeatPosition: beatPosition,
		SlotIndex:    slotIndex,
	}
}

func ClearCurrentOnset() {
	currentOnsetMu.Lock()
	defer currentOnsetMu.Unlock()
	currentOnset = OnsetContext{HasOnset: false}
}

func GetCurrentOnset() OnsetContext {
	currentOnsetMu.RLock()
	defer currentOnsetMu.RUnlock()
	return currentOnset
}
