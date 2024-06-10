package main

import (
	"errors"
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

func createApiDirIfNotExists(api string) (string, error) {
	if api == "" {
		return "", errors.New("api name should be a non-empty string")
	}

	apiName := strings.ToLower(api)
	apiDir := filepath.Join("api", apiName)

	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		if err := os.MkdirAll(apiDir, os.ModePerm); err != nil {
			return "", err
		}
	}
	return apiDir, nil
}

func generateApi(api string) error {
	if api == "" {
		return errors.New("api name should be a non-empty string")
	}

	apiName := strings.ToLower(api)
	apiDir := filepath.Join("api", apiName)
	if _, err := os.Stat(apiDir); err == nil {
		fmt.Println(apiName, "already exists")
		return nil
	}

	if err := os.MkdirAll(apiDir, os.ModePerm); err != nil {
		return err
	}

	if err := generateDto(apiDir, apiName); err != nil {
		return err
	}
	if err := generateSchema(apiDir, apiName); err != nil {
		return err
	}
	if err := generateService(apiDir, apiName); err != nil {
		return err
	}
	if err := generateController(apiDir, apiName); err != nil {
		return err
	}
	return nil
}

func generateService(apiDir, apiName string) error {
	apiLower := strings.ToLower(apiName)
	apiCaps := capitalizeFirstLetter(apiName)
	servicePath := filepath.Join(apiDir, fmt.Sprintf("%sservice.go", ""))

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
	Find%s(id primitive.ObjectID) (*schema.%s, error)
}

type service struct {
	network.BaseService
	%sQuery mongo.Query[schema.%s]
}

func New%sService(db mongo.Database, dbQueryTimeout time.Duration) %sService {
	s := service{
		BaseService:  network.NewBaseService(dbQueryTimeout),
		%sQuery: mongo.NewQuery[schema.%s](db, schema.CollectionName),
	}
	return &s
}

func (s *service) Find%s(id primitive.ObjectID) (*schema.%s, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"_id": id}

	msg, err := s.%sQuery.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
`, apiLower, apiLower, apiCaps, apiCaps, apiCaps, apiLower, apiCaps, apiCaps, apiCaps, apiLower, apiCaps, apiCaps, apiCaps, apiLower)

	return os.WriteFile(servicePath, []byte(template), os.ModePerm)
}

func generateController(apiDir, apiName string) error {
	apiLower := strings.ToLower(apiName)
	apiCaps := capitalizeFirstLetter(apiName)
	controllerPath := filepath.Join(apiDir, fmt.Sprintf("%scontroller.go", ""))

	template := fmt.Sprintf(`package %s

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/%s/dto"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
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
	mongoId, err := network.ReqParams(ctx, &coredto.MongoId{})
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	%s, err := c.%sService.Find%s(mongoId.ID)
	if err != nil {
		panic(network.NotFoundError("%s not found", err))
	}

	data, err := network.MapToDto[dto.Info%s](%s)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	network.SuccessResponse("success", data).Send(ctx)
}
`, apiLower, apiLower, apiLower, apiCaps, apiCaps, apiCaps, apiLower, apiLower, apiCaps, apiCaps, apiLower, apiLower, apiCaps, apiLower, apiCaps, apiLower)

	return os.WriteFile(controllerPath, []byte(template), os.ModePerm)
}

func generateSchema(apiDir, apiName string) error {
	schemaDirPath := filepath.Join(apiDir, "schema")
	if err := os.MkdirAll(schemaDirPath, os.ModePerm); err != nil {
		return err
	}

	apiLower := strings.ToLower(apiName)
	apiCaps := capitalizeFirstLetter(apiName)
	schemaPath := filepath.Join(apiDir, fmt.Sprintf("schema/%s.go", apiLower))

	tStr := `package schema

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
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
	doc := %s{
		Field:     field,
		Status:    true,
		CreatedAt: time,
		UpdatedAt: time,
	}
	if err := doc.Validate(); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (doc *%s) GetValue() *%s {
	return doc
}

func (doc *%s) Validate() error {
	validate := validator.New()
	return validate.Struct(doc)
}

func (*%s) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	q := mongo.NewQuery[%s](db, CollectionName)
	q.CreateIndexes(context.Background(), indexes)
}

`

	template := fmt.Sprintf(tStr, apiLower, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps)

	return os.WriteFile(schemaPath, []byte(template), os.ModePerm)
}

func generateDto(apiDir, apiName string) error {
	dtoDirPath := filepath.Join(apiDir, "dto")
	if err := os.MkdirAll(dtoDirPath, os.ModePerm); err != nil {
		return err
	}

	apiLower := strings.ToLower(apiName)
	apiCaps := capitalizeFirstLetter(apiName)
	dtoPath := filepath.Join(apiDir, fmt.Sprintf("dto/create_%s.go", apiLower))

	tStr := `package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Info%s struct {
	ID        primitive.ObjectID ` + "`" + `json:"_id" binding:"required"` + "`" + `
	Field     string             ` + "`" + `json:"field" binding:"required"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"createdAt" binding:"required"` + "`" + `
}

func EmptyInfo%s() *Info%s {
	return &Info%s{}
}

func (d *Info%s) GetValue() *Info%s {
	return d
}

func (d *Info%s) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%%s must be min %%s", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%%s must be max %%s", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
`
	template := fmt.Sprintf(tStr, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps, apiCaps)

	return os.WriteFile(dtoPath, []byte(template), os.ModePerm)
}

type Command string

const (
	Api        Command = "api"
	Controller Command = "controller"
	Service    Command = "service"
	Dto        Command = "dto"
	Schema     Command = "schema"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Error: command from [%s, %s, %s, %s, %s] is required", Api, Controller, Service, Dto, Schema)
		return
	}

	cmd := Command(os.Args[1])

	switch cmd {
	case Api, Controller, Service, Dto, Schema:
	default:
		fmt.Printf("Error: cmd must be from [%s, %s, %s, %s, %s]", Api, Controller, Service, Dto, Schema)
		return
	}

	if len(os.Args) < 3 {
		switch cmd {
		case Api:
			fmt.Println("Error: api name is required")
			return
		case Controller, Service, Dto, Schema:
			fmt.Println("Error: api name for " + cmd + " is required")
			return
		}
	}

	apiName := os.Args[2]

	if ok := validateName(apiName); !ok {
		fmt.Println("Error: api {" + apiName + "} should starts with a-z")
	}

	switch cmd {
	case Api:
		if err := generateApi(apiName); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	case Controller:
		dirName, err := createApiDirIfNotExists(apiName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if err := generateController(dirName, apiName); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	case Service:
		dirName, err := createApiDirIfNotExists(apiName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if err := generateService(dirName, apiName); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	case Dto:
		dirName, err := createApiDirIfNotExists(apiName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if err := generateDto(dirName, apiName); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	case Schema:
		dirName, err := createApiDirIfNotExists(apiName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if err := generateSchema(dirName, apiName); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}

func validateName(s string) bool {
	for _, char := range s[:1] {
		if char < 'a' || char > 'z' {
			return false
		}
	}
	return true
}
