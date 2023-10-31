package zapsdb_repository

import (
	"log"

	"github.com/zapscloud/golib-business-repository.git/business_common"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/zapsdb_utils"
	"github.com/zapscloud/golib-utils/utils"
)

// UserZapsDBDao - User DAO Repository
type UserZapsDBDao struct {
	client utils.Map
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (t *UserZapsDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize Zaps DAO")

	connection, _ := zapsdb_utils.GetConnection(t.client)
	collection, err := connection.GetCollection(business_common.DbBusinessUsers)
	if err != nil {
		log.Println("Business User Table not found, createTable in ZapsDB")
		collection, err = connection.CreateCollection(business_common.DbBusinessUsers, business_common.FLD_USER_ID, "Business User")
		if err != nil {
			log.Println("Failed to create Business User Table in ZapsDB")
		}
	}
	log.Println("Business User Collection", collection)
}

// List - List all Collections
func (t *UserZapsDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	log.Println("Begin - Find All Collection Dao", business_common.DbBusinessUsers)

	connection, _ := zapsdb_utils.GetConnection(t.client)
	response, err := connection.GetMany(business_common.DbBusinessUsers, filter, sort, skip, limit)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	log.Println("End - Find All Collection Dao", response)

	return response, nil
}

// Get - Get details
func (t *UserZapsDBDao) Get(userid string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("UserDBDao::Find:: Begin ", userid)

	connection, _ := zapsdb_utils.GetConnection(t.client)
	singleResult, err := connection.GetOne(business_common.DbBusinessUsers, userid, "")
	if err != nil {
		log.Println("Find:: Record not found ", err)
		return result, err
	}
	result = singleResult
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	log.Printf("UserDBDao::Find:: End Found a single document: %+v\n", result)
	return result, nil
}

// Find - Find by code
func (t *UserZapsDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("UserDBDao::Find:: Begin ", filter)

	connection, _ := zapsdb_utils.GetConnection(t.client)
	singleResult, err := connection.FindOne(business_common.DbBusinessUsers, filter, "")
	if err != nil {
		log.Println("Find:: Record not found ", err)
		return result, err
	}
	result = singleResult
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	log.Printf("UserDBDao::Find:: End Found a single document: %+v\n", result)
	return result, nil
}

// Create - Create Collection
func (t *UserZapsDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Business User Save - Begin", indata)
	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	connection, txnid := zapsdb_utils.GetConnection(t.client)
	insertResult, err := connection.Insert(business_common.DbBusinessUsers, indata, txnid)
	if err != nil {
		log.Println("Error in insert ", err)
		return utils.Map{}, err

	}
	log.Println("Inserted a single document: ", insertResult[business_common.FLD_USER_ID])
	log.Println("Save - End", indata[business_common.FLD_USER_ID])

	return indata, err
}

// Update - Update Collection
func (t *UserZapsDBDao) Update(userid string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)
	log.Printf("Update - Values %v", indata)

	connection, txnid := zapsdb_utils.GetConnection(t.client)
	updateResult, err := connection.UpdateOne(business_common.DbBusinessUsers, userid, indata, txnid)
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult)

	log.Println("Update - End")
	return t.Find(userid)
}

// Delete - Delete Collection
func (t *UserZapsDBDao) Delete(userid string) (int64, error) {

	log.Println("UserDBDao::Delete - Begin ", userid)

	connection, txnid := zapsdb_utils.GetConnection(t.client)
	res, err := connection.DeleteOne(business_common.DbBusinessUsers, userid, txnid)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("UserDBDao::Delete - End deleted %v documents\n", res)
	return 1, nil
}
