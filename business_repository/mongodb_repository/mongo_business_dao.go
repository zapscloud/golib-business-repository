package mongodb_repository

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-business-repository.git/business_common"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AdvertisementDao - Card DAO Repository
type BusinessMongoDBDao struct {
	client     utils.Map
	businessID string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *BusinessMongoDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize Mongodb DAO")
	p.client = client
	p.businessID = businessId
}

// Insert - Insert Collection
func (p *BusinessMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Business Save - Begin", indata)

	//insert business profile
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessProfiles)
	if err != nil {
		return indata, err
	}
	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	// Insert a single document
	insertResult, err := collection.InsertOne(ctx, indata)
	if err != nil {
		log.Println("Error in insert ", err)
		return indata, err
	}

	log.Println("Inserted a single document: ", insertResult.InsertedID)
	log.Println("Save - End", indata[business_common.FLD_BUSINESS_ID])

	return indata, err
}

// Get - Get user details
func (p *BusinessMongoDBDao) Get(businessid string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("BusinessMongoDBDao::Find:: Begin ", businessid)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessProfiles)
	log.Println("Find:: Got Collection ")

	filter := bson.D{
		{Key: business_common.FLD_BUSINESS_ID, Value: businessid},
		{Key: db_common.FLD_IS_DELETED, Value: false}, {}}

	log.Println("Find:: Got filter ", filter)

	singleResult := collection.FindOne(ctx, filter)
	if singleResult.Err() != nil {
		log.Println("Find:: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}
	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Println("BusinessMongoDBDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Find - Find by code
func (p *BusinessMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("ContactDBDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbAppContacts)
	log.Println("Find:: Got Collection ", err)

	bfilter := bson.D{}
	err = bson.UnmarshalExtJSON([]byte(filter), true, &bfilter)
	if err != nil {
		fmt.Println("Error on filter Unmarshal", err)
	}
	bfilter = append(bfilter,
		bson.E{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	log.Println("Find:: Got filter ", bfilter)
	singleResult := collection.FindOne(ctx, bfilter)
	if singleResult.Err() != nil {
		log.Println("Find:: Record not found ", singleResult.Err())
		return result, singleResult.Err()
	}
	singleResult.Decode(&result)
	if err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Println("ContactDBDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Update - Update Collection
func (p *BusinessMongoDBDao) Update(indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessProfiles)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)

	// Update a single document
	log.Printf("Update - Values %v", indata)

	filter := bson.D{{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID}}
	//filter = append(filter, bson.E{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID})

	updateResult, err := collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult.ModifiedCount)

	log.Println("Update - End")
	return p.Get(p.businessID)
}

// Delete - Delete Collection
func (p *BusinessMongoDBDao) Delete(businessid string) (int64, error) {

	log.Println("BusinessMongoDBDao::Delete - Begin ", businessid)
	//Business profile
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessProfiles)
	if err != nil {
		return 0, err
	}
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filter := bson.D{{Key: business_common.FLD_BUSINESS_ID, Value: businessid}}
	res, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("BusinessMongoDBDao::Delete - End deleted %v documents\n", res.DeletedCount)

	//Business_user
	collectionUser, ctxUser, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	log.Println("user inside ")
	if err != nil {
		return 0, err
	}
	optsUser := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filterUser := bson.D{{Key: business_common.FLD_BUSINESS_ID, Value: businessid}}
	resUser, err := collectionUser.DeleteOne(ctxUser, filterUser, optsUser)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("BusinessMongoDBDao::Delete - End deleted %v documents\n", resUser.DeletedCount)

	return res.DeletedCount, nil
}
