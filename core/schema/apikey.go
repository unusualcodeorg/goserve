package schema

import (
	"time"
)

const ApiKeyCollectionName = "api_keys"

type Permission string

const (
	GeneralPermission Permission = "GENERAL"
)

type ApiKey struct {
	ID          string       `bson:"_id,omitempty"`
	Key         string       `bson:"key" validate:"required,max=1024"`
	Version     int          `bson:"version" validate:"required,min=1,max=100"`
	Permissions []Permission `bson:"permissions" validate:"required"`
	Comments    []string     `bson:"comments" validate:"required,max=1000"`
	Status      bool         `bson:"status" validate:"-"`
	CreatedAt   time.Time    `bson:"createdAt" validate:"-"`
	UpdatedAt   time.Time    `bson:"updatedAt" validate:"-"`
}

func NewApiKey(key string, version int, permissions []Permission, comments []string) *ApiKey {
	currentTime := time.Now()
	return &ApiKey{
		Key:         key,
		Version:     version,
		Permissions: permissions,
		Comments:    comments,
		Status:      true,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}
}
