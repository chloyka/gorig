package configTypes

import "testing"

func TestStateConfig(t *testing.T) {
	t.Run("SetInputDevice", func(t *testing.T) {
		t.Run("should update input device and trigger save", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &StateConfig{}
			sut.SetSaveChan(saveChan)

			sut.SetInputDevice("USB Audio")

			if sut.InputDevice != "USB Audio" {
				t.Errorf("got InputDevice=%q, want %q", sut.InputDevice, "USB Audio")
			}

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})
	})

	t.Run("SetOutputDevice", func(t *testing.T) {
		t.Run("should update output device and trigger save", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &StateConfig{}
			sut.SetSaveChan(saveChan)

			sut.SetOutputDevice("Speakers")

			if sut.OutputDevice != "Speakers" {
				t.Errorf("got OutputDevice=%q, want %q", sut.OutputDevice, "Speakers")
			}

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})
	})

	t.Run("SetEffectsEnabled", func(t *testing.T) {
		t.Run("should update effects enabled and trigger save", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &StateConfig{EffectsEnabled: false}
			sut.SetSaveChan(saveChan)

			sut.SetEffectsEnabled(true)

			if !sut.EffectsEnabled {
				t.Error("expected EffectsEnabled=true")
			}

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})

		t.Run("should set to false and trigger save", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &StateConfig{EffectsEnabled: true}
			sut.SetSaveChan(saveChan)

			sut.SetEffectsEnabled(false)

			if sut.EffectsEnabled {
				t.Error("expected EffectsEnabled=false")
			}

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})
	})
}
