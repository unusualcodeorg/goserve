package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Helper function to capitalize the first letter of a string
func capitalizeFirstLetter(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

func generateFeature(featureTemplate string) error {
	if featureTemplate == "" {
		return fmt.Errorf("api name should be a non-empty string")
	}

	featureName := strings.ToLower(featureTemplate)
	featureDir := filepath.Join("api", featureName)
	if _, err := os.Stat(featureDir); err == nil {
		fmt.Println(featureName, "already exists")
		return nil
	}

	// Create api directory
	if err := os.MkdirAll(featureDir, os.ModePerm); err != nil {
		return err
	}

	if err := generateDto(featureDir, featureName); err != nil {
		return err
	}
	if err := generateSchema(featureDir, featureName); err != nil {
		return err
	}
	if err := generateService(featureDir, featureName); err != nil {
		return err
	}
	if err := generateController(featureDir, featureName); err != nil {
		return err
	}
	return nil
}

func generateService(featureDir, featureName string) error {
	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	servicePath := filepath.Join(featureDir, fmt.Sprintf("%sservice.go", ""))

	template := fmt.Sprintf(`package %s

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/%s/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type %sService interface {
	Find%s(id *primitive.ObjectID) (*schema.%s, error)
}

type service struct {
	network.BaseService
	%sQuery mongo.DatabaseQuery[schema.%s]
}

func New%sService(db mongo.Database, dbQueryTimeout time.Duration) %sService {
	s := service{
		BaseService:  network.NewBaseService(dbQueryTimeout),
		%sQuery: mongo.NewDatabaseQuery[schema.%s](db, schema.CollectionName),
	}
	return &s
}

func (s *service) Find%s(id *primitive.ObjectID) (*schema.%s, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"_id": id}

	msg, err := s.%sQuery.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
`, featureLower, featureLower, featureCaps, featureCaps, featureCaps, featureLower, featureCaps, featureCaps, featureCaps, featureLower, featureCaps, featureCaps, featureCaps, featureLower)

	return os.WriteFile(servicePath, []byte(template), os.ModePerm)
}

func generateController(featureDir, featureName string) error {
	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	controllerPath := filepath.Join(featureDir, fmt.Sprintf("%scontroller.go", ""))

	template := fmt.Sprintf(`package %s

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/%s/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	%sService %sService
}

func New%sController(
	authMFunc network.GroupMiddlewareFunc,
	authorizeMFunc network.GroupMiddlewareFunc,
	service %sService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/%s", authMFunc, authorizeMFunc),
		%sService:  service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.get%sHandler)
}

func (c *controller) get%sHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := mongo.NewObjectID(id)
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msg, err := c.%sService.Find%s(objectId)
	if err != nil {
		panic(network.NotFoundError("message not found", err))
	}

	data, err := network.MapToDto[dto.Info%s](msg)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	network.SuccessResponse("success", data).Send(ctx)
}
`, featureLower, featureLower, featureLower, featureCaps, featureCaps, featureCaps, featureLower, featureLower, featureCaps, featureCaps, featureLower, featureCaps, featureCaps)

	return os.WriteFile(controllerPath, []byte(template), os.ModePerm)
}

func generateSchema(featureDir, featureName string) error {
	schemaDirPath := filepath.Join(featureDir, "schema")
	if err := os.MkdirAll(schemaDirPath, os.ModePerm); err != nil {
		return err
	}

	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	schemaPath := filepath.Join(featureDir, fmt.Sprintf("schema/%s.go", featureLower))

	tStr := `package schema

import (
	"context"
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "%ss"

type %s struct {
	ID        primitive.ObjectID ` + "`" + `bson:"_id,omitempty" validate:"-"` + "`" + `
	Field     string             ` + "`" + `bson:"field" validate:"required"` + "`" + `
	Status    bool               ` + "`" + `bson:"status" validate:"required"` + "`" + `
	CreatedAt time.Time          ` + "`" + `bson:"createdAt" validate:"required"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `bson:"updatedAt" validate:"required"` + "`" + `
}` + `

func New%s(field string) (*%s, error) {
	time := time.Now()
	m := %s{
		Field:     field,
		Status:    true,
		CreatedAt: time,
		UpdatedAt: time,
	}
	if err := utils.Validate(m); err != nil {
		return nil, err
	}
	return &m, nil
}

func EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	q := mongo.NewDatabaseQuery[%s](db, CollectionName)
	q.CreateIndexes(context.Background(), indexes)
}
`

	template := fmt.Sprintf(tStr, featureLower, featureCaps, featureCaps, featureCaps, featureCaps, featureCaps)

	return os.WriteFile(schemaPath, []byte(template), os.ModePerm)
}

func generateDto(featureDir, featureName string) error {
	dtoDirPath := filepath.Join(featureDir, "dto")
	if err := os.MkdirAll(dtoDirPath, os.ModePerm); err != nil {
		return err
	}

	featureLower := strings.ToLower(featureName)
	featureCaps := capitalizeFirstLetter(featureName)
	dtoPath := filepath.Join(featureDir, fmt.Sprintf("dto/create_%s.go", featureLower))

	tStr := `package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Info%s struct {
	ID        primitive.ObjectID ` + "`" + `json:"_id" binding:"required"` + "`" + `
	Field     string             ` + "`" + `json:"field" binding:"required"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"createdAt" binding:"required"` + "`" + `
}
`

	template := fmt.Sprintf(tStr, featureCaps)

	return os.WriteFile(dtoPath, []byte(template), os.ModePerm)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("api name should be non-empty string")
		return
	}

	featureName := os.Args[1]
	if err := generateFeature(featureName); err != nil {
		fmt.Println("Error:", err)
	}
}