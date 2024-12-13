package validator

import (
	"testing"
)

type User struct {
	Name    *string `validate:"required,min=3,max=50"`
	Email   string  `validate:"required,email"`
	Age     int     `validate:"min=18,max=100"`
	Address string  `validate:"len=10"`
}

func TestValidator(t *testing.T) {
	var name string = "John"
	user := User{
		Name:    &name,
		Email:   "invalidemailcom",
		Age:     17,
		Address: "Short",
	}

	validator := New().WithCustomErrors(CustomErrors{
		"Email": {
			"required": "Email is required",
			"email":    "Please provide a valid email",
		},
		"Age": {
			"min": "You must be at least 18 years old",
			"max": "Age cannot exceed 100",
		},
		"Address": {
			"len": "Address must be exactly 10 characters",
		},
	})

	// Test for various validation failures with custom error messages
	err := validator.Validate(user)
	if err == nil {
		t.Errorf("Expected validation errors, but got none")
	} else {
		t.Log("Validation Error:", err)
	}

	// Modify user to pass validation
	name = "John Doe"
	user.Name = &name
	user.Email = "john.doe@example.com"
	user.Age = 25
	user.Address = "1234567890"

	err = validator.Validate(user)
	if err != nil {
		t.Errorf("Expected no validation errors, but got: %s", err)
	} else {
		t.Log("Validation passed!")
	}
}

func TestPointerValidation(t *testing.T) {
	var name string = "Valid Name"
	user := User{
		Name:    &name,
		Email:   "valid@example.com",
		Age:     20,
		Address: "1234567890",
	}

	// Test: Validate when pointer is nil (should trigger "required" error)
	user.Name = nil
	validator := New().WithCustomErrors(CustomErrors{
		"Name": {
			"required": "Name is required",
		},
	})

	err := validator.Validate(user)
	if err == nil {
		t.Errorf("Expected 'Name is required' error, but got none")
	} else {
		t.Log("Validation Error (Pointer nil):", err)
	}

	// Test: Validate when pointer is not nil but violates other rules
	name = "A"
	user.Name = &name // violates min=3 rule
	err = validator.Validate(user)
	if err == nil {
		t.Errorf("Expected 'Name is too short' error, but got none")
	} else {
		t.Log("Validation Error (Pointer not nil but validation failed):", err)
	}

	// Test: Validate when pointer is not nil and satisfies all rules
	name = "Valid Name"
	user.Name = &name
	err = validator.Validate(user)
	if err != nil {
		t.Errorf("Expected no validation errors, but got: %s", err)
	} else {
		t.Log("Validation passed (Pointer valid)!")
	}
}

func TestCustomErrors(t *testing.T) {
	var name string = "Short"
	user := User{
		Name:    &name,
		Email:   "invalidemailcom",
		Age:     17,
		Address: "Short",
	}

	validator := New().WithCustomErrors(CustomErrors{
		"Email": {
			"required": "Email is required",
			"email":    "Please provide a valid email",
		},
		"Age": {
			"min": "You must be at least 18 years old",
			"max": "Age cannot exceed 100",
		},
		"Address": {
			"len": "Address must be exactly 10 characters",
		},
	})

	// Test: Custom error message for required field
	err := validator.Validate(user)
	if err == nil {
		t.Errorf("Expected validation errors, but got none")
	} else {
		t.Log("Custom Validation Error:", err)
	}

	// Modify user to pass validation
	name = "John Doe"
	user.Name = &name
	user.Email = "valid@example.com"
	user.Age = 25
	user.Address = "1234567890"

	// Test: Custom error message for passing validation
	err = validator.Validate(user)
	if err != nil {
		t.Errorf("Expected no validation errors, but got: %s", err)
	} else {
		t.Log("Validation passed!")
	}
}

func TestMaxMinValidation(t *testing.T) {
	var name string = "John Doe"
	user := User{
		Name:    &name,
		Email:   "valid@example.com",
		Age:     101, // exceeds max=100
		Address: "1234567890",
	}

	validator := New().WithCustomErrors(CustomErrors{
		"Age": {
			"max": "Age cannot exceed 100",
		},
	})

	err := validator.Validate(user)
	if err == nil {
		t.Errorf("Expected 'Age cannot exceed 100' error, but got none")
	} else {
		t.Log("Validation Error (Age exceeds max):", err)
	}

	user.Age = 99 // within the valid range
	err = validator.Validate(user)
	if err != nil {
		t.Errorf("Expected no validation errors, but got: %s", err)
	} else {
		t.Log("Validation passed (Age valid)!")
	}
}
