package business_repository

import (
	"github.com/zapscloud/golib-business/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// ContactDao - Contact DAO Repository
type ContactDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Contact Details
	Get(contact_id string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Contact
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(contact_id string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(contact_id string) (int64, error)

	// DeleteAll - DeleteAll Collection
	DeleteAll() (int64, error)
}

// NewContactDBDao - Contruct Contact Dao
func NewContactDao(client utils.Map, business_id string) ContactDao {
	var daoContact ContactDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoContact = &mongodb_repository.ContactMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoContact != nil {
		// Initialize the Dao
		daoContact.InitializeDao(client, business_id)
	}

	return daoContact
}
