package zapsdb_repository

import (
	"log"

	"github.com/zapscloud/golib-business-repository.git/business_common"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/zapsdb_utils"
	"github.com/zapscloud/golib-utils/utils"
)

// RoleDBDao - Role DAO Repository
type RoleZapsDBDao struct {
	client utils.Map
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (t *RoleZapsDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize Zaps DAO")

	connection, _ := zapsdb_utils.GetConnection(t.client)
	collection, err := connection.GetCollection(business_common.DbBusinessRoles)
	if err != nil {
		log.Println("Business Role Table not found, createTable in ZapsDB")
		collection, err = connection.CreateCollection(business_common.DbBusinessRoles, business_common.FLD_ROLE_ID, "Business Role")
		if err != nil {
			log.Println("Failed to create Business Role Table in ZapsDB")
		}
	}
	log.Println("Business Role Collection", collection)
}

// List - List all Collections
func (t *RoleZapsDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	log.Println("Begin - Find All Collection Dao", business_common.DbBusinessRoles)

	connection, _ := zapsdb_utils.GetConnection(t.client)
	response, err := connection.GetMany(business_common.DbBusinessRoles, filter, sort, skip, limit)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	log.Println("End - Find All Collection Dao", response)

	return response, nil
}

// Get - Get details
func (t *RoleZapsDBDao) GetDetails(roleid string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("RoleDBDao::Find:: Begin ", roleid)

	connection, _ := zapsdb_utils.GetConnection(t.client)
	singleResult, err := connection.GetOne(business_common.DbBusinessRoles, roleid, "")
	if err != nil {
		log.Println("Find:: Record not found ", err)
		return result, err
	}
	result = singleResult
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	log.Printf("RoleDBDao::Find:: End Found a single document: %+v\n", result)
	return result, nil
}

// Find - Find by code
func (t *RoleZapsDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("RoleDBDao::Find:: Begin ", filter)

	connection, _ := zapsdb_utils.GetConnection(t.client)
	singleResult, err := connection.FindOne(business_common.DbBusinessRoles, filter, "")
	if err != nil {
		log.Println("Find:: Record not found ", err)
		return result, err
	}
	result = singleResult
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	log.Printf("RoleDBDao::Find:: End Found a single document: %+v\n", result)
	return result, nil
}

// Create - Create Collection
func (t *RoleZapsDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Business Role Save - Begin", indata)

	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	connection, txnid := zapsdb_utils.GetConnection(t.client)
	insertResult, err := connection.Insert(business_common.DbBusinessRoles, indata, txnid)
	if err != nil {
		log.Println("Error in insert ", err)
		return utils.Map{}, err

	}
	log.Println("Inserted a single document: ", insertResult[business_common.FLD_ROLE_ID])
	log.Println("Save - End", indata[business_common.FLD_ROLE_ID])

	return indata, err
}

// Update - Update Collection
func (t *RoleZapsDBDao) Update(roleid string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)

	log.Printf("Update - Values %v", indata)

	connection, txnid := zapsdb_utils.GetConnection(t.client)
	updateResult, err := connection.UpdateOne(business_common.DbBusinessRoles, roleid, indata, txnid)
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult)

	log.Println("Update - End")
	return t.Find(roleid)
}

// Delete - Delete Collection
func (t *RoleZapsDBDao) Delete(roleid string) (int64, error) {

	log.Println("RoleDBDao::Delete - Begin ", roleid)

	connection, txnid := zapsdb_utils.GetConnection(t.client)
	res, err := connection.DeleteOne(business_common.DbBusinessRoles, roleid, txnid)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("RoleDBDao::Delete - End deleted %v documents\n", res)
	return 1, nil
}
