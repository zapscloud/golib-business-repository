package business_repository

import (
	"github.com/zapscloud/golib-business-repository.git/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// UserTypeDao - Contact DAO Repository
type UserTypeDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(userTypeId string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(userTypeId string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(userTypeId string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewUserTypeDao - Contruct UserType Dao
func NewUserTypeDao(client utils.Map, business_id string) UserTypeDao {
	var daoUserType UserTypeDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoUserType = &mongodb_repository.UserTypeMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoUserType != nil {
		// Initialize the Dao
		daoUserType.InitializeDao(client, business_id)
	}

	return daoUserType
}
