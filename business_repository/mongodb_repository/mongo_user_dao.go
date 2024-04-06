package mongodb_repository

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-business-repository/business_common"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/mongo_utils"
	"github.com/zapscloud/golib-hr-repository/hr_common"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-utils/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserMongoDBDao - User DAO Repository
type UserMongoDBDao struct {
	client     utils.Map
	businessID string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func (p *UserMongoDBDao) InitializeDao(client utils.Map, businessId string) {
	log.Println("Initialize User Mongodb DAO")
	p.client = client
	p.businessID = businessId
}

// List - List all Collections
func (p *UserMongoDBDao) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	var results []utils.Map
	var bFilter bool = false

	log.Println("Begin - Find All Collection Dao", business_common.DbBusinessUsers)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	if err != nil {
		return nil, err
	}

	log.Println("Get Collection - Find All Collection Dao", filter, len(filter), sort, len(sort))

	filterdoc := bson.D{}
	if len(filter) > 0 {
		// filters, _ := strconv.Unquote(string(filter))
		err = bson.UnmarshalExtJSON([]byte(filter), true, &filterdoc)
		if err != nil {
			log.Println("Unmarshal Ext JSON error", err)
			log.Println(filterdoc)
		}
	}

	// All Stages
	stages := []bson.M{}
	// Remove unwanted fields
	unsetStage := bson.M{db_common.MONGODB_UNSET: db_common.FLD_DEFAULT_ID}
	stages = append(stages, unsetStage)

	// Match Stage ==================================
	filterdoc = append(filterdoc,
		bson.E{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID},
		bson.E{Key: db_common.FLD_IS_DELETED, Value: false})

	matchStage := bson.M{db_common.MONGODB_MATCH: filterdoc}
	stages = append(stages, matchStage)
	// ==================================================

	// Add Lookup stages ================================
	stages = p.appendListLookups(stages)
	// ==================================================

	if len(sort) > 0 {
		var sortdoc interface{}
		err = bson.UnmarshalExtJSON([]byte(sort), true, &sortdoc)
		if err != nil {
			log.Println("Sort Unmarshal Error ", sort)
		} else {
			sortStage := bson.M{db_common.MONGODB_SORT: sortdoc}
			stages = append(stages, sortStage)
		}
	}

	var filtercount int64 = 0
	if bFilter {
		// Prepare Filter Stages
		filterStages := stages

		// Add Count aggregate
		countStage := bson.M{db_common.MONGODB_COUNT: business_common.FLD_FILTERED_COUNT}
		filterStages = append(filterStages, countStage)

		//log.Println("Aggregate for Count ====>", filterStages, stages)

		// Execute aggregate to find the count of filtered_size
		cursor, err := collection.Aggregate(ctx, filterStages)
		if err != nil {
			log.Println("Error in Aggregate", err)
			return nil, err
		}
		var countResult []utils.Map
		if err = cursor.All(ctx, &countResult); err != nil {
			log.Println("Error in cursor.all", err)
			return nil, err
		}

		//log.Println("Count Filter ===>", countResult)
		if len(countResult) > 0 {
			if dataVal, dataOk := countResult[0][business_common.FLD_FILTERED_COUNT]; dataOk {
				filtercount = int64(dataVal.(int32))
			}
		}
		// log.Println("Count ===>", filtercount)

	} else {
		filtercount, err = collection.CountDocuments(ctx, filterdoc)
		if err != nil {
			return nil, err
		}
	}

	if skip > 0 {
		skipStage := bson.M{db_common.MONGODB_SKIP: skip}
		stages = append(stages, skipStage)
	}

	if limit > 0 {
		limitStage := bson.M{db_common.MONGODB_LIMIT: limit}
		stages = append(stages, limitStage)
	}

	cursor, err := collection.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}

	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	basefilterdoc := bson.D{
		{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID},
		{Key: db_common.FLD_IS_DELETED, Value: false}}

	totalcount, err := collection.CountDocuments(ctx, basefilterdoc)
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []utils.Map{}
	}

	response := utils.Map{
		db_common.LIST_SUMMARY: utils.Map{
			db_common.LIST_TOTALSIZE:    totalcount,
			db_common.LIST_FILTEREDSIZE: filtercount,
			db_common.LIST_RESULTSIZE:   len(results),
		},
		db_common.LIST_RESULT: results,
	}

	return response, nil

}

// Get - Get user details
func (p *UserMongoDBDao) Get(userid string) (utils.Map, error) {
	// Find a single document
	var result utils.Map
	stages := []bson.M{}
	log.Println("UserDBDao::Get:: Begin ", userid)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	if err != nil {
		return nil, err
	}
	log.Println("Find:: Got Collection ")

	filter := bson.D{{Key: business_common.FLD_USER_ID, Value: userid},
		{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID},
		{Key: db_common.FLD_IS_DELETED, Value: false}}

	matchStage := bson.M{db_common.MONGODB_MATCH: filter}
	stages = append(stages, matchStage)

	// Append Lookups
	stages = p.appendListLookups(stages)

	// Aggregate the stages
	singleResult, err := collection.Aggregate(ctx, stages)
	if err != nil {
		log.Println("GetDetails:: Error in aggregation: ", err)
		return result, err
	}

	if !singleResult.Next(ctx) {
		// No matching document found
		log.Println("GetDetails:: Record not found")
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Record Not Found", ErrorDetail: "Given UserID is not found"}
		return result, err
	}

	if err := singleResult.Decode(&result); err != nil {
		log.Println("Error in decode", err)
		return result, err
	}

	// Remove fields from result
	result = db_common.AmendFldsForGet(result)

	log.Println("UserDBDao::Get:: End Found a single document: \n", err)
	return result, nil
}

// Find - Find by code
func (p *UserMongoDBDao) Find(filter string) (utils.Map, error) {
	// Find a single document
	var result utils.Map

	log.Println("UserDBDao::Find:: Begin ", filter)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
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

	log.Println("UserDBDao::Find:: End Found a single document: \n", err)
	return result, nil
}

// Create - Create Collection
func (p *UserMongoDBDao) Create(indata utils.Map) (utils.Map, error) {

	log.Println("Business User Save - Begin", indata)
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	if err != nil {
		return utils.Map{}, err
	}
	// Add Fields for Create
	indata = db_common.AmendFldsforCreate(indata)

	// Insert a single document
	insertResult, err := collection.InsertOne(ctx, indata)
	if err != nil {
		log.Println("Error in insert ", err)
		return utils.Map{}, err

	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
	log.Println("Save - End", indata[business_common.FLD_USER_ID])

	return indata, err
}

// Update - Update Collection
func (p *UserMongoDBDao) Update(userid string, indata utils.Map) (utils.Map, error) {

	log.Println("Update - Begin")
	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	if err != nil {
		return utils.Map{}, err
	}
	// Modify Fields for Update
	indata = db_common.AmendFldsforUpdate(indata)

	// Update a single document
	log.Printf("Update - Values %v", indata)

	filter := bson.D{{Key: business_common.FLD_USER_ID, Value: userid}}
	filter = append(filter, bson.E{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID})

	updateResult, err := collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: indata}})
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("Update a single document: ", updateResult.ModifiedCount)

	log.Println("Update - End")
	return p.Get(userid)
}

// Delete - Delete Collection
func (p *UserMongoDBDao) Delete(userid string) (int64, error) {

	log.Println("UserDBDao::Delete - Begin ", userid)

	collection, ctx, err := mongo_utils.GetMongoDbCollection(p.client, business_common.DbBusinessUsers)
	if err != nil {
		return 0, err
	}
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    db_common.LOCALE,
		Strength:  1,
		CaseLevel: false,
	})

	filter := bson.D{{Key: business_common.FLD_USER_ID, Value: userid}}

	filter = append(filter, bson.E{Key: business_common.FLD_BUSINESS_ID, Value: p.businessID})

	res, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		log.Println("Error in delete ", err)
		return 0, err
	}
	log.Printf("UserDBDao::Delete - End deleted %v documents\n", res.DeletedCount)
	return res.DeletedCount, nil
}

func (p *UserMongoDBDao) appendListLookups(stages []bson.M) []bson.M {
	// Lookup Stage for Token PlatformAppUsers ========================================
	lookupStage := bson.M{
		db_common.MONGODB_LOOKUP: bson.M{
			db_common.MONGODB_STR_FROM:         business_common.DbBusinessRoles,
			db_common.MONGODB_STR_LOCALFIELD:   business_common.FLD_USER_ROLES + "." + business_common.FLD_ROLE_ID,
			db_common.MONGODB_STR_FOREIGNFIELD: business_common.FLD_ROLE_ID,
			db_common.MONGODB_STR_AS:           business_common.FLD_ROLD_INFO,
			db_common.MONGODB_STR_PIPELINE: []bson.M{
				// Match BusinessId
				{db_common.MONGODB_MATCH: bson.M{
					business_common.FLD_BUSINESS_ID: p.businessID,
				}},
				// Remove following fields from result-set
				{db_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID: 0,
					db_common.FLD_IS_DELETED: 0,
					db_common.FLD_CREATED_AT: 0,
					db_common.FLD_UPDATED_AT: 0,
				}},
			},
		},
	}
	// Add it to Aggregate Stage Shift
	stages = append(stages, lookupStage)

	// Lookup Stage for HR staff ==========================================
	lookupStage = bson.M{
		db_common.MONGODB_LOOKUP: bson.M{
			db_common.MONGODB_STR_FROM:         hr_common.DbHrStaffs,
			db_common.MONGODB_STR_LOCALFIELD:   platform_common.FLD_APP_USER_ID,
			db_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_STAFF_ID,
			db_common.MONGODB_STR_AS:           business_common.FLD_HR_STAFF_INFO,
			db_common.MONGODB_STR_PIPELINE: []bson.M{
				// Match BusinessId
				{db_common.MONGODB_MATCH: bson.M{
					business_common.FLD_BUSINESS_ID: p.businessID,
				}},
				// Remove following fields from result-set
				{db_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID: 0,
					db_common.FLD_IS_DELETED: 0,
					db_common.FLD_CREATED_AT: 0,
					db_common.FLD_UPDATED_AT: 0}},
			},
		},
	}
	// //Add it to Aggregate Stage
	stages = append(stages, lookupStage)

	// Lookup Stage for Department ==========================================
	lookupStage = bson.M{
		db_common.MONGODB_LOOKUP: bson.M{
			db_common.MONGODB_STR_FROM:         hr_common.DbHrDepartments,
			db_common.MONGODB_STR_LOCALFIELD:   business_common.FLD_HR_STAFF_INFO + "." + hr_common.FLD_STAFF_DATA + "." + hr_common.FLD_DEPARTMENT_ID,
			db_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_DEPARTMENT_ID,
			db_common.MONGODB_STR_AS:           business_common.FLD_DEPARTMENT_INFO,
			db_common.MONGODB_STR_PIPELINE: []bson.M{
				// Match BusinessId
				{db_common.MONGODB_MATCH: bson.M{
					business_common.FLD_BUSINESS_ID: p.businessID,
				}},
				// Remove following fields from result-set
				{db_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID: 0,
					db_common.FLD_IS_DELETED: 0,
					db_common.FLD_CREATED_AT: 0,
					db_common.FLD_UPDATED_AT: 0}},
			},
		},
	}
	// //Add it to Aggregate Stage
	stages = append(stages, lookupStage)

	// Lookup Stage for Positions ==========================================
	lookupStage = bson.M{
		db_common.MONGODB_LOOKUP: bson.M{
			db_common.MONGODB_STR_FROM:         hr_common.DbHrPositions,
			db_common.MONGODB_STR_LOCALFIELD:   business_common.FLD_HR_STAFF_INFO + "." + hr_common.FLD_STAFF_DATA + "." + hr_common.FLD_POSITION_ID,
			db_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_POSITION_ID,
			db_common.MONGODB_STR_AS:           business_common.FLD_POSITION_INFO,
			db_common.MONGODB_STR_PIPELINE: []bson.M{
				// Match BusinessId
				{db_common.MONGODB_MATCH: bson.M{
					business_common.FLD_BUSINESS_ID: p.businessID,
				}},
				// Remove following fields from result-set
				{db_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID: 0,
					db_common.FLD_IS_DELETED: 0,
					db_common.FLD_CREATED_AT: 0,
					db_common.FLD_UPDATED_AT: 0}},
			},
		},
	}
	// //Add it to Aggregate Stage
	stages = append(stages, lookupStage)

	// Lookup Stage for Designations ==========================================
	lookupStage = bson.M{
		db_common.MONGODB_LOOKUP: bson.M{
			db_common.MONGODB_STR_FROM:         hr_common.DbHrDesignations,
			db_common.MONGODB_STR_LOCALFIELD:   business_common.FLD_HR_STAFF_INFO + "." + hr_common.FLD_STAFF_DATA + "." + hr_common.FLD_DESIGNATION_ID,
			db_common.MONGODB_STR_FOREIGNFIELD: hr_common.FLD_DESIGNATION_ID,
			db_common.MONGODB_STR_AS:           business_common.FLD_DESIGNATION_INFO,
			db_common.MONGODB_STR_PIPELINE: []bson.M{
				// Match BusinessId
				{db_common.MONGODB_MATCH: bson.M{
					business_common.FLD_BUSINESS_ID: p.businessID,
				}},
				// Remove following fields from result-set
				{db_common.MONGODB_PROJECT: bson.M{
					db_common.FLD_DEFAULT_ID: 0,
					db_common.FLD_IS_DELETED: 0,
					db_common.FLD_CREATED_AT: 0,
					db_common.FLD_UPDATED_AT: 0}},
			},
		},
	}
	// //Add it to Aggregate Stage

	return stages
}
