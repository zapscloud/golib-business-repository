package business_repository

import (
	"github.com/zapscloud/golib-business/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-business/business_repository/zapsdb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// UserDao - User DAO Repository
type UserDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get User Details
	Get(userid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create User
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(userid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(userid string) (int64, error)
}

// NewUserDao - Contruct User Dao
func NewUserDao(client utils.Map, businessid string) UserDao {
	var daoUser UserDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoUser = &mongodb_repository.UserMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		daoUser = &zapsdb_repository.UserZapsDBDao{}
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoUser != nil {
		// Initialize the Dao
		daoUser.InitializeDao(client, businessid)
	}

	return daoUser
}
