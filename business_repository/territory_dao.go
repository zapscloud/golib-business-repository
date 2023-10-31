package business_repository

import (
	"github.com/zapscloud/golib-business/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// AdvertisementDao - Card DAO Repository
type TerritoryDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Territory Details
	Get(territoryid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Territory
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(territoryid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(territoryid string) (int64, error)

	// DeleteAll - Delete All Collection
	DeleteAll() (int64, error)
}

// NewTerritoryDao - Contruct Territory Dao
func NewTerritoryDao(client utils.Map, businessid string) TerritoryDao {
	var daoTerritory TerritoryDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoTerritory = &mongodb_repository.TerritoryMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoTerritory != nil {
		// Initialize the Dao
		daoTerritory.InitializeDao(client, businessid)
	}

	return daoTerritory
}
