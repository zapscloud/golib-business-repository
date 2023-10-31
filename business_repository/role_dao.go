package business_repository

import (
	"github.com/zapscloud/golib-business/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-business/business_repository/zapsdb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// RoleDao - Role DAO Repository
type RoleDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// GetDetails - GetDetails Role Details
	GetDetails(roleid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Role
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(roleid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(roleid string) (int64, error)
}

// NewRoleDBDao - Contruct Role Dao
func NewRoleDao(client utils.Map, businessid string) RoleDao {
	var daoRole RoleDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoRole = &mongodb_repository.RoleMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		daoRole = &zapsdb_repository.RoleZapsDBDao{}
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoRole != nil {
		// Initialize the Dao
		daoRole.InitializeDao(client, businessid)
	}

	return daoRole
}
