package effects

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	errs "github.com/chloyka/gorig/utils/errors"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type EffectRegistry struct {
	effects map[string]*InterpretedEffect
}

func NewEffectRegistry() *EffectRegistry {
	return &EffectRegistry{
		effects: make(map[string]*InterpretedEffect),
	}
}

func (r *EffectRegistry) GetEffect(name string) *InterpretedEffect {
	return r.effects[name]
}

func (r *EffectRegistry) GetAvailableEffectNames() []string {
	names := make([]string, 0, len(r.effects))
	for name := range r.effects {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (r *EffectRegistry) Count() int {
	return len(r.effects)
}

func loadEffectsFromDirRecursive(dir string) (*EffectRegistry, error) {
	registry := NewEffectRegistry()

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {

			if os.IsNotExist(err) && path == dir {
				return filepath.SkipAll
			}
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		effect, err := loadEffect(path)
		if err != nil {

			return nil
		}

		if _, exists := registry.effects[effect.Name()]; exists {
			return errs.ErrEffectsDuplicateName
		}

		registry.effects[effect.Name()] = effect.(*InterpretedEffect)
		return nil
	})

	if err != nil {
		return nil, errs.Wrap(errs.ErrEffectsWalkDir, err)
	}

	return registry, nil
}

func loadEffect(filePath string) (Effect, error) {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errs.Wrap(errs.ErrEffectsReadFile, err)
	}

	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, errs.Wrap(errs.ErrEffectsStdlib, err)
	}

	_, err = i.Eval(string(code))
	if err != nil {
		return nil, errs.Wrap(errs.ErrEffectsEval, err)
	}

	nameVal, err := i.Eval("effects.Name")
	if err != nil {
		return nil, errs.Wrap(errs.ErrEffectsGetName, err)
	}
	name := nameVal.Interface().(string)

	processVal, err := i.Eval("effects.Process")
	if err != nil {
		return nil, errs.Wrap(errs.ErrEffectsGetProcess, err)
	}

	return newInterpretedEffect(name, true, processVal), nil
}

func loadEffectsFromDir(dir string) ([]Effect, error) {
	registry, err := loadEffectsFromDirRecursive(dir)
	if err != nil {
		return nil, err
	}

	var effects []Effect
	for _, name := range registry.GetAvailableEffectNames() {
		effects = append(effects, registry.GetEffect(name))
	}
	return effects, nil
}
