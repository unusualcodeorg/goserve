package startup

import (
	authSchema "github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	userSchema "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	coreSchema "github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
)

func EnsureDbIndexes(db mongo.Database) {
	go authSchema.EnsureRoleIndexes(db)
	go authSchema.EnsureKeystoreIndexes(db)
	go userSchema.EnsureIndexes(db)
	go coreSchema.EnsureIndexes(db)
}
