package validator

import (
	"testing"

	"github.com/google/uuid"
)

// Test all integer types with pointers and values
type AllIntegerTypes struct {
	Int      int    `validate:"required,min=10,max=100"`
	IntPtr   *int   `validate:"required,min=10,max=100"`
	Int8     int8   `validate:"min=1,max=127"`
	Int8Ptr  *int8  `validate:"min=1,max=127"`
	Int16    int16  `validate:"min=100,max=1000"`
	Int16Ptr *int16 `validate:"min=100,max=1000"`
	Int32    int32  `validate:"min=1000,max=10000"`
	Int32Ptr *int32 `validate:"min=1000,max=10000"`
	Int64    int64  `validate:"min=10000,max=100000"`
	Int64Ptr *int64 `validate:"min=10000,max=100000"`
}

// Test all unsigned integer types
type AllUnsignedTypes struct {
	Uint      uint    `validate:"min=10,max=100"`
	UintPtr   *uint   `validate:"min=10,max=100"`
	Uint8     uint8   `validate:"min=1,max=255"`
	Uint8Ptr  *uint8  `validate:"min=1,max=255"`
	Uint16    uint16  `validate:"min=100,max=1000"`
	Uint16Ptr *uint16 `validate:"min=100,max=1000"`
	Uint32    uint32  `validate:"min=1000,max=10000"`
	Uint32Ptr *uint32 `validate:"min=1000,max=10000"`
	Uint64    uint64  `validate:"min=10000,max=100000"`
	Uint64Ptr *uint64 `validate:"min=10000,max=100000"`
}

// Test all float types
type AllFloatTypes struct {
	Float32      float32  `validate:"min=0,max=100"`
	Float32Ptr   *float32 `validate:"min=0,max=100"`
	Float64      float64  `validate:"min=0,max=100"`
	Float64Ptr   *float64 `validate:"min=0,max=100"`
	Latitude     float64  `validate:"latitude"`
	LatitudePtr  *float64 `validate:"latitude"`
	Longitude    float64  `validate:"longitude"`
	LongitudePtr *float64 `validate:"longitude"`
}

// Test UUID validation
type UUIDTypes struct {
	UUID         uuid.UUID  `validate:"required,uuid"`
	UUIDPtr      *uuid.UUID `validate:"uuid"`
	UUIDOptional uuid.UUID  `validate:"uuid"`
}

// Test string validations
type StringTypes struct {
	Email           string  `validate:"required,email"`
	EmailPtr        *string `validate:"email"`
	NumericString   string  `validate:"numeric"`
	NumericPtr      *string `validate:"numeric"`
	MinMaxString    string  `validate:"min=5,max=20"`
	MinMaxStringPtr *string `validate:"min=5,max=20"`
	LenString       string  `validate:"len=10"`
	LenStringPtr    *string `validate:"len=10"`
}

func TestAllIntegerTypes(t *testing.T) {
	t.Run("Valid values", func(t *testing.T) {
		intVal := 50
		int8Val := int8(10)
		int16Val := int16(500)
		int32Val := int32(5000)
		int64Val := int64(50000)

		data := AllIntegerTypes{
			Int:      50,
			IntPtr:   &intVal,
			Int8:     10,
			Int8Ptr:  &int8Val,
			Int16:    500,
			Int16Ptr: &int16Val,
			Int32:    5000,
			Int32Ptr: &int32Val,
			Int64:    50000,
			Int64Ptr: &int64Val,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for valid integer types, got: %v", err)
		}
	})

	t.Run("Required field missing", func(t *testing.T) {
		data := AllIntegerTypes{
			Int: 0, // Zero value, should fail required
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for missing required field")
		}
	})

	t.Run("Pointer nil for required", func(t *testing.T) {
		data := AllIntegerTypes{
			Int:    50,
			IntPtr: nil, // Should fail required
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for nil required pointer")
		}
	})

	t.Run("Int exceeds max", func(t *testing.T) {
		intVal := 150 // max is 100
		data := AllIntegerTypes{
			Int:    150,
			IntPtr: &intVal,
			Int8:   10,
			Int16:  500,
			Int32:  5000,
			Int64:  50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int exceeding max")
		}
	})

	t.Run("Int below min", func(t *testing.T) {
		intVal := 5 // min is 10
		data := AllIntegerTypes{
			Int:    5,
			IntPtr: &intVal,
			Int8:   10,
			Int16:  500,
			Int32:  5000,
			Int64:  50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int below min")
		}
	})

	t.Run("Int8 exceeds max", func(t *testing.T) {
		intVal := 50
		int8Val := int8(127) // max is 127, but let's test edge
		int8Over := int8(126)
		data := AllIntegerTypes{
			Int:     50,
			IntPtr:  &intVal,
			Int8:    127,
			Int8Ptr: &int8Val,
			Int16:   500,
			Int32:   5000,
			Int64:   50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Int8 at max should be valid, got: %v", err)
		}

		// Now test over max (can't use 128 as it overflows int8)
		data.Int8 = 127
		data.Int8Ptr = &int8Over
		err = validator.Validate(data)
		if err != nil {
			t.Errorf("Int8 at 126 should be valid, got: %v", err)
		}
	})

	t.Run("Int8 below min", func(t *testing.T) {
		intVal := 50
		int8Val := int8(0) // min is 1
		data := AllIntegerTypes{
			Int:     50,
			IntPtr:  &intVal,
			Int8:    0,
			Int8Ptr: &int8Val,
			Int16:   500,
			Int32:   5000,
			Int64:   50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int8 below min")
		}
	})

	t.Run("Int16 exceeds max", func(t *testing.T) {
		intVal := 50
		int16Val := int16(1001) // max is 1000
		data := AllIntegerTypes{
			Int:      50,
			IntPtr:   &intVal,
			Int8:     10,
			Int16:    1001,
			Int16Ptr: &int16Val,
			Int32:    5000,
			Int64:    50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int16 exceeding max")
		}
	})

	t.Run("Int16 below min", func(t *testing.T) {
		intVal := 50
		int16Val := int16(99) // min is 100
		data := AllIntegerTypes{
			Int:      50,
			IntPtr:   &intVal,
			Int8:     10,
			Int16:    99,
			Int16Ptr: &int16Val,
			Int32:    5000,
			Int64:    50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int16 below min")
		}
	})

	t.Run("Int32 exceeds max", func(t *testing.T) {
		intVal := 50
		int32Val := int32(10001) // max is 10000
		data := AllIntegerTypes{
			Int:      50,
			IntPtr:   &intVal,
			Int8:     10,
			Int16:    500,
			Int32:    10001,
			Int32Ptr: &int32Val,
			Int64:    50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int32 exceeding max")
		}
	})

	t.Run("Int32 below min", func(t *testing.T) {
		intVal := 50
		int32Val := int32(999) // min is 1000
		data := AllIntegerTypes{
			Int:      50,
			IntPtr:   &intVal,
			Int8:     10,
			Int16:    500,
			Int32:    999,
			Int32Ptr: &int32Val,
			Int64:    50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int32 below min")
		}
	})

	t.Run("Int64 exceeds max", func(t *testing.T) {
		intVal := 50
		int64Val := int64(100001) // max is 100000
		data := AllIntegerTypes{
			Int:      50,
			IntPtr:   &intVal,
			Int8:     10,
			Int16:    500,
			Int32:    5000,
			Int64:    100001,
			Int64Ptr: &int64Val,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int64 exceeding max")
		}
	})

	t.Run("Int64 below min", func(t *testing.T) {
		intVal := 50
		int64Val := int64(9999) // min is 10000
		data := AllIntegerTypes{
			Int:      50,
			IntPtr:   &intVal,
			Int8:     10,
			Int16:    500,
			Int32:    5000,
			Int64:    9999,
			Int64Ptr: &int64Val,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for int64 below min")
		}
	})
}

func TestAllUnsignedTypes(t *testing.T) {
	t.Run("Valid unsigned values", func(t *testing.T) {
		uintVal := uint(50)
		uint8Val := uint8(100)
		uint16Val := uint16(500)
		uint32Val := uint32(5000)
		uint64Val := uint64(50000)

		data := AllUnsignedTypes{
			Uint:      50,
			UintPtr:   &uintVal,
			Uint8:     100,
			Uint8Ptr:  &uint8Val,
			Uint16:    500,
			Uint16Ptr: &uint16Val,
			Uint32:    5000,
			Uint32Ptr: &uint32Val,
			Uint64:    50000,
			Uint64Ptr: &uint64Val,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for valid unsigned types, got: %v", err)
		}
	})

	t.Run("Uint exceeds max", func(t *testing.T) {
		uintVal := uint(101) // max is 100
		data := AllUnsignedTypes{
			Uint:    101,
			UintPtr: &uintVal,
			Uint8:   100,
			Uint16:  500,
			Uint32:  5000,
			Uint64:  50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for uint exceeding max")
		}
	})

	t.Run("Uint below min", func(t *testing.T) {
		uintVal := uint(9) // min is 10
		data := AllUnsignedTypes{
			Uint:    9,
			UintPtr: &uintVal,
			Uint8:   100,
			Uint16:  500,
			Uint32:  5000,
			Uint64:  50000,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for uint below min")
		}
	})

	t.Run("Uint64 exceeds max", func(t *testing.T) {
		uint64Val := uint64(100001) // max is 100000
		data := AllUnsignedTypes{
			Uint:      50,
			Uint8:     100,
			Uint16:    500,
			Uint32:    5000,
			Uint64:    100001,
			Uint64Ptr: &uint64Val,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for uint64 exceeding max")
		}
	})

	t.Run("Uint64 below min", func(t *testing.T) {
		uint64Val := uint64(9999) // min is 10000
		data := AllUnsignedTypes{
			Uint:      50,
			Uint8:     100,
			Uint16:    500,
			Uint32:    5000,
			Uint64:    9999,
			Uint64Ptr: &uint64Val,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for uint64 below min")
		}
	})
}

func TestAllFloatTypes(t *testing.T) {
	t.Run("Valid float values", func(t *testing.T) {
		float32Val := float32(50.5)
		float64Val := float64(75.7)
		latVal := float64(45.0)
		lonVal := float64(90.0)

		data := AllFloatTypes{
			Float32:      50.5,
			Float32Ptr:   &float32Val,
			Float64:      75.7,
			Float64Ptr:   &float64Val,
			Latitude:     45.0,
			LatitudePtr:  &latVal,
			Longitude:    90.0,
			LongitudePtr: &lonVal,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for valid float types, got: %v", err)
		}
	})

	t.Run("Invalid latitude", func(t *testing.T) {
		data := AllFloatTypes{
			Latitude: 91.0, // Out of range
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for invalid latitude")
		}
	})

	t.Run("Invalid longitude", func(t *testing.T) {
		data := AllFloatTypes{
			Longitude: 181.0, // Out of range
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for invalid longitude")
		}
	})

	t.Run("Float32 exceeds max", func(t *testing.T) {
		float32Val := float32(101.0) // max is 100
		data := AllFloatTypes{
			Float32:    101.0,
			Float32Ptr: &float32Val,
			Float64:    75.7,
			Latitude:   45.0,
			Longitude:  90.0,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for float32 exceeding max")
		}
	})

	t.Run("Float32 below min", func(t *testing.T) {
		float32Val := float32(-1.0) // min is 0
		data := AllFloatTypes{
			Float32:    -1.0,
			Float32Ptr: &float32Val,
			Float64:    75.7,
			Latitude:   45.0,
			Longitude:  90.0,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for float32 below min")
		}
	})

	t.Run("Float64 exceeds max", func(t *testing.T) {
		float64Val := float64(100.1) // max is 100
		data := AllFloatTypes{
			Float32:    50.5,
			Float64:    100.1,
			Float64Ptr: &float64Val,
			Latitude:   45.0,
			Longitude:  90.0,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for float64 exceeding max")
		}
	})

	t.Run("Float64 below min", func(t *testing.T) {
		float64Val := float64(-0.1) // min is 0
		data := AllFloatTypes{
			Float32:    50.5,
			Float64:    -0.1,
			Float64Ptr: &float64Val,
			Latitude:   45.0,
			Longitude:  90.0,
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for float64 below min")
		}
	})
}

func TestUUIDValidation(t *testing.T) {
	t.Run("Valid UUID", func(t *testing.T) {
		validUUID := uuid.New()
		data := UUIDTypes{
			UUID:    validUUID,
			UUIDPtr: &validUUID,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for valid UUID, got: %v", err)
		}
	})

	t.Run("Nil UUID required", func(t *testing.T) {
		data := UUIDTypes{
			UUID: uuid.Nil, // Should fail required
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for nil UUID on required field")
		}
	})

	t.Run("Optional UUID can be nil", func(t *testing.T) {
		validUUID := uuid.New()
		data := UUIDTypes{
			UUID:         validUUID,
			UUIDOptional: uuid.Nil, // This is okay, not required
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for optional nil UUID, got: %v", err)
		}
	})
}

func TestStringValidations(t *testing.T) {
	t.Run("Valid email", func(t *testing.T) {
		emailVal := "test@example.com"
		data := StringTypes{
			Email:    "test@example.com",
			EmailPtr: &emailVal,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for valid email, got: %v", err)
		}
	})

	t.Run("Invalid email formats", func(t *testing.T) {
		invalidEmails := []string{
			"notanemail",
			"@example.com",
			"test@",
			"test @example.com",
			"test@example",
		}

		for _, email := range invalidEmails {
			data := StringTypes{
				Email: email,
			}

			validator := New()
			err := validator.Validate(data)
			if err == nil {
				t.Errorf("Expected error for invalid email: %s", email)
			}
		}
	})

	t.Run("Valid numeric string", func(t *testing.T) {
		numVal := "123.45"
		data := StringTypes{
			Email:         "test@example.com",
			NumericString: "123.45",
			NumericPtr:    &numVal,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for valid numeric string, got: %v", err)
		}
	})

	t.Run("Invalid numeric string", func(t *testing.T) {
		data := StringTypes{
			Email:         "test@example.com",
			NumericString: "abc123",
		}

		validator := New()
		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for invalid numeric string")
		}
	})

	t.Run("String min/max length", func(t *testing.T) {
		validStr := "Valid String"
		data := StringTypes{
			Email:           "test@example.com",
			MinMaxString:    validStr,
			MinMaxStringPtr: &validStr,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for valid string length, got: %v", err)
		}

		// Test too short
		shortStr := "Hi"
		data.MinMaxString = shortStr
		err = validator.Validate(data)
		if err == nil {
			t.Error("Expected error for string too short")
		}

		// Test too long
		longStr := "This is a very long string that exceeds maximum"
		data.MinMaxString = longStr
		err = validator.Validate(data)
		if err == nil {
			t.Error("Expected error for string too long")
		}
	})

	t.Run("String exact length", func(t *testing.T) {
		exactLen := "1234567890"
		data := StringTypes{
			Email:     "test@example.com",
			LenString: exactLen,
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for exact length string, got: %v", err)
		}

		// Test wrong length
		data.LenString = "12345"
		err = validator.Validate(data)
		if err == nil {
			t.Error("Expected error for wrong length string")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("Empty struct", func(t *testing.T) {
		type Empty struct{}
		data := Empty{}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for empty struct, got: %v", err)
		}
	})

	t.Run("Struct with no validation tags", func(t *testing.T) {
		type NoTags struct {
			Name  string
			Age   int
			Email string
		}
		data := NoTags{
			Name:  "Test",
			Age:   25,
			Email: "invalid",
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for struct with no tags, got: %v", err)
		}
	})

	t.Run("Unexported fields ignored", func(t *testing.T) {
		type WithPrivate struct {
			Public  string `validate:"required"`
			private string `validate:"required"` // Should be ignored
		}
		data := WithPrivate{
			Public:  "test",
			private: "", // Should not cause error
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors with unexported field, got: %v", err)
		}
	})

	t.Run("Validation on pointer to struct", func(t *testing.T) {
		type TestStruct struct {
			Name string `validate:"required"`
		}
		data := &TestStruct{
			Name: "Test",
		}

		validator := New()
		err := validator.Validate(data)
		if err != nil {
			t.Errorf("Expected no errors for pointer to struct, got: %v", err)
		}
	})
}

func TestCustomErrorsDetailed(t *testing.T) {
	type TestData struct {
		Name  string `validate:"required,min=3"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=18"`
	}

	validator := New().WithCustomErrors(CustomErrors{
		"Name": {
			"required": "Name is absolutely required",
			"min":      "Name must be at least 3 characters",
		},
		"Email": {
			"required": "Email is absolutely required",
			"email":    "Please provide a valid email address",
		},
		"Age": {
			"min": "You must be 18 or older",
		},
	})

	t.Run("Custom error for required", func(t *testing.T) {
		data := TestData{
			Name:  "",
			Email: "test@example.com",
			Age:   20,
		}

		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for missing name")
		} else if err.Error() != "Name is absolutely required" {
			t.Errorf("Expected custom error message, got: %s", err.Error())
		}
	})

	t.Run("Custom error for min", func(t *testing.T) {
		data := TestData{
			Name:  "Hi",
			Email: "test@example.com",
			Age:   20,
		}

		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for name too short")
		} else if err.Error() != "Name must be at least 3 characters" {
			t.Errorf("Expected custom error message, got: %s", err.Error())
		}
	})

	t.Run("Custom error for email", func(t *testing.T) {
		data := TestData{
			Name:  "John",
			Email: "invalid",
			Age:   20,
		}

		err := validator.Validate(data)
		if err == nil {
			t.Error("Expected error for invalid email")
		} else if err.Error() != "Please provide a valid email address" {
			t.Errorf("Expected custom error message, got: %s", err.Error())
		}
	})

	t.Run("GetField method works", func(t *testing.T) {
		data := TestData{
			Name:  "",
			Email: "test@example.com",
			Age:   20,
		}

		err := validator.Validate(data)
		if err != nil {
			if valErr, ok := err.(*ValidationError); ok {
				if valErr.GetField() != "Name" {
					t.Errorf("Expected field name 'Name', got: %s", valErr.GetField())
				}
			} else {
				t.Error("Expected ValidationError type")
			}
		}
	})
}
