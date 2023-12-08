package business_repository

import (
	"github.com/zapscloud/golib-business-repository/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// AdvertisementDao - Card DAO Repository
type PaymentDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)

	// List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get - Get Payment Details
	Get(paymentid string) (utils.Map, error)

	// Find - Find by code
	Find(filter string) (utils.Map, error)

	// Create - Create Payment
	Create(indata utils.Map) (utils.Map, error)

	// Update - Update Collection
	Update(paymentid string, indata utils.Map) (utils.Map, error)

	// Delete - Delete Collection
	Delete(paymentid string) (int64, error)

	// // DeleteAll - Delete All Collection
	// DeleteAll() (int64, error)
}

// NewPaymentDao - Contruct Payment Dao
func NewPaymentDao(client utils.Map, businessid string) PaymentDao {
	var daoPayment PaymentDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoPayment = &mongodb_repository.PaymentMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoPayment != nil {
		// Initialize the Dao
		daoPayment.InitializeDao(client, businessid)
	}

	return daoPayment
}
