package validation

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type ValidationFunc func(interface{}) ValidationErrorItem

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationRule struct {
	Field       string
	Description string
	Validations []func(interface{}) ValidationErrorItem
}

func Validate(data interface{}, rules []ValidationRule) []ValidationErrorItem {
	var validationErrs []ValidationErrorItem

	for _, rule := range rules {
		fieldValue, err := getField(data, rule.Field)
		if err != nil {
			validationErrs = append(validationErrs, ValidationErrorItem{
				Field:   rule.Field,
				Message: err.Error(),
			})
			continue
		}

		for _, validation := range rule.Validations {
			validationErr := validation(fieldValue)
			if validationErr != (ValidationErrorItem{}) {
				validationErr.Field = rule.Field // set the field name
				validationErrs = append(validationErrs, validationErr)
			}
		}
	}

	return validationErrs
}

func getField(data interface{}, field string) (interface{}, error) {
	value := getFieldByJsonTag(data, field)
	if value == nil {
		return nil, fmt.Errorf("field not found: %s", field)
	}

	return value, nil
}

func getFieldByJsonTag(data interface{}, field string) interface{} {
	val := reflect.ValueOf(data)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		if tag == field {
			return val.Field(i).Interface()
		}
	}

	return nil
}

func RequiredValidation(value interface{}) ValidationErrorItem {
	if value == nil || value == "" {
		return ValidationErrorItem{
			Message: "Field is required",
		}
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		if value.(string) == "" {
			return ValidationErrorItem{
				Message: "Field is required",
			}
		}
	case reflect.Slice, reflect.Map:
		if reflect.ValueOf(value).Len() == 0 {
			return ValidationErrorItem{
				Message: "Field is required",
			}
		}
	}

	return ValidationErrorItem{}
}

func MinLengthValidation(minLength int) func(interface{}) ValidationErrorItem {
	return func(value interface{}) ValidationErrorItem {
		if value == nil || value == "" {
			return ValidationErrorItem{}
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			if len(value.(string)) < minLength {
				return ValidationErrorItem{
					Message: fmt.Sprintf("Field must be at least %d characters long", minLength),
				}
			}
		case reflect.Slice, reflect.Map:
			if reflect.ValueOf(value).Len() < minLength {
				return ValidationErrorItem{
					Message: fmt.Sprintf("Field must have at least %d items", minLength),
				}
			}
		}

		return ValidationErrorItem{}
	}
}

func MaxLengthValidation(maxLength int) func(interface{}) ValidationErrorItem {
	return func(value interface{}) ValidationErrorItem {
		if value == nil || value == "" {
			return ValidationErrorItem{}
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			if len(value.(string)) > maxLength {
				return ValidationErrorItem{
					Message: fmt.Sprintf("Field must not exceed %d characters", maxLength),
				}
			}
		case reflect.Slice, reflect.Map:
			if reflect.ValueOf(value).Len() > maxLength {
				return ValidationErrorItem{
					Message: fmt.Sprintf("Field must not have more than %d items", maxLength),
				}
			}
		}

		return ValidationErrorItem{}
	}
}

func EmailValidation(value interface{}) ValidationErrorItem {
	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	email, ok := value.(string)
	if !ok {
		return ValidationErrorItem{}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ValidationErrorItem{
			Message: "Field must be a valid email address",
		}
	}

	return ValidationErrorItem{}
}

func PhoneValidation(value interface{}) ValidationErrorItem {
	if value == nil || value == "" {
		return ValidationErrorItem{}
	}
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(value.(string)) {
		return ValidationErrorItem{
			Message: "Invalid phone number format",
		}
	}

	return ValidationErrorItem{}
}

func URLValidation(value interface{}) ValidationErrorItem {
	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	urlStr, ok := value.(string)
	if !ok {
		return ValidationErrorItem{
			Message: "Value is not a string",
		}
	}

	// Regular expression pattern for matching URLs
	pattern := `^(http|https)://[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`
	match, err := regexp.MatchString(pattern, urlStr)
	if err != nil {
		return ValidationErrorItem{
			Message: fmt.Sprintf("Error while matching URL: %v", err),
		}
	}

	if !match {
		return ValidationErrorItem{
			Message: "Value is not a valid URL",
		}
	}

	return ValidationErrorItem{}
}

func StringValidation(value interface{}) ValidationErrorItem {
	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		// do nothing, value is already a string
	default:
		return ValidationErrorItem{
			Message: "Field must be a string",
		}
	}

	return ValidationErrorItem{}
}

func NumericValidation(value interface{}) ValidationErrorItem {
	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		match, _ := regexp.MatchString(`^[0-9]+$`, value.(string))
		if !match {
			return ValidationErrorItem{
				Message: "Field must be numeric",
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// do nothing, value is already numeric
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// do nothing, value is already numeric
	default:
		return ValidationErrorItem{
			Message: "Field must be numeric",
		}
	}

	return ValidationErrorItem{}
}

func DateValidation(value interface{}) ValidationErrorItem {

	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	dateStr, ok := value.(string)
	if !ok {
		return ValidationErrorItem{
			Message: "Invalid date format",
		}
	}

	match, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr)
	if err != nil {
		return ValidationErrorItem{
			Message: "Invalid date format",
		}
	}
	if !match {
		return ValidationErrorItem{
			Message: "Invalid date format",
		}
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ValidationErrorItem{
			Message: "Invalid date",
		}
	}

	// Check that the year, month, and day are valid
	if date.Year() < 1 || date.Year() > 9999 {
		return ValidationErrorItem{
			Message: "Invalid year",
		}
	}
	if date.Month() < 1 || date.Month() > 12 {
		return ValidationErrorItem{
			Message: "Invalid month",
		}
	}
	if date.Day() < 1 || date.Day() > 31 {
		return ValidationErrorItem{
			Message: "Invalid day",
		}
	}

	return ValidationErrorItem{}
}

func ImageValidation(value interface{}) ValidationErrorItem {

	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return ValidationErrorItem{
			Message: "Invalid file format",
		}
	}

	// Check file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return ValidationErrorItem{
			Message: "File must be an image",
		}
	}

	return ValidationErrorItem{}
}

func FileSizeValidation(maxSize int64) func(interface{}) ValidationErrorItem {
	return func(value interface{}) ValidationErrorItem {

		if value == nil || value == "" {
			return ValidationErrorItem{}
		}

		file, ok := value.(*multipart.FileHeader)
		if !ok {
			return ValidationErrorItem{
				Message: "Invalid file format",
			}
		}

		if file.Size > maxSize {
			return ValidationErrorItem{
				Message: fmt.Sprintf("File size must be less than %d bytes", maxSize),
			}
		}

		return ValidationErrorItem{}
	}
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func ImageMimeValidation(value interface{}) ValidationErrorItem {

	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return ValidationErrorItem{
			Message: "Invalid file format",
		}
	}

	allowedMimeTypes := []string{"image/png", "image/jpg", "image/jpeg", "image/svg+xml"}

	// Check file type
	if !contains(allowedMimeTypes, file.Header.Get("Content-Type")) {
		return ValidationErrorItem{
			Message: "File must be a PNG, JPG, JPEG or SVG",
		}
	}

	return ValidationErrorItem{}
}

func FileValidation(value interface{}) ValidationErrorItem {

	if value == nil || value == "" {
		return ValidationErrorItem{}
	}
	
	_, ok := value.(*multipart.FileHeader)
	if !ok {
		return ValidationErrorItem{
			Message: "Invalid file format",
		}
	}

	return ValidationErrorItem{}
}

func FileTypeValidation(value interface{}, validTypes []string) ValidationErrorItem {

	if value == nil || value == "" {
		return ValidationErrorItem{}
	}

	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return ValidationErrorItem{
			Message: "Invalid file format",
		}
	}

	// Check file type
	fileType := strings.ToLower(strings.TrimPrefix(file.Header.Get("Content-Type"), "application/"))
	if !contains(validTypes, fileType) {
		return ValidationErrorItem{
			Message: "File type is not allowed",
		}
	}

	return ValidationErrorItem{}
}
