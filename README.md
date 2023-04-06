# Validation Package

[![GoDoc](https://godoc.org/github.com/JubaerHossain/validation?status.svg)](https://godoc.org/github.com/JubaerHossain/validation)
[![Go Report Card](https://goreportcard.com/badge/github.com/JubaerHossain/validation)](https://goreportcard.com/report/github.com/JubaerHossain/validation)
[![Build Status](https://travis-ci.org/JubaerHossain/validation.svg?branch=master)](https://travis-ci.org/JubaerHossain/validation)
[![codecov](https://codecov.io/gh/JubaerHossain/validation/branch/master/graph/badge.svg)](https://codecov.io/gh/JubaerHossain/validation)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/JubaerHossain/validation.svg?style=flat-square)](https://github.com/JubaerHossain/validation/releases/latest)


This is a simple Go package for validating structs based on a set of validation rules. It can be used to ensure that incoming data from clients, databases or other sources meets certain criteria before being processed further.

# Installation

To use this package in your Go project, you can install it using the go get command:

```golang
go get github.com/JubaerHossain/validation
```

# Getting Started

Here's an example of how to use this package to validate a User struct:

```golang
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
```

In this example, we define a User struct with two fields: Name and Password. We then define an array of ValidationRule structs that describe the validation rules for each field. Each ValidationRule has a Field name, a human-readable Description of the field, and an array of validation functions that take a value and return a ValidationErrorItem if the value is invalid.

We then call the Validate function with the User struct and the array of validation rules. This function returns an array of ValidationErrorItems if there are any validation errors, or an empty array if the input is valid.

# Available Validations

Here are the available validation functions that can be used in a ValidationRule:

- **RequiredValidation** - Returns an error if the value is nil or an empty string.
- **MinLengthValidation(min int)** - Returns an error if the string length is less than the minimum value.
- **MaxLengthValidation(max int)** - Returns an error if the string length is greater than the maximum value.
- You can also define your own custom validation functions by implementing the func(interface{}) ValidationErrorItem signature.

# License

This package is licensed under the MIT License. See the LICENSE file for details.
