package utils

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Validate(obj any) error {
	if err := validator.New().Struct(obj); err != nil {
		return err
	}
	return nil
}

func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
