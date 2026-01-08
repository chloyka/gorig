# Logging Conventions

## Typed Keys

Use typed keys from `internal/logger/keys` instead of raw `zap.String`, `zap.Int`, etc.

```go
import "github.com/chloyka/distortion-pedal-macos/internal/logger/keys"

// Good
logger.Info("device found", keys.DeviceName(name))

// Bad
logger.Info("device found", zap.String("name", name))
```

## Key Naming

Format: `domain.key_name` (snake_case)

| Domain | Example |
|--------|---------|
| device | `device.input_name`, `device.output_count` |
| audio | `audio.sample_rate`, `audio.input_latency_ms` |
| effect | `effect.name`, `effect.enabled` |
| path | `path.config`, `path.dir` |
| ui | `ui.key`, `ui.width` |

## Log Levels

| Level | When to use |
|-------|-------------|
| Debug | Internal state changes, loop iterations, intermediate steps |
| Info | Significant events: startup, shutdown, device changes, config saves |
| Warn | Recoverable issues: fallback to default, missing optional config |
| Error | Failed operations that affect functionality |

## What to Log

### Must log (Info)
- Component startup/shutdown
- Device selection changes
- Configuration saves
- External resource connections

### Should log (Debug)
- Signal received in event loops
- Intermediate processing steps

### Must log (Error)
- All errors with `keys.Error(err)` and relevant context

## Message Format

- Lowercase start (unless abbreviation)
- Describe the event, not the data
- Concise

```go
// Good
logger.Info("config saved", keys.PathConfig(path))
logger.Warn("saved device not found", keys.DeviceSavedName(saved), keys.DeviceUsingName(fallback))
logger.Error("failed to open stream", keys.Error(err))

// Bad
logger.Info("Config Saved", ...)           // uppercase
logger.Info("the path is", ...)            // describes data
logger.Error("error", keys.Error(err))     // not specific
```

## Context Enrichment

Use `logger.With()` when logging multiple related entries:

```go
log := logger.With(
    keys.AudioSampleRate(cfg.SampleRate),
    keys.AudioFramesPerBuffer(cfg.FramesPerBuffer),
)
log.Info("audio engine started")
```

## Adding New Keys

1. Find or create domain file in `internal/logger/keys/`
2. Add typed key:

```go
var DeviceNewKey = String("device.new_key")
```

Key types: `String`, `Int`, `Bool`, `Float64`, `Strings`

## Linter

`depguard` blocks direct `zap.String/Int/...` usage outside keys package.
