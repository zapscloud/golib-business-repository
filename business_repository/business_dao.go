package business_repository

import (
	"github.com/zapscloud/golib-business-repository.git/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-business-repository.git/business_repository/zapsdb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// AdvertisementDao - Card DAO Repository
type BusinessDao interface {

	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// Insert - Insert Collection
	Create(indata utils.Map) (utils.Map, error)

	// Get Business Detail
	Get(businessid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Update Business Detail
	Update(indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(businessid string) (int64, error)
}

func NewBusinessDao(client utils.Map, businessid string) BusinessDao {
	var daoBusiness BusinessDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoBusiness = &mongodb_repository.BusinessMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		daoBusiness = &zapsdb_repository.BusinessZapsDBDao{}
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoBusiness != nil {
		// Initialize the Dao
		daoBusiness.InitializeDao(client, businessid)
	}

	return daoBusiness
}
