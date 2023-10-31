package zapsdb_repository

import (
	"log"

	"github.com/zapscloud/golib-business-repository/business_common"
	"github.com/zapscloud/golib-dbutils/zapsdb_utils"
	"github.com/zapscloud/golib-utils/utils"
)

type BusinessZapsDBDao struct {
	client utils.Map
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *BusinessZapsDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize Zaps DAO")

	connection, _ := zapsdb_utils.GetConnection(p.client)
	collection, err := connection.GetCollection(business_common.DbBusinessProfiles)
	if err != nil {
		log.Println("BusinessProfile Table not found, createTable in ZapsDB")
		collection, err = connection.CreateCollection(business_common.DbBusinessProfiles, business_common.FLD_BUSINESS_ID, "Business Profile")
		if err != nil {
			log.Println("Failed to create BusinessProfile Table in ZapsDB")
		}
	}
	log.Println("BusinessProfile Collection", collection)
}

// Create - Create Collection
func (t *BusinessZapsDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Business Save - Begin", indata)

	connection, txnid := zapsdb_utils.GetConnection(t.client)
	insertResult, err := connection.Insert(business_common.DbBusinessProfiles, indata, txnid)
	if err != nil {
		log.Println("Error in insert ", err)
		return indata, err
	}

	log.Println("Inserted a single document: ", insertResult[business_common.FLD_BUSINESS_ID])
	log.Println("Save - End", indata[business_common.FLD_BUSINESS_ID])

	return indata, err
}

// Get - Get details
func (t *BusinessZapsDBDao) Get(businessid string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("userZapsDao::Find:: Begin ", businessid)

	connection, _ := zapsdb_utils.GetConnection(t.client)
	singleResult, err := connection.GetOne(business_common.DbBusinessProfiles, businessid, "")
	if err != nil {
		log.Println("Find:: Record not found ", err)
		return result, err
	}
	result = singleResult
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	log.Printf("userZapsDao::Find:: End Found a single document: %+v\n", result)
	return result, nil
}

// Delete - Delete Collection
func (t *BusinessZapsDBDao) Delete(businessid string) (int64, error) {

	log.Println("BusinessDBDao::Delete - Begin ", businessid)
	//Business profile
	connection, txnid := zapsdb_utils.GetConnection(t.client)
	res, err := connection.DeleteOne(business_common.DbBusinessProfiles, businessid, txnid)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("BusinessDBDao::Delete - End deleted %v documents\n", res)

	//Business_user
	resUser, err := connection.DeleteOne(business_common.DbBusinessUsers, businessid, txnid)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("BusinessDBDao::Delete - End deleted %v documents\n", resUser)

	return 1, nil
}

// Find - Find by code
func (t *BusinessZapsDBDao) Find(filter string) (utils.Map, error) {
	return utils.Map{}, nil
}

// Update - Update Collection
func (t *BusinessZapsDBDao) Update(indata utils.Map) (utils.Map, error) {
	return utils.Map{}, nil
}
