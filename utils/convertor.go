package utils

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertUint16(str string) uint16 {
	u, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(u)
}

func ConvertUint8(str string) uint8 {
	u, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(u)
}

// CopyAndSetField creates a copy of the provided struct and sets the specified field to the new value.
func CopyAndSetField[T any, V any](input *T, fieldName string, newValue *V) *T {
	// Get the reflect.Value of the input struct
	inputValue := reflect.ValueOf(*input)

	// Create a new struct of the same type as the input
	outputValue := reflect.New(inputValue.Type()).Elem()

	// Copy the field values from the input struct to the new struct
	outputValue.Set(inputValue)

	// Get the reflect.Value of the field by name
	fieldValue := outputValue.FieldByName(fieldName)

	// Check if the field exists and is settable
	if fieldValue.IsValid() && fieldValue.CanSet() {
		// Convert the new value to the field's type
		newValueReflect := reflect.ValueOf(*newValue).Convert(fieldValue.Type())

		// Set the value of the field in the new struct
		fieldValue.Set(newValueReflect)
	}

	// Return the new struct as an interface{}
	output := outputValue.Interface().(T)

	return &output
}

func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

func MapTo[T any, V any](from *V) (*T, error) {
	var to T
	err := copier.Copy(&to, from)
	if err != nil {
		return nil, err
	}
	return &to, nil
}

func ExtractBearerToken(authHeader string) string {
	const prefix = "Bearer "
	tokenIndex := strings.Index(authHeader, prefix)
	if tokenIndex == -1 || tokenIndex != 0 {
		return ""
	}
	return authHeader[tokenIndex+len(prefix):]
}

func FormatEndpoint(endpoint string) string {
	endpoint = strings.ReplaceAll(endpoint, " ", "")
	endpoint = strings.ReplaceAll(endpoint, "/", "-")
	endpoint = strings.ReplaceAll(endpoint, "?", "")
	return endpoint
}
