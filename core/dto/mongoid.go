package coredto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoId struct {
	Id string             `uri:"id" binding:"required" validate:"required,len=24"`
	ID primitive.ObjectID `uri:"-" validate:"-"`
}

func (d *MongoId) GetValue() *MongoId {
	id, err := mongo.NewObjectID(d.Id)
	if err != nil {
		panic(network.BadRequestError("id is invalid", errors.New("mongo id is invalid")))
	}
	d.ID = id
	return d
}

// strings.ToLower because gin query param validation does not give back form:"page"
func (d *MongoId) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", strings.ToLower(err.Field())))
		case "len":
			msgs = append(msgs, fmt.Sprintf("%s must be of length %s", strings.ToLower(err.Field()), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", strings.ToLower(err.Field())))
		}
	}
	return msgs, nil
}
