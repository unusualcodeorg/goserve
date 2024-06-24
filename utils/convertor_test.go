package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestConvertUint16(t *testing.T) {
	tests := []struct {
		input    string
		expected uint16
	}{
		{"65535", 65535},
		{"0", 0},
		{"12345", 12345},
		{"invalid", 0},
	}

	for _, tt := range tests {
		result := ConvertUint16(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

func TestConvertUint8(t *testing.T) {
	tests := []struct {
		input    string
		expected uint8
	}{
		{"255", 255},
		{"0", 0},
		{"123", 123},
		{"invalid", 0},
	}

	for _, tt := range tests {
		result := ConvertUint8(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestCopyAndSetField(t *testing.T) {
	input := &TestStruct{Field1: "original", Field2: 123}
	newValue := "new value"
	result := CopyAndSetField(input, "Field1", &newValue)

	assert.Equal(t, "new value", result.Field1)
	assert.Equal(t, 123, result.Field2)
}

func TestIsValidObjectID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()
	invalidID := "invalid"

	assert.True(t, IsValidObjectID(validID))
	assert.False(t, IsValidObjectID(invalidID))
}

func TestMapTo(t *testing.T) {
	type From struct {
		Field1 string
		Field2 int
	}

	type To struct {
		Field1 string
		Field2 int
	}

	from := &From{Field1: "value", Field2: 42}
	to, err := MapTo[To](from)

	assert.NoError(t, err)
	assert.Equal(t, from.Field1, to.Field1)
	assert.Equal(t, from.Field2, to.Field2)
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Bearer token123", "token123"},
		{"Bearer ", ""},
		{"Invalid token123", ""},
		{"BearerBearer token123", ""},
	}

	for _, tt := range tests {
		result := ExtractBearerToken(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

func TestFormatEndpoint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"endpoint /path?query", "endpoint-pathquery"},
		{"no changes", "nochanges"},
		{"spaces only", "spacesonly"},
		{"slashes/only", "slashes-only"},
	}

	for _, tt := range tests {
		result := FormatEndpoint(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}