package rhythm

import (
	"fmt"

	"github.com/chloyka/gorig/internal/onset"
)

type Subdivision int

const (
	Sub2  Subdivision = 2
	Sub4  Subdivision = 4
	Sub8  Subdivision = 8
	Sub16 Subdivision = 16
	Sub32 Subdivision = 32
	Sub64 Subdivision = 64
)

func (s Subdivision) String() string {
	return fmt.Sprintf("1/%d", s)
}

func (s Subdivision) Next() Subdivision {
	switch s {
	case Sub2:
		return Sub4
	case Sub4:
		return Sub8
	case Sub8:
		return Sub16
	case Sub16:
		return Sub32
	case Sub32:
		return Sub64
	case Sub64:
		return Sub2
	default:
		return Sub8
	}
}

func (s Subdivision) Prev() Subdivision {
	switch s {
	case Sub64:
		return Sub32
	case Sub32:
		return Sub16
	case Sub16:
		return Sub8
	case Sub8:
		return Sub4
	case Sub4:
		return Sub2
	case Sub2:
		return Sub64
	default:
		return Sub8
	}
}

func (s Subdivision) IsValid() bool {
	switch s {
	case Sub2, Sub4, Sub8, Sub16, Sub32, Sub64:
		return true
	default:
		return false
	}
}

func SubdivisionFromInt(val int) Subdivision {
	s := Subdivision(val)
	if s.IsValid() {
		return s
	}
	return Sub8
}

type QuantizedOnset struct {
	OriginalEvent onset.Event
	BeatPosition  float64
	SlotIndex     int
	WasQueued     bool
}
