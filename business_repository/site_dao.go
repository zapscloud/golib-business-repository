package business_repository

import (
	"github.com/zapscloud/golib-business-repository/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// SiteDao - Site DAO Repository
type SiteDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Site Details
	Get(siteid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Site
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(siteid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(siteid string) (int64, error)

	// DeleteAll - Delete All Collection
	DeleteAll() (int64, error)
}

// NewSiteDao - Contruct Site Dao
func NewSiteDao(client utils.Map, businessid string) SiteDao {
	var daoSite SiteDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoSite = &mongodb_repository.SiteMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoSite != nil {
		// Initialize the Dao
		daoSite.InitializeDao(client, businessid)
	}

	return daoSite
}
