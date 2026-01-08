package configTypes

import "testing"

func TestPresetsConfig(t *testing.T) {
	t.Run("SetActivePreset", func(t *testing.T) {
		t.Run("should update active preset and trigger save", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &PresetsConfig{}
			sut.SetSaveChan(saveChan)

			sut.SetActivePreset("test-preset")

			if sut.ActivePreset != "test-preset" {
				t.Errorf("got ActivePreset=%q, want %q", sut.ActivePreset, "test-preset")
			}

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})
	})

	t.Run("AddPreset", func(t *testing.T) {
		t.Run("should append preset to list", func(t *testing.T) {
			sut := &PresetsConfig{Presets: []Preset{}}

			sut.AddPreset(Preset{Name: "new", EffectChain: []string{"dist"}})

			if len(sut.Presets) != 1 {
				t.Fatalf("got len=%d, want 1", len(sut.Presets))
			}
			if sut.Presets[0].Name != "new" {
				t.Errorf("got Name=%q, want %q", sut.Presets[0].Name, "new")
			}
		})

		t.Run("should trigger save", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &PresetsConfig{}
			sut.SetSaveChan(saveChan)

			sut.AddPreset(Preset{Name: "test"})

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})
	})

	t.Run("UpdatePreset", func(t *testing.T) {
		t.Run("should update existing preset chain", func(t *testing.T) {
			sut := &PresetsConfig{
				Presets: []Preset{{Name: "test", EffectChain: []string{"old"}}},
			}

			sut.UpdatePreset("test", []string{"new1", "new2"})

			got := sut.Presets[0].EffectChain
			if len(got) != 2 || got[0] != "new1" || got[1] != "new2" {
				t.Errorf("got EffectChain=%v, want [new1 new2]", got)
			}
		})

		t.Run("should return true on success", func(t *testing.T) {
			sut := &PresetsConfig{
				Presets: []Preset{{Name: "test"}},
			}

			got := sut.UpdatePreset("test", []string{})

			if !got {
				t.Error("expected true")
			}
		})

		t.Run("should return false for non-existent preset", func(t *testing.T) {
			sut := &PresetsConfig{Presets: []Preset{}}

			got := sut.UpdatePreset("missing", []string{})

			if got {
				t.Error("expected false")
			}
		})

		t.Run("should trigger save on success", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &PresetsConfig{Presets: []Preset{{Name: "test"}}}
			sut.SetSaveChan(saveChan)

			sut.UpdatePreset("test", []string{"updated"})

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})
	})

	t.Run("DeletePreset", func(t *testing.T) {
		t.Run("should remove preset from list", func(t *testing.T) {
			sut := &PresetsConfig{
				Presets: []Preset{{Name: "keep"}, {Name: "delete"}, {Name: "also-keep"}},
			}

			sut.DeletePreset("delete")

			if len(sut.Presets) != 2 {
				t.Fatalf("got len=%d, want 2", len(sut.Presets))
			}
			if sut.Presets[0].Name != "keep" || sut.Presets[1].Name != "also-keep" {
				t.Errorf("wrong presets remaining: %v", sut.Presets)
			}
		})

		t.Run("should return true on success", func(t *testing.T) {
			sut := &PresetsConfig{Presets: []Preset{{Name: "test"}}}

			got := sut.DeletePreset("test")

			if !got {
				t.Error("expected true")
			}
		})

		t.Run("should return false for non-existent preset", func(t *testing.T) {
			sut := &PresetsConfig{Presets: []Preset{}}

			got := sut.DeletePreset("missing")

			if got {
				t.Error("expected false")
			}
		})

		t.Run("should trigger save on success", func(t *testing.T) {
			saveChan := make(chan struct{}, 1)
			sut := &PresetsConfig{Presets: []Preset{{Name: "test"}}}
			sut.SetSaveChan(saveChan)

			sut.DeletePreset("test")

			select {
			case <-saveChan:

			default:
				t.Error("expected save signal")
			}
		})
	})

	t.Run("GetPreset", func(t *testing.T) {
		t.Run("should return preset by name", func(t *testing.T) {
			sut := &PresetsConfig{
				Presets: []Preset{
					{Name: "first", EffectChain: []string{"a"}},
					{Name: "second", EffectChain: []string{"b"}},
				},
			}

			got := sut.GetPreset("second")

			if got == nil {
				t.Fatal("expected non-nil preset")
			}
			if got.Name != "second" {
				t.Errorf("got Name=%q, want %q", got.Name, "second")
			}
		})

		t.Run("should return nil for non-existent preset", func(t *testing.T) {
			sut := &PresetsConfig{Presets: []Preset{{Name: "other"}}}

			got := sut.GetPreset("missing")

			if got != nil {
				t.Errorf("expected nil, got %v", got)
			}
		})
	})

	t.Run("GetActivePresetConfig", func(t *testing.T) {
		t.Run("should return active preset", func(t *testing.T) {
			sut := &PresetsConfig{
				Presets:      []Preset{{Name: "active", EffectChain: []string{"fx"}}},
				ActivePreset: "active",
			}

			got := sut.GetActivePresetConfig()

			if got == nil {
				t.Fatal("expected non-nil preset")
			}
			if got.Name != "active" {
				t.Errorf("got Name=%q, want %q", got.Name, "active")
			}
		})

		t.Run("should return nil when no active preset set", func(t *testing.T) {
			sut := &PresetsConfig{
				Presets:      []Preset{{Name: "test"}},
				ActivePreset: "",
			}

			got := sut.GetActivePresetConfig()

			if got != nil {
				t.Errorf("expected nil, got %v", got)
			}
		})

		t.Run("should return nil when active preset not found", func(t *testing.T) {
			sut := &PresetsConfig{
				Presets:      []Preset{{Name: "other"}},
				ActivePreset: "missing",
			}

			got := sut.GetActivePresetConfig()

			if got != nil {
				t.Errorf("expected nil, got %v", got)
			}
		})
	})
}
