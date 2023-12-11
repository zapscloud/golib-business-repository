package business_repository

import (
	"github.com/zapscloud/golib-business-repository/business_repository/mongodb_repository"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-utils/utils"
)

// PaymentTxnDao - Card DAO Repository
type PaymentTxnDao interface {
	// InitializeDao
	InitializeDao(client utils.Map, businessId string)
	//List - List all Collections
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	// Get - Get Payment Details
	Get(paymentTxnId string) (utils.Map, error)
	// Find - Find by filter
	Find(filter string) (utils.Map, error)
	// Create - Create Payment
	Create(indata utils.Map) (utils.Map, error)
	// Update - Update Collection
	Update(paymentTxnId string, indata utils.Map) (utils.Map, error)
	// Delete - Delete Collection
	Delete(paymentTxnId string) (int64, error)
}

// NewPaymentTxnDao - Contruct Business Payment Dao
func NewPaymentTxnDao(client utils.Map, business_id string) PaymentTxnDao {
	var daoPaymentTxn PaymentTxnDao = nil

	// Get DatabaseType and no need to validate error
	// since the dbType was assigned with correct value after dbService was created
	dbType, _ := db_common.GetDatabaseType(client)

	switch dbType {
	case db_common.DATABASE_TYPE_MONGODB:
		daoPaymentTxn = &mongodb_repository.PaymentTxnMongoDBDao{}
	case db_common.DATABASE_TYPE_ZAPSDB:
		// *Not Implemented yet*
	case db_common.DATABASE_TYPE_MYSQLDB:
		// *Not Implemented yet*
	}

	if daoPaymentTxn != nil {
		// Initialize the Dao
		daoPaymentTxn.InitializeDao(client, business_id)
	}

	return daoPaymentTxn
}
