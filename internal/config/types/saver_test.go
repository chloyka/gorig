package configTypes

import "testing"

func TestConfigSaver(t *testing.T) {
	t.Run("SetSaveChan", func(t *testing.T) {
		t.Run("should store the channel", func(t *testing.T) {
			sut := &configSaver{}
			saveChan := make(chan struct{}, 1)

			sut.SetSaveChan(saveChan)

			if sut.saveChan == nil {
				t.Error("expected saveChan to be set")
			}
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("should send signal to channel", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &configSaver{saveChan: saveChan}

			sut.Save()

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal to be sent")
			}
		})

		t.Run("should not block when channel buffer is full", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			saveChan <- struct{}{}
			sut := &configSaver{saveChan: saveChan}

			sut.Save()
		})

		t.Run("should not panic when channel is nil", func(t *testing.T) {
			sut := &configSaver{saveChan: nil}

			sut.Save()
		})
	})
}
