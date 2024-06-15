package startup

import (
	auth "github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/model"
	contact "github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/model"
	user "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/framework/mongo"
)

func EnsureDbIndexes(db mongo.Database) {
	go mongo.Document[auth.Keystore](&auth.Keystore{}).EnsureIndexes(db)
	go mongo.Document[auth.ApiKey](&auth.ApiKey{}).EnsureIndexes(db)
	go mongo.Document[user.User](&user.User{}).EnsureIndexes(db)
	go mongo.Document[user.Role](&user.Role{}).EnsureIndexes(db)
	go mongo.Document[contact.Message](&contact.Message{}).EnsureIndexes(db)
}
