package api

import (
	"context"
	"strings"

	errs "github.com/chloyka/gorig/utils/errors"
)

type MockEffectRepository struct {
	effects []EffectDefinition
}

func NewMockEffectRepository() *MockEffectRepository {
	return &MockEffectRepository{
		effects: []EffectDefinition{
			{
				ID:          "dist-001",
				Name:        "Heavy Distortion",
				Description: "High-gain distortion for metal tones",
				Author:      "EffectsCommunity",
				Version:     "1.0.0",
				Category:    "distortion",
				DownloadURL: "http://localhost",
			},
			{
				ID:          "chorus-001",
				Name:        "Analog Chorus",
				Description: "Warm analog-style chorus effect",
				Author:      "EffectsCommunity",
				Version:     "1.0.0",
				Category:    "modulation",
				DownloadURL: "http://localhost",
			},
			{
				ID:          "delay-001",
				Name:        "Digital Delay",
				Description: "Clean digital delay with feedback control",
				Author:      "EffectsCommunity",
				Version:     "1.0.0",
				Category:    "delay",
				DownloadURL: "http://localhost",
			},
			{
				ID:          "reverb-001",
				Name:        "Hall Reverb",
				Description: "Spacious hall reverb effect",
				Author:      "EffectsCommunity",
				Version:     "1.0.0",
				Category:    "reverb",
				DownloadURL: "http://localhost",
			},
			{
				ID:          "comp-001",
				Name:        "Studio Compressor",
				Description: "Transparent studio-grade compressor",
				Author:      "EffectsCommunity",
				Version:     "1.0.0",
				Category:    "dynamics",
				DownloadURL: "http://localhost",
			},
		},
	}
}

func (m *MockEffectRepository) SearchEffects(ctx context.Context, filters SearchFilters) ([]EffectDefinition, error) {
	if filters.Query == "" && filters.Category == "" {
		return m.effects, nil
	}

	var results []EffectDefinition
	for _, e := range m.effects {
		if filters.Query != "" {
			if containsIgnoreCase(e.Name, filters.Query) ||
				containsIgnoreCase(e.Description, filters.Query) {
				results = append(results, e)
				continue
			}
		}
		if filters.Category != "" && strings.EqualFold(e.Category, filters.Category) {
			results = append(results, e)
		}
	}
	return results, nil
}

func (m *MockEffectRepository) GetEffect(ctx context.Context, id string) (*EffectDefinition, error) {
	for _, e := range m.effects {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, errs.Wrap(errs.ErrAPIEffectNotFound, id)
}

func (m *MockEffectRepository) DownloadEffect(ctx context.Context, id string, targetDir string) error {

	return errs.Wrap(errs.ErrAPINotImplemented, []string{id, targetDir})
}

type MockPresetRepository struct {
	presets []PresetDefinition
}

func NewMockPresetRepository() *MockPresetRepository {
	return &MockPresetRepository{
		presets: []PresetDefinition{
			{
				ID:          "preset-001",
				Name:        "Metal Tone",
				Description: "High gain setup for metal and djent",
				Author:      "PresetsCommunity",
				EffectChain: []string{"noise gate", "overdrive", "heavy distortion", "eq"},
				Tags:        []string{"metal", "high-gain", "djent"},
			},
			{
				ID:          "preset-002",
				Name:        "Clean Blues",
				Description: "Warm clean tone with subtle overdrive",
				Author:      "PresetsCommunity",
				EffectChain: []string{"compressor", "light overdrive", "chorus", "reverb"},
				Tags:        []string{"blues", "clean", "warm"},
			},
			{
				ID:          "preset-003",
				Name:        "80s Rock",
				Description: "Classic 80s rock tone with chorus and delay",
				Author:      "PresetsCommunity",
				EffectChain: []string{"overdrive", "chorus", "delay"},
				Tags:        []string{"rock", "80s", "classic"},
			},
			{
				ID:          "preset-004",
				Name:        "Ambient Soundscape",
				Description: "Dreamy ambient textures",
				Author:      "PresetsCommunity",
				EffectChain: []string{"reverb", "delay", "chorus", "tremolo"},
				Tags:        []string{"ambient", "atmospheric", "dreamy"},
			},
		},
	}
}

func (m *MockPresetRepository) SearchPresets(ctx context.Context, filters SearchFilters) ([]PresetDefinition, error) {
	if filters.Query == "" && len(filters.Tags) == 0 {
		return m.presets, nil
	}

	var results []PresetDefinition
	for _, p := range m.presets {
		if filters.Query != "" {
			if containsIgnoreCase(p.Name, filters.Query) ||
				containsIgnoreCase(p.Description, filters.Query) {
				results = append(results, p)
				continue
			}
		}

		if len(filters.Tags) > 0 {
			for _, filterTag := range filters.Tags {
				for _, presetTag := range p.Tags {
					if strings.EqualFold(filterTag, presetTag) {
						results = append(results, p)
						break
					}
				}
			}
		}
	}
	return results, nil
}

func (m *MockPresetRepository) GetPreset(ctx context.Context, id string) (*PresetDefinition, error) {
	for _, p := range m.presets {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, errs.Wrap(errs.ErrAPIPresetNotFound, id)
}

func (m *MockPresetRepository) ImportPreset(ctx context.Context, id string) (*PresetDefinition, error) {
	p, err := m.GetPreset(ctx, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (m *MockPresetRepository) SharePreset(ctx context.Context, preset PresetDefinition) (string, error) {

	return "shared-" + strings.ReplaceAll(preset.Name, " ", "-"), nil
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
