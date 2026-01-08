package api

import "context"

type EffectDefinition struct {
	ID          string
	Name        string
	Description string
	Author      string
	Version     string
	Category    string
	DownloadURL string
}

type PresetDefinition struct {
	ID          string
	Name        string
	Description string
	Author      string
	EffectChain []string
	Tags        []string
}

type SearchFilters struct {
	Query    string
	Category string
	Author   string
	Tags     []string
	Limit    int
	Offset   int
}

type EffectRepository interface {
	SearchEffects(ctx context.Context, filters SearchFilters) ([]EffectDefinition, error)

	GetEffect(ctx context.Context, id string) (*EffectDefinition, error)

	DownloadEffect(ctx context.Context, id string, targetDir string) error
}

type PresetRepository interface {
	SearchPresets(ctx context.Context, filters SearchFilters) ([]PresetDefinition, error)

	GetPreset(ctx context.Context, id string) (*PresetDefinition, error)

	ImportPreset(ctx context.Context, id string) (*PresetDefinition, error)

	SharePreset(ctx context.Context, preset PresetDefinition) (string, error)
}
