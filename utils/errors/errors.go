package errors

import (
	"errors"
	"fmt"
)

func Wrap(sentinel error, context any) error {
	if sentinel == nil {
		return nil
	}
	if context == nil {
		return sentinel
	}
	if err, ok := context.(error); ok {
		return fmt.Errorf("%w: %w", sentinel, err)
	}
	return fmt.Errorf("%w: %v", sentinel, context)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func New(text string) error {
	return errors.New(text)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
