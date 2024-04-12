package utils

import "reflect"

// JsonDataHasKey checks if the given data has a field with the specified key.
// It returns the value of the field and a boolean indicating whether the field exists.
func JsonDataHasKey(data interface{}, key string) (interface{}, bool) {
	value := reflect.ValueOf(data)

	// Check if the data is a struct
	if value.Kind() != reflect.Struct {
		return "null", false
	}

	// Get the field by name
	fieldValue := value.FieldByName(key)

	// Check if the field exists
	if !fieldValue.IsValid() {
		return "null", false
	}

	// Return the field value
	return fieldValue.Interface(), true
}
