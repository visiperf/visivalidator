package visivalidator

import (
	"errors"
	"fmt"
	"sort"
)

// ErrNilValidator is error occured when validator is nil
var ErrNilValidator = errors.New("validator cannot be nil")

// ErrNilOrEmptyMapper is error occured when mapper is nil or empty
var ErrNilOrEmptyMapper = errors.New("mapper cannot be nil or empty")

// ValidateFunc is type representing a func validating something
type ValidateFunc func() error

// Validator is interface implemented by types that can be validated
type Validator interface {
	ValidationMapper() map[string]ValidateFunc
}

// Validate return stacked error of validator using field mask
// if field mask is empty, the operation applies to all fields (as if a field mask of all fields has been specified).
func Validate(validator Validator, mask ...string) error {
	mapper, err := getMapperFromValidator(validator)
	if err != nil {
		return err
	}

	fields := mask
	if len(fields) <= 0 {
		fields = extractMapKeys(mapper)
	}

	sort.Strings(fields)

	var errs []error
	for _, f := range fields {
		fn, ok := mapper[f]
		if !ok {
			errs = append(errs, fmt.Errorf("field %s is unknown", f))
			continue
		}

		if err := fn(); err != nil {
			errs = append(errs, fmt.Errorf("field %s is invalid: %w", f, err))
		}
	}

	if len(errs) > 0 {
		return ValidationErrors(errs)
	}

	return nil
}

func getMapperFromValidator(validator Validator) (map[string]ValidateFunc, error) {
	if validator == nil {
		return nil, ErrNilValidator
	}

	mapper := validator.ValidationMapper()
	if len(mapper) <= 0 {
		return nil, ErrNilOrEmptyMapper
	}

	return mapper, nil
}

func extractMapKeys(m map[string]ValidateFunc) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
