package business_common

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-platform-repository/platform_common"
)

// Business module tables =================================
const (
	DbPrefix      = db_common.DB_COLLECTION_PREFIX
	DbAppAccess   = DbPrefix + "app_access"
	DbAppContacts = DbPrefix + "app_contacts"

	DbAppSites       = DbPrefix + "app_sites"
	DbAppTerritories = DbPrefix + "app_territories"

	DbBusinessProfiles = DbPrefix + "business_profiles"
	DbBusinessRoles    = DbPrefix + "business_roles"
	DbBusinessUsers    = DbPrefix + "business_users"

	DbBusinessUserTypes   = DbPrefix + "business_user_types"
	DbBusinessPayments    = DbPrefix + "business_payments"
	DbBusinessPaymentTxns = DbPrefix + "business_payment_txns"
)

// Business module table fields
const (
	FLD_BUSINESS_ID       = platform_common.FLD_BUSINESS_ID
	FLD_BUSINESS_NAME     = platform_common.FLD_BUSINESS_NAME
	FLD_BUSINESS_EMAILID  = platform_common.FLD_BUSINESS_EMAILID
	FLD_BUSINESS_TIMEZONE = "business_timezone"

	FLD_USER_ID                = platform_common.FLD_APP_USER_ID
	FLD_IS_USER_BUSINESS_ADMIN = "is_user_business_admin"
	FLD_USER_ROLES             = "user_roles"

	FLD_ROLE_ID   = platform_common.FLD_APP_ROLE_ID
	FLD_ROLE_NAME = platform_common.FLD_APP_ROLE_NAME
	FLD_ROLE_DESC = platform_common.FLD_APP_ROLE_DESC

	FLD_APP_ACCESS_ID  = "app_access_id"
	FLD_APP_CONTACT_ID = "app_contact_id"

	FLD_APP_SITE_ID      = "app_site_id"
	FLD_APP_TERRITORY_ID = "app_territory_id"

	FLD_USERTYPE_ID   = "usertype_id"
	FLD_USERTYPE_NAME = "usertype_name"
	FLD_USERTYPE_DESC = "usertype_desc"

	FLD_PAYMENT_ID   = "payment_id"
	FLD_PAYMENT_NAME = "payment_name"

	FLD_PAYMENT_TXN_ID = "payment_txn_id"

	FLD_DATE_TIME = "date_time"
)

const (
	DEF_BUSINESS_TIMEZONE = "Asia/Calcutta"

	FLD_USER_INFO            = "user_info"
	FLD_ROLD_INFO            = "role_info"
	FLD_DEPARTMENT_INFO      = "department_info"
	FLD_DESIGNATION_INFO     = "designation_info"
	FLD_POSITION_INFO        = "position_info"
	FLD_HR_STAFF_INFO        = "hr_staff_info"
	FLD_REPORTING_STAFF_INFO = "reporting_staff_info"
	FLD_FILTERED_COUNT       = "filtered_count"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

}

func GetServiceModuleCode() string {
	return "S4"
}
