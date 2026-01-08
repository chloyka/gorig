# Log Keys Convention

This package provides typed logging keys for structured logging with Uber's zap.

## Why Typed Keys?

1. **Type Safety**: Each key enforces a specific value type at compile time
2. **Consistency**: All log keys follow the same naming convention
3. **Discoverability**: IDE autocompletion shows available keys
4. **Maintainability**: Centralized key definitions prevent typos and inconsistencies

## Naming Convention

Keys follow the dot notation pattern: `domain.key_name`

- **Domain**: Category of the key (e.g., `device`, `audio`, `effect`, `ui`, `path`)
- **Key Name**: Specific identifier in snake_case

Examples:
- `device.input_name`
- `audio.sample_rate`
- `effect.enabled`
- `ui.key`
- `path.config`

## Usage

### Basic Usage

```go
import "github.com/chloyka/distortion-pedal-macos/internal/logger/keys"

// Instead of:
logger.Info("device found", zap.String("name", deviceName))

// Use:
logger.Info("device found", keys.DeviceName(deviceName))
```

### With Context (log.With)

```go
log := logger.With(
    keys.AudioSampleRate(cfg.SampleRate),
    keys.AudioFramesPerBuffer(cfg.FramesPerBuffer),
)
log.Info("audio engine started")
```

### Error Logging

```go
if err != nil {
    logger.Error("failed to load config", keys.Error(err))
}
```

## Available Key Types

| Type | Function | Example |
|------|----------|---------|
| `StringKey` | `keys.String(name)` | `keys.DeviceName = String("device.name")` |
| `IntKey` | `keys.Int(name)` | `keys.UIWidth = Int("ui.width")` |
| `BoolKey` | `keys.Bool(name)` | `keys.EffectEnabled = Bool("effect.enabled")` |
| `Float64Key` | `keys.Float64(name)` | `keys.AudioInputLatencyMs = Float64("audio.input_latency_ms")` |
| `StringsKey` | `keys.Strings(name)` | `keys.EffectList = Strings("effect.list")` |

## Adding New Keys

1. Choose the appropriate domain file (e.g., `device.go`, `audio.go`)
2. Add a new variable with documentation:

```go
// DeviceNewKey logs the new device property.
var DeviceNewKey = String("device.new_key")
```

3. If you need a new domain, create a new file (e.g., `newdomain.go`)

## Log Message Convention

- Messages start with a lowercase letter (unless an abbreviation)
- Messages describe the event, not the data
- Keep messages concise

```go
// Good
logger.Info("device connected", keys.DeviceName(name))
logger.Debug("config saved", keys.PathConfig(path))

// Bad
logger.Info("Device connected", keys.DeviceName(name))  // Uppercase
logger.Info("the device name is", keys.DeviceName(name)) // Data description
```

## Linter

The project uses `depguard` to prevent direct usage of `zap.String`, `zap.Int`, etc. outside the keys package. Run `golangci-lint run` to check for violations.
