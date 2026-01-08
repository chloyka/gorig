# Testing Conventions

## Test Structure

Tests are organized as suites using `testing.T` and `t.Run()`. Each test name must start with "should".

```go
func TestComponentName(t *testing.T) {
    t.Run("MethodName", func(t *testing.T) {
        t.Run("should do something specific", func(t *testing.T) {
            // test implementation
        })

        t.Run("should handle edge case", func(t *testing.T) {
            // test implementation
        })
    })
}
```

## Test Naming

### Function Names
- `TestComponentName` - main test function for a component
- `TestComponentName_Integration` - integration tests (if needed)

### Subtest Names (t.Run)
- First level: method or behavior group name (`"AddPreset"`, `"Save"`)
- Second level: must start with `"should"` (`"should add preset to list"`)

```go
// Good
t.Run("should return nil for non-existent preset", ...)
t.Run("should trigger save after update", ...)

// Bad
t.Run("test add preset", ...)
t.Run("AddPresetWorks", ...)
t.Run("returns nil", ...)
```

## Variable Naming

Use semantic names that describe the role:

| Variable | Purpose |
|----------|---------|
| `sut` | System Under Test - the main object being tested |
| `got` | Actual result from the function |
| `want` | Expected result |
| `err` | Error returned from function |
| `saveChan` | Channel for testing save signals |

```go
func TestPresetsConfig(t *testing.T) {
    t.Run("GetPreset", func(t *testing.T) {
        t.Run("should return preset by name", func(t *testing.T) {
            sut := &PresetsConfig{
                Presets: []Preset{{Name: "test", EffectChain: []string{"dist"}}},
            }

            got := sut.GetPreset("test")

            want := &Preset{Name: "test", EffectChain: []string{"dist"}}
            if got.Name != want.Name {
                t.Errorf("got %v, want %v", got, want)
            }
        })
    })
}
```

## What to Test

### Unit Tests
- Public methods and their behavior
- Edge cases (nil, empty, not found)
- Side effects (Save() calls)
- Error conditions

### Integration Tests
- Component interaction
- FX dependency injection

## Testing Save Signals

For components that trigger `Save()`, use a buffered channel:

```go
func TestStateConfig(t *testing.T) {
    t.Run("SetInputDevice", func(t *testing.T) {
        t.Run("should trigger save", func(t *testing.T) {
            saveChan := make(chan struct{}, 1)
            sut := &StateConfig{}
            sut.SetSaveChan(saveChan)

            sut.SetInputDevice("USB Audio")

            select {
            case <-saveChan:
                // success
            default:
                t.Error("expected save signal")
            }
        })
    })
}
```

## Table-Driven Tests

Use table-driven tests for multiple similar cases:

```go
func TestGetPreset(t *testing.T) {
    tests := []struct {
        name       string
        presets    []Preset
        searchName string
        wantNil    bool
    }{
        {
            name:       "should return preset when exists",
            presets:    []Preset{{Name: "test"}},
            searchName: "test",
            wantNil:    false,
        },
        {
            name:       "should return nil when not exists",
            presets:    []Preset{{Name: "other"}},
            searchName: "test",
            wantNil:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sut := &PresetsConfig{Presets: tt.presets}
            got := sut.GetPreset(tt.searchName)

            if (got == nil) != tt.wantNil {
                t.Errorf("got nil=%v, want nil=%v", got == nil, tt.wantNil)
            }
        })
    }
}
```

## FX Testing

Use `fx.ValidateApp()` to verify dependency injection graph:

```go
func TestFXIntegration(t *testing.T) {
    t.Run("should validate complete app dependencies", func(t *testing.T) {
        err := fx.ValidateApp(
            config.ProvideConfig(),
            logger.Module,
            config.ProvideManager(),
            // ... other modules
        )
        if err != nil {
            t.Fatalf("fx validation failed: %v", err)
        }
    })
}
```

This single test validates:
- All providers return correct types
- All dependencies are satisfied
- FX groups are populated correctly
- No circular dependencies

## File Organization

```
internal/
└── config/
    ├── module.go
    ├── module_test.go      # Tests for provideConfig
    ├── saver.go
    ├── saver_test.go       # Tests for ConfigManager
    ├── loader.go
    ├── loader_test.go      # Tests for file discovery
    └── types/
        ├── presets.go
        ├── presets_test.go
        ├── state.go
        ├── state_test.go
        ├── saver.go
        └── saver_test.go
```
