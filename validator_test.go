package visivalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type productWithNilMapper struct{}

func newProductWithNilMapper() *productWithNilMapper {
	return &productWithNilMapper{}
}

func (p productWithNilMapper) ValidationMap() map[string]ValidateFunc {
	return nil
}

type productWithEmptyMapper struct{}

func newProductWithEmptyMapper() *productWithEmptyMapper {
	return &productWithEmptyMapper{}
}

func (p productWithEmptyMapper) ValidationMap() map[string]ValidateFunc {
	return map[string]ValidateFunc{}
}

type product struct {
	name      string
	unitPrice price
}

func newProduct(name string, unitPrice int) *product {
	return &product{
		name:      name,
		unitPrice: price(unitPrice),
	}
}

func (p product) validateName() error {
	if len(p.name) <= 0 {
		return errors.New("name must not be empty")
	}

	return nil
}

func (p product) ValidationMap() map[string]ValidateFunc {
	return map[string]ValidateFunc{
		"name":  p.validateName,
		"price": p.unitPrice.validate,
	}
}

type price int

func (p price) validate() error {
	if p <= 0 {
		return errors.New("price must be positive")
	}

	return nil
}

func TestValidate(t *testing.T) {
	type in struct {
		validator ValidationMapper
		mask      []string
	}

	type out struct {
		errs []string
	}

	type validatorState struct {
		validator ValidationMapper
		state     string
	}

	validatorStates := []validatorState{{
		state:     "validator is nil",
		validator: nil,
	}, {
		state:     "validator with nil mapper",
		validator: newProductWithNilMapper(),
	}, {
		state:     "validator with empty mapper",
		validator: newProductWithEmptyMapper(),
	}, {
		state:     "validator with invalid values",
		validator: newProduct("", 0),
	}, {
		state:     "validator with valid values",
		validator: newProduct("my product", 123),
	}}

	type maskState struct {
		mask  []string
		state string
	}

	maskStates := []maskState{{
		state: "mask is nil",
		mask:  nil,
	}, {
		state: "mask is empty",
		mask:  []string{},
	}, {
		state: "mask with invalid field",
		mask:  []string{"id"},
	}, {
		state: "mask with all fields",
		mask:  []string{"name", "price"},
	}, {
		state: "mask with some fields",
		mask:  []string{"price"},
	}}

	tests := []struct {
		message string
		in      in
		out     out
	}{{
		message: fmt.Sprintf("%s", validatorStates[0].state),
		in:      in{validator: validatorStates[0].validator},
		out: out{errs: []string{
			ErrNilMapper.Error(),
		}},
	}, {
		message: fmt.Sprintf("%s", validatorStates[1].state),
		in:      in{validator: validatorStates[1].validator},
		out: out{errs: []string{
			ErrNilOrEmptyValidationMap.Error(),
		}},
	}, {
		message: fmt.Sprintf("%s", validatorStates[2].state),
		in:      in{validator: validatorStates[2].validator},
		out: out{errs: []string{
			ErrNilOrEmptyValidationMap.Error(),
		}},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[3].state, maskStates[0].state),
		in:      in{validator: validatorStates[3].validator, mask: maskStates[0].mask},
		out: out{errs: []string{
			"field name is invalid: name must not be empty",
			"field price is invalid: price must be positive",
		}},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[3].state, maskStates[1].state),
		in:      in{validator: validatorStates[3].validator, mask: maskStates[1].mask},
		out: out{errs: []string{
			"field name is invalid: name must not be empty",
			"field price is invalid: price must be positive",
		}},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[3].state, maskStates[2].state),
		in:      in{validator: validatorStates[3].validator, mask: maskStates[2].mask},
		out: out{errs: []string{
			"field id is unknown",
		}},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[3].state, maskStates[3].state),
		in:      in{validator: validatorStates[3].validator, mask: maskStates[3].mask},
		out: out{errs: []string{
			"field name is invalid: name must not be empty",
			"field price is invalid: price must be positive",
		}},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[3].state, maskStates[4].state),
		in:      in{validator: validatorStates[3].validator, mask: maskStates[4].mask},
		out: out{errs: []string{
			"field price is invalid: price must be positive",
		}},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[4].state, maskStates[0].state),
		in:      in{validator: validatorStates[4].validator, mask: maskStates[0].mask},
		out:     out{errs: nil},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[4].state, maskStates[1].state),
		in:      in{validator: validatorStates[4].validator, mask: maskStates[1].mask},
		out:     out{errs: nil},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[4].state, maskStates[2].state),
		in:      in{validator: validatorStates[4].validator, mask: maskStates[2].mask},
		out: out{errs: []string{
			"field id is unknown",
		}},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[4].state, maskStates[3].state),
		in:      in{validator: validatorStates[4].validator, mask: maskStates[3].mask},
		out:     out{errs: nil},
	}, {
		message: fmt.Sprintf("%s • %s", validatorStates[4].state, maskStates[4].state),
		in:      in{validator: validatorStates[4].validator, mask: maskStates[4].mask},
		out:     out{errs: nil},
	}}

	for _, test := range tests {
		err := Validate(test.in.validator, test.in.mask...)

		if test.out.errs != nil {
			var errs []string
			if ve, ok := err.(ValidationErrors); ok {
				for _, e := range ve {
					errs = append(errs, e.Error())
				}
			} else {
				errs = []string{err.Error()}
			}

			for i, e := range test.out.errs {
				assert.Equal(t, e, errs[i], test.message)
			}
		} else {
			assert.Nil(t, err, test.message)
		}
	}
}
