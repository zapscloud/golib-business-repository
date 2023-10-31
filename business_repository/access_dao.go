package business_repository

import (
	"github.com/zapscloud/golib-business-repository.git/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// AccessDao - Access DAO Repository
type AccessDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(sys_filter string, filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Access Details
	Get(accessid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	GrantPermission(indata utils.Map) (utils.Map, error)

	// RevokePermission - RevokePermission Collection
	RevokePermission(accessid string) (int64, error)
}

// NewaccessMongoDao - Contruct Access Dao
func NewAccessDao(client utils.Map, businessid string) AccessDao {
	var daoAccess AccessDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoAccess = &mongodb_repository.AccessMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoAccess != nil {
		// Initialize the Dao
		daoAccess.InitializeDao(client, businessid)
	}

	return daoAccess
}
