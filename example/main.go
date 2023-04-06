package main

import (
	"fmt"

	"github.com/JubaerHossain/validation"
)

type User struct {
	Name     string `json:"first_name"`
	Password string `json:"password"`
}

func main() {

	userInput := User{
		Name:     "John",
		Password: "123456",
	}
	validationErrors := ValidateUser(userInput)
	if validationErrors != nil {
		var errorMsgs []string
		for _, validationErr := range validationErrors {
			errorMsgs = append(errorMsgs, validationErr.Field+" : "+validationErr.Message)
		}
		fmt.Println(errorMsgs)
	}
}

func ValidateUser(user User) []validation.ValidationErrorItem {

	rules := []validation.ValidationRule{
		{
			Field:       "name",
			Description: "Name",
			Validations: []func(interface{}) validation.ValidationErrorItem{
				validation.RequiredValidation,
				validation.MinLengthValidation(3),
				validation.MaxLengthValidation(50),
			},
		},
		{
			Field:       "password",
			Description: "Password",
			Validations: []func(interface{}) validation.ValidationErrorItem{
				validation.RequiredValidation,
				validation.MinLengthValidation(6),
				validation.MaxLengthValidation(50),
			},
		},
	}

	return validation.Validate(user, rules)

}
