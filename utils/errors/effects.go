package errors

var (
	ErrEffectsWalkDir       = New("effects: failed to walk directory")
	ErrEffectsReadFile      = New("effects: failed to read file")
	ErrEffectsStdlib        = New("effects: failed to use stdlib")
	ErrEffectsEval          = New("effects: failed to eval")
	ErrEffectsGetName       = New("effects: failed to get Name")
	ErrEffectsGetProcess    = New("effects: failed to get Process")
	ErrEffectsDuplicateName = New("effects: duplicate effect name")
	ErrEffectsLoad          = New("effects: failed to load")
)
