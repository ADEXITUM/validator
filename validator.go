package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Field string
type Rule string
type ErrorMsg string

type CustomErrors map[Field]map[Rule]ErrorMsg

type ValidationError struct {
	Field   string
	Message ErrorMsg
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Field '%s' validation failed: %s", e.Field, e.Message)
}

type Validator struct {
	customErrors CustomErrors
}

func New() *Validator {
	return &Validator{
		customErrors: make(CustomErrors),
	}
}

func (v *Validator) WithCustomErrors(errors CustomErrors) *Validator {
	for field, validationErrors := range errors {
		if _, exists := v.customErrors[field]; !exists {
			v.customErrors[field] = make(map[Rule]ErrorMsg)
		}
		for validationType, message := range validationErrors {
			v.customErrors[field][validationType] = message
		}
	}
	return v
}

func (v *Validator) Validate(i interface{}) error {
	val := reflect.ValueOf(i)
	typ := reflect.TypeOf(i)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag

		if fieldType.PkgPath != "" {
			continue
		}

		validationTag := tag.Get("validate")
		if validationTag != "" {
			if err := v.validateField(field, fieldType.Name, validationTag); err != nil {
				if customError, ok := v.customErrors[Field(fieldType.Name)]["required"]; ok {
					if err.Error() == "field is required" {
						return &ValidationError{
							Field:   fieldType.Name,
							Message: ErrorMsg(customError),
						}
					}
				}

				if customError, ok := v.customErrors[Field(fieldType.Name)]["max"]; ok {
					if err.Error() == fmt.Sprintf("value exceeds maximum of %d", getValidationMaxValue(validationTag)) {
						return &ValidationError{
							Field:   fieldType.Name,
							Message: customError,
						}
					}
				}

				return err
			}
		}
	}

	return nil
}

func (v *Validator) validateField(field reflect.Value, fieldName string, validationTag string) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return &ValidationError{
				Field:   fieldName,
				Message: "field is required",
			}
		}
		field = field.Elem()
	}

	rules := parseValidationTag(validationTag)

	for _, rule := range rules {
		if rule == "required" && isZeroValue(field) {
			return &ValidationError{
				Field:   fieldName,
				Message: "field is required",
			}
		}

		if err := validateMaxMin(field, rule); err != nil {
			return err
		}

		if err := validateLen(field, rule); err != nil {
			return err
		}

		if err := validateEmail(field, rule); err != nil {
			return err
		}
	}

	return nil
}

func parseValidationTag(validationTag string) []string {
	return strings.Split(validationTag, ",")
}

func validateMaxMin(field reflect.Value, rule string) error {
	if strings.HasPrefix(rule, "max=") {
		max, err := strconv.Atoi(rule[len("max="):])
		if err == nil && field.Kind() == reflect.Int && field.Int() > int64(max) {
			return fmt.Errorf("value exceeds maximum of %d", max)
		} else if field.Kind() == reflect.String && len(field.String()) > max {
			return fmt.Errorf("length exceeds maximum of %d", max)
		}
	}

	if strings.HasPrefix(rule, "min=") {
		min, err := strconv.Atoi(rule[len("min="):])
		if err == nil && field.Kind() == reflect.Int && field.Int() < int64(min) {
			return fmt.Errorf("value is below minimum of %d", min)
		} else if field.Kind() == reflect.String && len(field.String()) < min {
			return fmt.Errorf("length is below minimum of %d", min)
		}
	}

	return nil
}

func validateLen(field reflect.Value, rule string) error {
	if strings.HasPrefix(rule, "len=") {
		expectedLen, err := strconv.Atoi(rule[len("len="):])
		if err == nil && field.Kind() == reflect.String && len(field.String()) != expectedLen {
			return fmt.Errorf("length must be exactly %d", expectedLen)
		}
	}

	return nil
}

func validateEmail(field reflect.Value, rule string) error {
	if rule == "email" && field.Kind() == reflect.String {
		email := field.String()
		if !isValidEmail(email) {
			return fmt.Errorf("invalid email format")
		}
	}
	return nil
}

func isZeroValue(field reflect.Value) bool {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	return (field.Kind() == reflect.String && field.String() == "") ||
		(field.Kind() == reflect.Int && field.Int() == 0) ||
		(field.Kind() == reflect.Slice && field.Len() == 0)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func getValidationMaxValue(validationTag string) int {
	if strings.HasPrefix(validationTag, "max=") {
		maxStr := validationTag[len("max="):]
		max, err := strconv.Atoi(maxStr)
		if err == nil {
			return max
		}
	}
	return 0
}
