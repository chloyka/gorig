package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/chloyka/gorig/internal/rhythm"
)

const (
	beatsToShow = 4

	hitMarkerDecay = 500 * time.Millisecond
)

type hitMarker struct {
	slotIndex int
	timestamp time.Time
}

type rhythmVisualizer struct {
	bpm         float64
	subdivision rhythm.Subdivision
	phase       float64
	beatCount   int64
	hitMarkers  []hitMarker
	width       int
}

func newRhythmVisualizer() rhythmVisualizer {
	return rhythmVisualizer{
		bpm:         rhythm.DefaultBPM,
		subdivision: rhythm.Sub8,
		hitMarkers:  make([]hitMarker, 0),
		width:       80,
	}
}

func (v *rhythmVisualizer) Update(bpm float64, sub rhythm.Subdivision, phase float64, beatCount int64) {
	v.bpm = bpm
	v.subdivision = sub
	v.phase = phase
	v.beatCount = beatCount
	v.cleanupExpiredMarkers()
}

func (v *rhythmVisualizer) AddHitMarker(beatInDisplay int, slotInBeat int) {
	globalSlot := beatInDisplay*int(v.subdivision) + slotInBeat
	v.hitMarkers = append(v.hitMarkers, hitMarker{
		slotIndex: globalSlot,
		timestamp: time.Now(),
	})
}

func (v *rhythmVisualizer) cleanupExpiredMarkers() {
	now := time.Now()
	active := make([]hitMarker, 0, len(v.hitMarkers))
	for _, m := range v.hitMarkers {
		if now.Sub(m.timestamp) < hitMarkerDecay {
			active = append(active, m)
		}
	}
	v.hitMarkers = active
}

func (v *rhythmVisualizer) isSlotHit(globalSlot int) bool {
	for _, m := range v.hitMarkers {
		if m.slotIndex == globalSlot {
			return true
		}
	}
	return false
}

func (v *rhythmVisualizer) SetWidth(width int) {
	v.width = width
}

func (v *rhythmVisualizer) View() string {
	if v.width < 40 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" BPM: %-3.0f  [%s]  TAP:[t]  +/-:[,/.]  Sub:[/]\n",
		v.bpm, v.subdivision.String()))

	totalSlots := beatsToShow * int(v.subdivision)

	sb.WriteString(" ")
	for beat := 0; beat < beatsToShow; beat++ {
		for slot := 0; slot < int(v.subdivision); slot++ {
			globalSlot := beat*int(v.subdivision) + slot

			if slot == 0 {

				sb.WriteString("|")
			} else if v.isSlotHit(globalSlot) {

				sb.WriteString("!")
			} else {

				sb.WriteString(".")
			}
		}
	}
	sb.WriteString("|\n")

	currentBeatInDisplay := int(v.beatCount) % beatsToShow
	slotInBeat := int(v.phase * float64(v.subdivision))
	if slotInBeat >= int(v.subdivision) {
		slotInBeat = int(v.subdivision) - 1
	}
	playheadPos := currentBeatInDisplay*int(v.subdivision) + slotInBeat

	sb.WriteString(" ")
	for i := 0; i <= totalSlots; i++ {
		if i == playheadPos {
			sb.WriteString("^")
		} else {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("\n")

	return sb.String()
}
