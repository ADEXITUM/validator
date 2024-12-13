### Validator Package Documentation

The `validator` package provides utilities for struct field validation in Go. It supports custom error messages, common validation rules (e.g., required, max, min, length, email), and pointer field handling (dereferencing non-nil pointers).

---

#### Types:
**CustomErrors**  
   A map containing custom error messages for specific fields and validation rules.

**ValidationError**  
   Struct representing a validation error for a field.

---

#### Functions:

1. **New()**  
   Creates and returns a new `Validator` instance.

   ```go
   v := New()
   ```

2. **WithCustomErrors(errors CustomErrors) *Validator**  
   Sets custom error messages for specific fields and rules.

   ```go
   v.WithCustomErrors(CustomErrors{
     "FieldName": {
       "required": "This field is required",
     },
   })
   ```

3. **Validate(i interface{}) error**  
   Validates the fields of a struct passed as the interface `i`. Returns an error if any validation fails.

   ```go
   err := v.Validate(&myStruct)
   ```
---

### Important Notes:
- **Pointer Fields**: If a struct field is a pointer, and it is not `nil`, the field will be dereferenced for validation. For example, if a pointer to an integer is provided, it is dereferenced to check its value.
- **Validation Tags**: Fields can have validation rules defined in their struct tags (e.g., `validate:"required,max=10"`). The package processes these tags and applies the corresponding validations.
- **Custom Error Messages**: You can define custom error messages for specific rules and fields using the `WithCustomErrors` method. This overrides default error messages for specific cases.

---

### Example Usage:

```go
package main

import (
	"fmt"
	"github.com/ADEXITUM/validator"
)

type User struct {
	Name  *string `validate:"required"`
	Age   int     `validate:"min=18,max=100"`
	Email string  `validate:"email"`
}

func main() {
	v := validator.New()

	// Custom error messages
	v.WithCustomErrors(validator.CustomErrors{
		"Name": {
			"required": "Name is mandatory",
		},
		"Age": {
			"min": "Age must be at least 18",
		},
	})

	// Example struct with invalid data
	user := &User{
		Age:   16,
		Email: "invalid-email",
	}

	// Validate the struct
	err := v.Validate(user)
	if err != nil {
		fmt.Println("Validation Error:", err)
	}
}
```