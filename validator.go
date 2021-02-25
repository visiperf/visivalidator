package visivalidator

import (
	"errors"
	"fmt"
	"sort"
)

// ErrNilMapper is error occured when mapper is nil
var ErrNilMapper = errors.New("mapper cannot be nil")

// ErrNilOrEmptyValidationMap is error occured when validation map is nil or empty
var ErrNilOrEmptyValidationMap = errors.New("validation map cannot be nil or empty")

// ValidateFunc is type representing a func validating something
type ValidateFunc func() error

// ValidationMapper is interface implemented by types that can map fields with validate funcs
type ValidationMapper interface {
	ValidationMap() map[string]ValidateFunc
}

// Validator is interface implemented by types that can be validated
type Validator interface {
	Validate() error
}

// Validate return stacked error of validator using field mask
// if field mask is empty, the operation applies to all fields (as if a field mask of all fields has been specified).
func Validate(mapper ValidationMapper, mask ...string) error {
	vm, err := getValidationMap(mapper)
	if err != nil {
		return err
	}

	fields := mask
	if len(fields) <= 0 {
		fields = extractMapKeys(vm)
	}

	sort.Strings(fields)

	var errs []error
	for _, f := range fields {
		fn, ok := vm[f]
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

func getValidationMap(mapper ValidationMapper) (map[string]ValidateFunc, error) {
	if mapper == nil {
		return nil, ErrNilMapper
	}

	vm := mapper.ValidationMap()
	if len(vm) <= 0 {
		return nil, ErrNilOrEmptyValidationMap
	}

	return vm, nil
}

func extractMapKeys(vm map[string]ValidateFunc) []string {
	keys := make([]string, 0, len(vm))
	for k := range vm {
		keys = append(keys, k)
	}

	return keys
}
