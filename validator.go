package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Field string
type Rule string
type ErrorMsg string

type CustomErrors map[Field]map[Rule]ErrorMsg

type ValidationError struct {
	// Internal field name - should never be exposed to API responses
	field   string
	Message ErrorMsg
}

func (e *ValidationError) Error() string {
	// Return only the message, never expose field names
	return string(e.Message)
}

// GetField returns the internal field name for logging purposes only
// NEVER expose this to API responses
func (e *ValidationError) GetField() string {
	return e.field
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

func (v *Validator) Validate(i any) error {
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
				return err
			}
		}
	}

	return nil
}

func (v *Validator) validateField(field reflect.Value, fieldName string, validationTag string) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			if strings.Contains(validationTag, "required") {
				return v.getCustomValidationError(fieldName, "required", "Invalid request data")
			} else {
				return nil
			}
		}
		field = field.Elem()
	}

	rules := parseValidationTag(validationTag)

	// Check if field is required
	isRequired := false
	for _, rule := range rules {
		if rule == "required" {
			isRequired = true
			if isZeroValue(field) {
				return v.getCustomValidationError(fieldName, "required", "Invalid request data")
			}
			break
		}
	}

	// For optional fields: skip validation only if the value is semantically "not provided"
	// - For strings: empty string means not provided
	// - For other types: they're always considered "provided" even if zero value
	if !isRequired && field.Kind() == reflect.String && field.String() == "" {
		return nil
	}

	// Validate all other rules
	for _, rule := range rules {
		if rule == "required" {
			continue // Already handled
		}

		if err := validateMaxMin(field, rule); err != nil {
			return v.getCustomValidationError(fieldName, extractRuleType(rule), "Invalid request data")
		}

		if err := validateLen(field, rule); err != nil {
			return v.getCustomValidationError(fieldName, extractRuleType(rule), "Invalid request data")
		}

		if err := validateEmail(field, rule); err != nil {
			return v.getCustomValidationError(fieldName, extractRuleType(rule), "Invalid email format")
		}

		if err := validateUUID(field, rule); err != nil {
			return v.getCustomValidationError(fieldName, extractRuleType(rule), "Invalid identifier format")
		}

		if err := validateLatitude(field, rule); err != nil {
			return v.getCustomValidationError(fieldName, extractRuleType(rule), "Invalid latitude value")
		}

		if err := validateLongitude(field, rule); err != nil {
			return v.getCustomValidationError(fieldName, extractRuleType(rule), "Invalid longitude value")
		}

		if err := validateNumeric(field, rule); err != nil {
			return v.getCustomValidationError(fieldName, extractRuleType(rule), "Invalid numeric format")
		}
	}

	return nil
}

func extractRuleType(rule string) string {
	if strings.Contains(rule, "=") {
		return strings.Split(rule, "=")[0]
	}
	return rule
}

func (v *Validator) getCustomValidationError(fieldName, rule, defaultMsg string) error {
	if customError, ok := v.customErrors[Field(fieldName)][Rule(rule)]; ok {
		return &ValidationError{
			field:   fieldName,
			Message: customError,
		}
	}

	return &ValidationError{
		field:   fieldName,
		Message: ErrorMsg(defaultMsg),
	}
}

func parseValidationTag(validationTag string) []string {
	return strings.Split(validationTag, ",")
}

func validateMaxMin(field reflect.Value, rule string) error {
	// Handle pointers
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil
		}
		field = field.Elem()
	}

	if strings.HasPrefix(rule, "max=") {
		maxStr := rule[len("max="):]
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			max, err := strconv.ParseInt(maxStr, 10, 64)
			if err != nil {
				return err
			}
			if field.Int() > max {
				return fmt.Errorf("value exceeds maximum of %d", max)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			max, err := strconv.ParseUint(maxStr, 10, 64)
			if err != nil {
				return err
			}
			if field.Uint() > max {
				return fmt.Errorf("value exceeds maximum of %d", max)
			}
		case reflect.Float32, reflect.Float64:
			max, err := strconv.ParseFloat(maxStr, 64)
			if err != nil {
				return err
			}
			if field.Float() > max {
				return fmt.Errorf("value exceeds maximum of %f", max)
			}
		case reflect.String:
			max, err := strconv.Atoi(maxStr)
			if err != nil {
				return err
			}
			if utf8.RuneCountInString(field.String()) > max {
				return fmt.Errorf("length exceeds maximum of %d", max)
			}
		}
	}

	if strings.HasPrefix(rule, "min=") {
		minStr := rule[len("min="):]
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			min, err := strconv.ParseInt(minStr, 10, 64)
			if err != nil {
				return err
			}
			if field.Int() < min {
				return fmt.Errorf("value is below minimum of %d", min)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			min, err := strconv.ParseUint(minStr, 10, 64)
			if err != nil {
				return err
			}
			if field.Uint() < min {
				return fmt.Errorf("value is below minimum of %d", min)
			}
		case reflect.Float32, reflect.Float64:
			min, err := strconv.ParseFloat(minStr, 64)
			if err != nil {
				return err
			}
			if field.Float() < min {
				return fmt.Errorf("value is below minimum of %f", min)
			}
		case reflect.String:
			min, err := strconv.Atoi(minStr)
			if err != nil {
				return err
			}
			if utf8.RuneCountInString(field.String()) < min {
				return fmt.Errorf("length is below minimum of %d", min)
			}
		}
	}

	return nil
}

func validateLen(field reflect.Value, rule string) error {
	if strings.HasPrefix(rule, "len=") {
		expectedLen, err := strconv.Atoi(rule[len("len="):])
		if err == nil && field.Kind() == reflect.String && utf8.RuneCountInString(field.String()) != expectedLen {
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

	// Handle UUID type
	if field.Type() == reflect.TypeOf(uuid.UUID{}) {
		uuidVal := field.Interface().(uuid.UUID)
		return uuidVal == uuid.Nil
	}

	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0.0
	case reflect.Bool:
		return !field.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len() == 0
	default:
		return false
	}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func validateUUID(field reflect.Value, rule string) error {
	if rule == "uuid" || rule == "required" {
		// Check if the field type is uuid.UUID
		if field.Type() == reflect.TypeOf(uuid.UUID{}) {
			uuidVal := field.Interface().(uuid.UUID)
			// uuid.Nil is the zero value for UUID
			if uuidVal == uuid.Nil && rule == "required" {
				return fmt.Errorf("uuid is required")
			}
		}
	}
	return nil
}

func validateLatitude(field reflect.Value, rule string) error {
	if rule == "latitude" && field.Kind() == reflect.Float64 {
		lat := field.Float()
		if lat < -90.0 || lat > 90.0 {
			return fmt.Errorf("latitude out of range")
		}
	}
	return nil
}

func validateLongitude(field reflect.Value, rule string) error {
	if rule == "longitude" && field.Kind() == reflect.Float64 {
		lon := field.Float()
		if lon < -180.0 || lon > 180.0 {
			return fmt.Errorf("longitude out of range")
		}
	}
	return nil
}

func validateNumeric(field reflect.Value, rule string) error {
	if rule == "numeric" && field.Kind() == reflect.String {
		str := field.String()
		if str == "" {
			return nil // Empty strings are handled by required rule
		}
		// Check if string contains only digits and dots
		for _, char := range str {
			if !((char >= '0' && char <= '9') || char == '.') {
				return fmt.Errorf("string is not numeric")
			}
		}
	}
	return nil
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
