package middleware

import (
	"DataSource/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"gopkg.in/mgo.v2/bson"
)

func dbHealthCheck() {
	err := config.Session.Ping()
	if err != nil {
		fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
		fmt.Println("MongoDB", err, "->", "URL:", config.MongodbURL, " Database:", config.MongodbDatabase)
		config.Session.Refresh()
		fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
		fmt.Println("MongoDB Reconnect ->", " URL:", config.MongodbURL, " Database:", config.MongodbDatabase)
	}
}

// TestConnection ...
func TestConnection(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/")
	w.Header().Set("Server", "Grafana Datasource Server")
	w.WriteHeader(200)
	msg := "This is DataSource for iFactory/Andon"
	w.Write([]byte(msg))
}

// GetGroup ...
func GetGroup(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/group")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	groupTopoCollection := config.DB.C(config.GroupTopo)
	var groupTopoResults []map[string]interface{}
	groupNameArray := []map[string]string{}

	parentID := requestBody.Get("parentID").MustString()
	fmt.Println("ParentID:", parentID)

	if parentID == "ALL" {
		groupTopoCollection.Find(bson.M{}).All(&groupTopoResults)
	} else {
		groupTopoCollection.Find(bson.M{"ParentID": parentID}).All(&groupTopoResults)
	}

	for _, groupTopoResult := range groupTopoResults {
		temp := map[string]string{"text": groupTopoResult["GroupName"].(string), "value": groupTopoResult["GroupID"].(string)}
		groupNameArray = append(groupNameArray, temp)
	}
	// fmt.Println(groupNameArray)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(groupNameArray)
}

// GetMachine ...
func GetMachine(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/machine")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	machineRawDataCollection := config.DB.C(config.MachineRawData)
	var machineRawDataResults []map[string]interface{}
	machineIDArray := []map[string]string{}

	groupID := requestBody.Get("groupID").MustString()
	fmt.Println("GroupID:", groupID)

	if groupID == "ALL" {
		machineRawDataCollection.Find(bson.M{}).All(&machineRawDataResults)
	} else {
		machineRawDataCollection.Find(bson.M{"GroupID": groupID}).All(&machineRawDataResults)
	}

	for _, machineRawDataResult := range machineRawDataResults {
		temp := map[string]string{"text": machineRawDataResult["MachineName"].(string), "value": machineRawDataResult["MachineID"].(string)}
		machineIDArray = append(machineIDArray, temp)
	}
	// fmt.Println(machineIDArray)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(machineIDArray)
}

// GetAssignee ...
func GetAssignee(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/assignee")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	eventHistCollection := config.DB.C(config.EventHist)
	var assigneeIDByGroupIDResults []string
	var assigneeIDByMachineIDResults []string
	assigneeIDArray := []map[string]string{}

	filterID := requestBody.Get("filterID").MustString()
	fmt.Println("FilterID:", filterID)

	eventHistCollection.Find(bson.M{"GroupID": filterID}).Distinct("PrincipalID", &assigneeIDByGroupIDResults)
	eventHistCollection.Find(bson.M{"MachineID": filterID}).Distinct("PrincipalID", &assigneeIDByMachineIDResults)

	for _, eventHistByGroupIDResult := range assigneeIDByGroupIDResults {
		var oneResult map[string]interface{}
		eventHistCollection.Find(bson.M{"PrincipalID": eventHistByGroupIDResult}).One(&oneResult)
		temp := map[string]string{"text": oneResult["PrincipalName"].(string), "value": oneResult["PrincipalID"].(string)}
		assigneeIDArray = append(assigneeIDArray, temp)
	}

	for _, eventHistByMachineIDResults := range assigneeIDByMachineIDResults {
		var oneResult map[string]interface{}
		eventHistCollection.Find(bson.M{"PrincipalID": eventHistByMachineIDResults}).One(&oneResult)
		temp := map[string]string{"text": oneResult["PrincipalName"].(string), "value": oneResult["PrincipalID"].(string)}
		assigneeIDArray = append(assigneeIDArray, temp)
	}
	// fmt.Println(assigneeIDArray)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(assigneeIDArray)
}

// Search ...
func Search(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/search")
	requestBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Body: ", string(requestBody))

	metrics := []string{"EventLatest", "EventHist", "EventList", "LastMonthAbReasonRank", "MTTD", "MTTR", "MTBF", "FlowCharting-Machines", "FlowCharting-TPC", "Singlestat-Machines", "Singlestat-Event", "Singlestat-MeanTimeCompute"}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(metrics)
}

// Query ...
func Query(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	grafnaResponseArray := []map[string]interface{}{}
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/query")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	fmt.Println(requestBody.Get("targets").MustArray())

	for indexOfTargets := 0; indexOfTargets < len(requestBody.Get("targets").MustArray()); indexOfTargets++ {
		refID := requestBody.Get("targets").GetIndex(indexOfTargets).Get("refId").MustString()
		dataType := requestBody.Get("targets").GetIndex(indexOfTargets).Get("type").MustString()
		groupID := requestBody.Get("targets").GetIndex(indexOfTargets).Get("groupID").MustString()
		if strings.HasPrefix(groupID, "$") {
			temp := strings.Split(groupID, "$")
			varName := temp[1]
			groupID = requestBody.Get("scopedVars").Get(varName).Get("value").MustString()
		}
		machineID := requestBody.Get("targets").GetIndex(indexOfTargets).Get("machineID").MustString()
		if strings.HasPrefix(machineID, "$") {
			temp := strings.Split(machineID, "$")
			varName := temp[1]
			machineID = requestBody.Get("scopedVars").Get(varName).Get("value").MustString()
		}
		metrics := requestBody.Get("targets").GetIndex(indexOfTargets).Get("metrics").MustString()
		fromTimeString := requestBody.Get("scopedVars").Get("__from").Get("value").MustString()
		fromUnixTime, _ := strconv.ParseInt(fromTimeString, 10, 64)
		// fmt.Println("FromUnixTime:", fromUnixTime)
		temp := fromUnixTime / 1000
		fromTime := time.Unix(temp, 0)
		fromTime = fromTime.In(config.TaipeiTimeZone)
		toTimeSting := requestBody.Get("scopedVars").Get("__to").Get("value").MustString()
		toUnixTime, _ := strconv.ParseInt(toTimeSting, 10, 64)
		// fmt.Println("ToUnixTime:", toUnixTime)
		temp = toUnixTime / 1000
		toTime := time.Unix(temp, 0)
		toTime = toTime.In(config.TaipeiTimeZone)
		// fmt.Println("  RefID:", refID, "DataType:", dataType, " GroupID:", groupID, " MachineID:", machineID, " Metrics:", metrics, " From:", fromUnixTime, fromTime.Format(time.RFC3339), " To:", toUnixTime, toTime.Format(time.RFC3339))

		if dataType == "table" {
			switch metrics {
			case "EventLatest":
				grafnaResponseArray = append(grafnaResponseArray, eventLatest(groupID, machineID))
			case "EventHist":
				grafnaResponseArray = append(grafnaResponseArray, eventHist(groupID, machineID))
			case "EventList":
				grafnaResponseArray = append(grafnaResponseArray, eventList(groupID, fromTime, toTime))
			case "Singlestat-Event":
				grafnaResponseArray = append(grafnaResponseArray, eventSinglestat(groupID, machineID))
			case "Singlestat-Machines":
				grafnaResponseArray = append(grafnaResponseArray, machinesSinglestat(groupID))
			case "LastMonthAbReasonRank":
				grafnaResponseArray = append(grafnaResponseArray, abnormalReasonRank(groupID, fromTime, toTime))
			case "FlowCharting-Machines":
				grafnaResponseArray = append(grafnaResponseArray, machinesFlowCharting(groupID, refID))
			case "FlowCharting-TPC":
				grafnaResponseArray = append(grafnaResponseArray, tpcFlowCharting(groupID, refID))
			case "MTTD", "MTTR", "MTBF":
				grafnaResponseArray = append(grafnaResponseArray, v3MeanTimeCompute(metrics, groupID, fromTime, toTime))
				// v2
				// grafnaResponseArrayElement, totalComputeValue := v2MeanTimeCompute(metrics, groupID, fromTime, toTime)
				// grafnaResponseArray = append(grafnaResponseArray, grafnaResponseArrayElement)
				// grafnaResponseArrayElement = meanTimeComputeSinglestat(metrics, groupID, totalComputeValue, fromTime, toTime)
				// grafnaResponseArray = append(grafnaResponseArray, grafnaResponseArrayElement)
			}
		} else {
			switch metrics {
			// case "MTTD", "MTTR", "MTBF":
			// 	grafnaResponseArray = append(grafnaResponseArray, v1MeanTimeCompute(metrics, groupID, fromTime, toTime))
			}
		}
	}

	// jsonStr, _ := json.Marshal(grafnaResponseArray)
	// fmt.Println(string(jsonStr))
	// fmt.Println(grafnaResponseArray)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(grafnaResponseArray)
}

// Table
func eventLatest(groupID string, machineID string) map[string]interface{} {
	eventLatestcollection := config.DB.C(config.EventLatest)
	var eventLatestResults []map[string]interface{}
	rows := []interface{}{}
	if machineID == "" {
		eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)
	} else {
		eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}}).All(&eventLatestResults)
	}

	for _, eventLatestResult := range eventLatestResults {

		// for key,value := range eventLatestResult {
		// 	var row []interface{}
		// 	row = append(row, value)
		// }

		var row []interface{}
		row = append(row, eventLatestResult["EventID"])
		row = append(row, eventLatestResult["EventCode"])
		row = append(row, eventLatestResult["Type"])
		row = append(row, eventLatestResult["GroupID"])
		row = append(row, eventLatestResult["GroupName"])
		row = append(row, eventLatestResult["MachineID"])
		row = append(row, eventLatestResult["MachineName"])
		row = append(row, eventLatestResult["AbnormalStartTime"])
		row = append(row, eventLatestResult["ProcessingStatusCode"])
		row = append(row, eventLatestResult["ProcessingProgress"])
		row = append(row, eventLatestResult["AbnormalLastingSecond"])
		row = append(row, eventLatestResult["ShouldRepairTime"])
		row = append(row, eventLatestResult["PlanRepairTime"])
		row = append(row, eventLatestResult["TPCID"])
		row = append(row, eventLatestResult["TPCName"])
		row = append(row, eventLatestResult["PrincipalID"])
		row = append(row, eventLatestResult["PrincipalName"])
		row = append(row, eventLatestResult["AbnormalReason"])
		row = append(row, eventLatestResult["AbnormalSolution"])
		row = append(row, eventLatestResult["AbnormalCode"])
		row = append(row, eventLatestResult["AbnormalPosition"])
		// fmt.Println(row)
		rows = append(rows, row)
	}
	// fmt.Println(rows)

	columns := []map[string]string{
		{"text": "EventID", "type": "string"},
		{"text": "EventCode", "type": "number"},
		{"text": "Type", "type": "string"},
		{"text": "GroupID", "type": "string"},
		{"text": "GroupName", "type": "string"},
		{"text": "MachineID", "type": "string"},
		{"text": "MachineName", "type": "string"},
		{"text": "AbnormalStartTime", "type": "time"},
		{"text": "ProcessingStatusCode", "type": "number"},
		{"text": "ProcessingProgress", "type": "string"},
		{"text": "AbnormalLastingSecond", "type": "number"},
		{"text": "ShouldRepairTime", "type": "time"},
		{"text": "PlanRepairTime", "type": "time"},
		{"text": "TPCID", "type": "string"},
		{"text": "TPCName", "type": "string"},
		{"text": "PrincipalID", "type": "string"},
		{"text": "PrincipalName", "type": "string"},
		{"text": "AbnormalReason", "type": "string"},
		{"text": "AbnormalSolution", "type": "string"},
		{"text": "AbnormalCode", "type": "number"},
		{"text": "AbnormalPosition", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}

	return grafanaData
}

func eventList(groupID string, startTime time.Time, endTime time.Time) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventHistCollection := config.DB.C(config.EventHist)
	rows := []interface{}{}

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {

		// for key,value := range eventLatestResult {
		// 	var row []interface{}
		// 	row = append(row, value)
		// }

		var row []interface{}
		row = append(row, eventLatestResult["EventID"])
		row = append(row, eventLatestResult["EventCode"])
		row = append(row, eventLatestResult["Type"])
		row = append(row, eventLatestResult["GroupID"])
		row = append(row, eventLatestResult["GroupName"])
		row = append(row, eventLatestResult["MachineID"])
		row = append(row, eventLatestResult["MachineName"])
		row = append(row, eventLatestResult["AbnormalStartTime"])
		row = append(row, eventLatestResult["ProcessingStatusCode"])
		row = append(row, eventLatestResult["ProcessingProgress"])
		row = append(row, eventLatestResult["AbnormalLastingSecond"])
		row = append(row, eventLatestResult["ShouldRepairTime"])
		row = append(row, eventLatestResult["PlanRepairTime"])
		row = append(row, eventLatestResult["TPCID"])
		row = append(row, eventLatestResult["PrincipalID"])
		row = append(row, eventLatestResult["PrincipalName"])
		row = append(row, eventLatestResult["AbnormalReason"])
		row = append(row, eventLatestResult["AbnormalSolution"])
		row = append(row, eventLatestResult["AbnormalCode"])
		row = append(row, eventLatestResult["AbnormalPosition"])
		// fmt.Println(row)
		rows = append(rows, row)
	}

	var eventHistResults []map[string]interface{}
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventHistResults)
	for _, eventHistResult := range eventHistResults {

		// for key,value := range eventLatestResult {
		// 	var row []interface{}
		// 	row = append(row, value)
		// }

		var row []interface{}
		row = append(row, eventHistResult["EventID"])
		row = append(row, eventHistResult["EventCode"])
		row = append(row, eventHistResult["Type"])
		row = append(row, eventHistResult["GroupID"])
		row = append(row, eventHistResult["GroupName"])
		row = append(row, eventHistResult["MachineID"])
		row = append(row, eventHistResult["MachineName"])
		row = append(row, eventHistResult["AbnormalStartTime"])
		row = append(row, eventHistResult["ProcessingStatusCode"])
		row = append(row, eventHistResult["ProcessingProgress"])
		row = append(row, eventHistResult["AbnormalLastingSecond"])
		row = append(row, eventHistResult["ShouldRepairTime"])
		row = append(row, eventHistResult["PlanRepairTime"])
		row = append(row, eventHistResult["TPCID"])
		row = append(row, eventHistResult["PrincipalID"])
		row = append(row, eventHistResult["PrincipalName"])
		row = append(row, eventHistResult["AbnormalReason"])
		row = append(row, eventHistResult["AbnormalSolution"])
		row = append(row, eventHistResult["AbnormalCode"])
		row = append(row, eventHistResult["AbnormalPosition"])
		// fmt.Println(row)
		rows = append(rows, row)
	}

	// fmt.Println(rows)

	columns := []map[string]string{
		{"text": "EventID", "type": "string"},
		{"text": "EventCode", "type": "number"},
		{"text": "Type", "type": "string"},
		{"text": "GroupID", "type": "string"},
		{"text": "GroupName", "type": "string"},
		{"text": "MachineID", "type": "string"},
		{"text": "MachineName", "type": "string"},
		{"text": "AbnormalStartTime", "type": "time"},
		{"text": "ProcessingStatusCode", "type": "number"},
		{"text": "ProcessingProgress", "type": "string"},
		{"text": "AbnormalLastingSecond", "type": "number"},
		{"text": "ShouldRepairTime", "type": "time"},
		{"text": "PlanRepairTime", "type": "time"},
		{"text": "TPCID", "type": "string"},
		{"text": "PrincipalID", "type": "string"},
		{"text": "PrincipalName", "type": "string"},
		{"text": "AbnormalReason", "type": "string"},
		{"text": "AbnormalSolution", "type": "string"},
		{"text": "AbnormalCode", "type": "number"},
		{"text": "AbnormalPosition", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}

	return grafanaData
}

func eventSinglestat(groupID string, machineID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventHistCollection := config.DB.C(config.EventHist)
	tpcListCollection := config.DB.C(config.TPCList)
	nowTime := time.Now().In(config.TaipeiTimeZone)
	starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)
	var rows []interface{}

	var eventHistResults []map[string]interface{}
	var tpcListResults []map[string]interface{}
	var eventLatestResults []map[string]interface{}
	if machineID == "" {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$match": bson.M{"$expr": bson.M{"$or": []bson.M{{"ProcessingStatusCode": 5}, {"ProcessingStatusCode": 6}}}}}}).All(&eventHistResults)
		tpcListCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&tpcListResults)
		eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)
	} else {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$match": bson.M{"$expr": bson.M{"$or": []bson.M{{"ProcessingStatusCode": 5}, {"ProcessingStatusCode": 6}}}}}}).All(&eventHistResults)
		tpcListCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}}).All(&tpcListResults)
		eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}}).All(&eventLatestResults)
	}
	// fmt.Println("EventLatestResults:", len(eventLatestResults))
	// fmt.Println("EventHistResults:", eventHistResults)
	// fmt.Println("TPCListResults:", tpcListResults)
	numOfCompleted := len(eventHistResults)
	numOfTPC := len(tpcListResults)
	numOfLatestEvents := len(eventLatestResults)
	numOfTotalEvents := len(eventLatestResults) + len(eventHistResults)

	numOfCompletedToday := 0
	numOfOverdue := 0
	numOfOpen := 0
	numOfInProgressing := 0
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["PlanRepairTime"] != nil {
			planRepairTime := eventLatestResult["PlanRepairTime"].(time.Time)
			if planRepairTime.After(starttimeT) && planRepairTime.Before(endtimeT) {
				numOfCompletedToday++
			}
		}

		processingStatusCode := eventLatestResult["ProcessingStatusCode"]
		if processingStatusCode == 0 || processingStatusCode == 1 || processingStatusCode == 4 {
			numOfOpen++
			if processingStatusCode == 4 {
				numOfOverdue++
			}
		} else if processingStatusCode == 3 {
			numOfInProgressing++
		}
	}
	// eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"PlanRepairTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}}).All(&eventLatestResults)
	// numOfCompletedToday := len(eventLatestResults)
	// fmt.Println("NumOfCompletedToday -> ", numOfCompletedToday)

	// eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"ProcessingStatusCode": 4}}}).All(&eventLatestResults)
	// numOfOverdue := len(eventLatestResults)
	// fmt.Println("NumOfOverdue -> ", numOfOverdue)

	var row []interface{}
	row = append(row, numOfTPC)
	row = append(row, numOfTotalEvents)
	row = append(row, numOfLatestEvents)
	row = append(row, numOfCompletedToday)
	row = append(row, numOfOverdue)
	row = append(row, numOfOpen)
	row = append(row, numOfInProgressing)
	row = append(row, numOfCompleted)
	rows = append(rows, row)
	// fmt.Println(rows)

	columns := []map[string]string{
		{"text": "NumOfTPC", "type": "number"},
		{"text": "NumOfTotalEvents", "type": "number"},
		{"text": "NumOfLatestEvents", "type": "number"},
		{"text": "CompletedToday", "type": "number"},
		{"text": "Overdue", "type": "number"},
		{"text": "Open", "type": "number"},
		{"text": "InProgressing", "type": "number"},
		{"text": "Completed", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	return grafanaData
}

func machinesSinglestat(groupID string) map[string]interface{} {
	statisticCollection := config.DB.C(config.Statistic)
	var result map[string]interface{}
	var rows []interface{}
	statisticCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"Dashboard": "DO"}}}).One(&result)
	// fmt.Println("Result: ", result)

	var row []interface{}
	row = append(row, result["Total"])
	row = append(row, result["Run"])
	row = append(row, result["Idle"])
	row = append(row, result["Down"])
	row = append(row, result["Off"])
	row = append(row, result["SumFailureCost"])
	rows = append(rows, row)
	// fmt.Println(rows)

	columns := []map[string]string{
		{"text": "Total", "type": "number"},
		{"text": "Run", "type": "number"},
		{"text": "Idle", "type": "number"},
		{"text": "Down", "type": "number"},
		{"text": "Off", "type": "number"},
		{"text": "SumFailureCost", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	return grafanaData
}

func eventHist(groupID string, machineID string) map[string]interface{} {
	eventHistCollection := config.DB.C(config.EventHist)
	rows := []interface{}{}
	var results []map[string]interface{}
	if machineID == "" {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"EventCode": 1}}}).All(&results)
	} else {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}, {"$match": bson.M{"EventCode": 1}}}).All(&results)
	}
	for indexOfResult := 0; indexOfResult < len(results); indexOfResult++ {
		var result map[string]interface{}
		result = results[indexOfResult]
		// fmt.Println("Result: ", result)
		var row []interface{}
		row = append(row, result["EventID"])
		row = append(row, result["EventCode"])
		row = append(row, result["Type"])
		row = append(row, result["GroupID"])
		row = append(row, result["GroupName"])
		row = append(row, result["MachineID"])
		row = append(row, result["MachineName"])
		row = append(row, result["ProcessingStatusCode"])
		row = append(row, result["ProcessingProgress"])
		row = append(row, result["AbnormalStartTime"])
		if result["ProcessingStatusCode"].(int) == 5 {
			if result["RepairStartTime"] != nil {
				row = append(row, result["RepairStartTime"].(time.Time).Sub(result["AbnormalStartTime"].(time.Time)).Seconds())
				row = append(row, result["RepairStartTime"])
			} else {
				row = append(row, result["AbnormalStartTime"].(time.Time).Sub(result["AbnormalStartTime"].(time.Time)).Seconds())
				row = append(row, result["AbnormalStartTime"])
			}
			row = append(row, result["CompleteTime"].(time.Time).Sub(result["RepairStartTime"].(time.Time)).Seconds())
			row = append(row, result["CompleteTime"])
		} else {
			row = append(row, nil)
			row = append(row, nil)
			row = append(row, nil)
			row = append(row, nil)
		}
		row = append(row, result["AbnormalLastingSecond"])
		row = append(row, result["ShouldRepairTime"])
		row = append(row, result["PlanRepairTime"])
		row = append(row, result["PrincipalID"])
		row = append(row, result["PrincipalName"])
		row = append(row, result["AbnormalReason"])
		row = append(row, result["AbnormalSolution"])
		row = append(row, result["AbnormalCode"])
		row = append(row, result["AbnormalPosition"])
		// fmt.Println(row)
		rows = append(rows, row)
	}
	// fmt.Println(rows)

	columns := []map[string]string{
		{"text": "EventID", "type": "string"},
		{"text": "EventCode", "type": "number"},
		{"text": "Type", "type": "string"},
		{"text": "GroupID", "type": "string"},
		{"text": "GroupName", "type": "string"},
		{"text": "MachineID", "type": "string"},
		{"text": "MachineName", "type": "string"},
		{"text": "ProcessingStatusCode", "type": "number"},
		{"text": "ProcessingProgress", "type": "string"},
		{"text": "AbnormalStartTime", "type": "time"},
		{"text": "TTD", "type": "number"},
		{"text": "RepairStartTime", "type": "time"},
		{"text": "TTR", "type": "number"},
		{"text": "CompleteTime", "type": "time"},
		{"text": "AbnormalLastingSecond", "type": "number"},
		{"text": "ShouldRepairTime", "type": "time"},
		{"text": "PlanRepairTime", "type": "time"},
		{"text": "PrincipalID", "type": "string"},
		{"text": "PrincipalName", "type": "string"},
		{"text": "AbnormalReason", "type": "string"},
		{"text": "AbnormalSolution", "type": "string"},
		{"text": "AbnormalCode", "type": "number"},
		{"text": "AbnormalPosition", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	return grafanaData
}

func abnormalReasonRank(groupID string, fromTime time.Time, toTime time.Time) map[string]interface{} {
	rows := []interface{}{}
	eventHistCollection := config.DB.C(config.EventHist)
	var results []map[string]interface{}
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"CompleteTime": bson.M{"$gte": fromTime, "$lt": toTime}}}, {"$group": bson.M{"_id": "$AbnormalReason", "count": bson.M{"$sum": 1}}}, {"$sort": bson.M{"count": -1}}}).All(&results)
	// fmt.Println(results)
	recordCount := 0
	for _, result := range results {
		if result["_id"] != nil && recordCount < 10 {
			// fmt.Println(result)
			row := []interface{}{result["_id"], result["count"]}
			rows = append(rows, row)
			recordCount++
		}
	}
	// fmt.Println(rows)

	columns := []map[string]string{
		{"text": "AbnormalReason", "type": "string"},
		{"text": "Count", "type": "float"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func lastMonthAbReasonRank(groupID string) map[string]interface{} {
	rows := []interface{}{}
	var starttimeT, endtimeT time.Time
	nowTime := time.Now().In(config.TaipeiTimeZone)
	eventHistCollection := config.DB.C(config.EventHist)
	var results []map[string]interface{}
	year := nowTime.Year()
	// month := int(nowTime.Month())
	endtimeFS := fmt.Sprintf("%02d-01-01T00:00:00+08:00", year)
	endtimeT, _ = time.Parse(time.RFC3339, endtimeFS)
	// if month-1 == 0 {
	// 	month = 12
	// 	year = year - 1
	// } else {
	// 	month = month - 1
	// }
	starttimeFS := fmt.Sprintf("%02d-01-01T00:00:00+08:00", year-1)
	starttimeT, _ = time.Parse(time.RFC3339, starttimeFS)
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"RepairStartTime": bson.M{"$gte": starttimeT, "$lt": endtimeT}}}, {"$match": bson.M{"CompleteTime": bson.M{"$gte": starttimeT, "$lt": endtimeT}}}, {"$group": bson.M{"_id": "$AbnormalReason", "count": bson.M{"$sum": 1}}}, {"$sort": bson.M{"count": -1}}}).All(&results)
	// fmt.Println(results)
	recordCount := 0
	for _, result := range results {
		if result["_id"] != nil && recordCount < 10 {
			// fmt.Println(result)
			row := []interface{}{result["_id"], result["count"]}
			rows = append(rows, row)
			recordCount++
		}
	}
	// fmt.Println(rows)

	columns := []map[string]string{
		{"text": "AbnormalReason", "type": "string"},
		{"text": "Count", "type": "float"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func machinesFlowCharting(groupID string, refID string) map[string]interface{} {
	grafanaData := map[string]interface{}{}
	machineRawDataCollection := config.DB.C(config.MachineRawData)
	columns := []map[string]string{}
	rows := []interface{}{}
	row := []interface{}{}

	var machineRawDataResults []map[string]interface{}
	machineRawDataCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&machineRawDataResults)
	for _, machineRawDataResult := range machineRawDataResults {
		if machineRawDataResult["ManualEvent"].(int) > 0 {
			row = append(row, 2000)
		} else {
			row = append(row, machineRawDataResult["StatusLay1Value"])
		}
		columns = append(columns, map[string]string{"text": machineRawDataResult["MachineName"].(string), "type": "string"})
		// fmt.Println(row)
	}
	rows = append(rows, row)
	// fmt.Println(columns, rows)

	grafanaData = map[string]interface{}{
		"refId":   refID,
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	return grafanaData
}

func tpcFlowCharting(groupID string, refID string) map[string]interface{} {
	grafanaData := map[string]interface{}{}
	tpcListCollection := config.DB.C(config.TPCList)
	eventLatestCollection := config.DB.C(config.EventLatest)
	columns := []map[string]string{}
	rows := []interface{}{}
	row := []interface{}{}

	var tpcListResults []map[string]interface{}
	tpcListCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&tpcListResults)
	for _, tpcListResult := range tpcListResults {
		var eventLatestResults []map[string]interface{}
		eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"TPCID": tpcListResult["TPCID"]}}, {"$match": bson.M{"EventCode": 2}}}).All(&eventLatestResults)
		columns = append(columns, map[string]string{"text": tpcListResult["TPCName"].(string), "type": "string"})
		row = append(row, len(eventLatestResults))
		// fmt.Println(row)
	}
	if len(tpcListResults) != 0 {
		rows = append(rows, row)
	} else {
		columns = append(columns, map[string]string{"text": "No TPC", "type": "string"})
	}
	// fmt.Println(columns, rows)

	grafanaData = map[string]interface{}{
		"refId":   refID,
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	return grafanaData
}

func meanTimeComputeSinglestat(computeType string, groupID string, totalComputeValue float64, fromTime time.Time, toTime time.Time) map[string]interface{} {
	eventHistResults := []map[string]interface{}{}
	rows := []interface{}{}
	eventHistCollection := config.DB.C(config.EventHist)
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": fromTime, "$lt": toTime}}}}).All(&eventHistResults)
	totalNumberOfFailures := len(eventHistResults)

	row := []interface{}{totalComputeValue, totalNumberOfFailures}
	rows = append(rows, row)

	columns := []map[string]string{
		{"text": "Group" + computeType, "type": "float"},
		{"text": "NumberOfFailures", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func v2MeanTimeCompute(computeType string, groupID string, fromTime time.Time, toTime time.Time) (map[string]interface{}, float64) {
	var results []map[string]interface{}
	rows := []interface{}{}
	eventHistCollection := config.DB.C(config.EventHist)
	machineRawDataHistCollection := config.DB.C(config.MachineRawDataHist)
	machineRawDataCollection := config.DB.C(config.MachineRawData)
	var machineRawDataResults []map[string]interface{}
	var resultValueArray []float64
	var resultMap = make(map[string]float64)
	// totalStartT = fromTime
	// totalEndT = toTime
	machineRawDataCollection.Find(bson.M{"GroupID": groupID}).All(&machineRawDataResults)
	// fmt.Println("MachineRawDataResults:", machineRawDataResults)

	totalComputeValue := 0.0
	for _, machineRawDataResult := range machineRawDataResults {
		// fmt.Println("MachineID:", machineRawDataResult["MachineID"])
		var sumOfSecond float64
		if computeType == "MTTD" {
			eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"MachineID": machineRawDataResult["MachineID"]}}, {"$match": bson.M{"RepairStartTime": bson.M{"$gte": fromTime, "$lt": toTime}}}}).All(&results)
			for _, result := range results {
				// fmt.Println(result)
				ast := result["AbnormalStartTime"].(time.Time)
				rst := result["RepairStartTime"].(time.Time)
				// fmt.Println("   ", "AST:", ast, " RST:", rst)
				mttd := rst.Sub(ast).Seconds()
				// fmt.Println("       ", "MTTD:", mttd)
				sumOfSecond = sumOfSecond + mttd
			}
		} else if computeType == "MTTR" {
			eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"MachineID": machineRawDataResult["MachineID"]}}, {"$match": bson.M{"CompleteTime": bson.M{"$gte": fromTime, "$lt": toTime}}}}).All(&results)
			// fmt.Println("Results:", results)
			for _, result := range results {
				if result["RepairStartTime"] != nil {
					rst := result["RepairStartTime"].(time.Time)
					rct := result["CompleteTime"].(time.Time)
					// fmt.Println(" RST:", rst, " RCT:", rct)
					mttr := rct.Sub(rst).Seconds()
					// fmt.Println("MTTR:", mttr)
					sumOfSecond = sumOfSecond + mttr
				}
			}
		} else if computeType == "MTBF" {
			eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"ProcessingStatusCode": 5}}, {"$match": bson.M{"MachineID": machineRawDataResult["MachineID"]}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": fromTime, "$lte": toTime}}}}).All(&results)
			// fmt.Println("EventHistCollectionResults:", results)
			var sumOfFailureSeconds float64
			for _, result := range results {
				// fmt.Println("  AbnormalStartTime ->", result["AbnormalStartTime"])
				ast := result["AbnormalStartTime"].(time.Time)
				var lastEvent map[string]interface{}
				eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"ProcessingStatusCode": 5}}, {"$match": bson.M{"MachineID": result["MachineID"]}}, {"$match": bson.M{"CompleteTime": bson.M{"$lte": ast}}}, {"$sort": bson.M{"CompleteTime": -1}}}).One(&lastEvent)
				// fmt.Println("  LastEventCompleteTime ->", lastEvent["CompleteTime"])
				if len(lastEvent) != 0 {
					rct := lastEvent["CompleteTime"].(time.Time)
					// fmt.Println("AST:", ast, " RCT:", rct)
					failureSeconds := ast.Sub(rct).Seconds()
					// fmt.Println("    FailureSeconds:", failureSeconds)
					var machineRawDataHistResults []map[string]interface{}
					machineRawDataHistCollection.Pipe([]bson.M{{"MachineID": result["MachineID"]}, {"$match": bson.M{"Timestamp": bson.M{"$gte": rct, "$lte": ast}}}, {"$sort": bson.M{"Timestamp": 1}}}).All(&machineRawDataHistResults)
					// fmt.Println("    MachineRawDataHistResults:", machineRawDataHistResults)
					var sumOfOffTimeSeconds float64
					for machineRawDataHistResultIndex, machineRawDataHistResult := range machineRawDataHistResults {
						// fmt.Println("      Timestamp:", machineRawDataHistResult["Timestamp"], "StatusLay1Value:", machineRawDataHistResult["StatusLay1Value"])
						if machineRawDataHistResult["StatusLay1Value"] == 4000 {
							if machineRawDataHistResultIndex != len(machineRawDataHistResults)-1 {
								// fmt.Println("        Next -> Timestamp:", machineRawDataHistResults[machineRawDataHistResultIndex+1]["Timestamp"], "StatusLay1Value:", machineRawDataHistResults[machineRawDataHistResultIndex+1]["StatusLay1Value"])
								offTimeSeconds := machineRawDataHistResults[machineRawDataHistResultIndex+1]["Timestamp"].(time.Time).Sub(machineRawDataHistResult["Timestamp"].(time.Time)).Seconds()
								// fmt.Println("        OffTimeSeconds:", offTimeSeconds)
								sumOfOffTimeSeconds += offTimeSeconds
							} else {
								// fmt.Println("        Next -> Timestamp:", ast)
								offTimeSeconds := ast.Sub(machineRawDataHistResult["Timestamp"].(time.Time)).Seconds()
								// fmt.Println("        OffTimeSeconds:", offTimeSeconds)
								sumOfOffTimeSeconds += offTimeSeconds
							}
						}
					}
					// fmt.Println("    SumOfOffTimeSeconds:", sumOfOffTimeSeconds)
					betweenFailureSeconds := failureSeconds - sumOfOffTimeSeconds
					// fmt.Println("  BetweenFailureSeconds :", betweenFailureSeconds)
					// fmt.Println()
					sumOfFailureSeconds = sumOfFailureSeconds + betweenFailureSeconds
				}
			}
			sumOfSecond = sumOfFailureSeconds
			// fmt.Println("SumOfSecond:", sumOfSecond, " NumberOfEvents:", len(results))
		}
		// var row []interface{}
		sumOfHours := sumOfSecond / 3600.0
		if sumOfSecond == 0 {
			resultValueArray = append(resultValueArray, 0.0)
			resultMap[machineRawDataResult["MachineName"].(string)] = 0.0
			// row = []interface{}{machineRawDataResult["MachineName"], nil}
		} else {
			computeValue := (sumOfHours / float64(len(results)))
			resultValueArray = append(resultValueArray, computeValue)
			resultMap[machineRawDataResult["MachineName"].(string)] = computeValue
			// row = []interface{}{machineRawDataResult["MachineName"], (sumOfSecond / float64(len(results)))}
			totalComputeValue += computeValue
		}
		// rows = append(rows, row)
	}
	totalComputeValue = totalComputeValue / float64(len(machineRawDataResults))
	if computeType == "MTBF" {
		sort.Float64s(resultValueArray)
		for _, iV := range resultValueArray {
			for k, v := range resultMap {
				if iV == v {
					row := []interface{}{k, v}
					rows = append(rows, row)
					delete(resultMap, k)
					break
				}
			}
		}
	} else {
		sort.Float64s(resultValueArray)
		sort.Sort(sort.Reverse(sort.Float64Slice(resultValueArray)))
		for _, iV := range resultValueArray {
			for k, v := range resultMap {
				if iV == v {
					row := []interface{}{k, v}
					rows = append(rows, row)
					delete(resultMap, k)
					break
				}
			}
		}
	}
	// fmt.Println(testArray)
	// sort.Float64s(testArray)
	// fmt.Println(testArray)
	// sort.Sort(sort.Reverse(sort.Float64Slice(testArray)))
	// fmt.Println(testArray)

	// for i, j := 0, len(datapoints)-1; i < j; i, j = i+1, j-1 {
	// 	datapoints[i], datapoints[j] = datapoints[j], datapoints[i]
	// }

	columns := []map[string]string{
		{"text": "MachineName", "type": "string"},
		{"text": computeType, "type": "float"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData, totalComputeValue
}

func v3MeanTimeCompute(computeType string, groupID string, fromTime time.Time, toTime time.Time) map[string]interface{} {
	var result map[string]interface{}
	var targetValueArray []float64
	targetMap := map[string]float64{}
	rows := []interface{}{}
	monthlyStatisticsCollection := config.DB.C(config.MonthlyStatistics)
	yearlyStatisticsCollection := config.DB.C(config.YearlyStatistics)

	if fromTime.Month() != toTime.Month() {
		yearlyStatisticsCollection.Pipe([]bson.M{{"$match": bson.M{"DateTime": fromTime}}, {"$match": bson.M{"GroupID": groupID}}}).One(&result)
	} else {
		monthlyStatisticsCollection.Pipe([]bson.M{{"$match": bson.M{"DateTime": fromTime}}, {"$match": bson.M{"GroupID": groupID}}}).One(&result)
	}
	if len(result) != 0 {
		machines := result["Machines"].([]interface{})
		for _, machine := range machines {
			machineMap := machine.(map[string]interface{})
			tempValue := machineMap[computeType].(float64) / 3600.0
			targetMap[machineMap["MachineName"].(string)] = tempValue
			targetValueArray = append(targetValueArray, tempValue)
		}
		nonZeroValueRow, zeroValueRow := []interface{}{}, []interface{}{}
		if computeType == "MTBF" {
			sort.Float64s(targetValueArray)
			for _, iV := range targetValueArray {
				for k, v := range targetMap {
					if iV == v {
						if v != 0 {
							row := []interface{}{k, v}
							nonZeroValueRow = append(nonZeroValueRow, row)
						} else {
							row := []interface{}{k, v}
							zeroValueRow = append(zeroValueRow, row)
						}
						delete(targetMap, k)
						break
					}
				}
			}
		} else {
			sort.Float64s(targetValueArray)
			sort.Sort(sort.Reverse(sort.Float64Slice(targetValueArray)))
			for _, iV := range targetValueArray {
				for k, v := range targetMap {
					if iV == v {
						if v != 0 {
							row := []interface{}{k, v}
							nonZeroValueRow = append(nonZeroValueRow, row)
						} else {
							row := []interface{}{k, v}
							zeroValueRow = append(zeroValueRow, row)
						}
						delete(targetMap, k)
						break
					}
				}
			}
		}
		rows = append(nonZeroValueRow, zeroValueRow...)
	}

	columns := []map[string]string{
		{"text": "MachineName", "type": "string"},
		{"text": computeType, "type": "float"},
	}
	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	return grafanaData
}

// Timeseries
func v1MeanTimeCompute(computeType string, groupID string) map[string]interface{} {
	var starttimeT, endtimeT time.Time
	nowTime := time.Now().In(config.TaipeiTimeZone)
	var results []map[string]interface{}
	var datapoints []interface{}
	eventHistCollection := config.DB.C(config.EventHist)
	machineRawDataHistCollection := config.DB.C(config.MachineRawDataHist)

	year := nowTime.Year()
	month := int(nowTime.Month())
	// for monthIndex := 0; monthIndex < 1; monthIndex++ {
	for monthIndex := 0; monthIndex < 12; monthIndex++ {
		var datapoint []interface{}
		if monthIndex == 0 {
			endtimeT = nowTime
		} else {
			endtimeT = starttimeT
		}
		starttimeFS := fmt.Sprintf("%02d-%02d-01T00:00:00+08:00", year, month)
		starttimeT, _ = time.Parse(time.RFC3339, starttimeFS)
		// fmt.Println("\nStartTime:", starttimeT, "EndTime:", endtimeT)
		var sumOfSecond float64
		if computeType == "MTTD" {
			eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lt": endtimeT}}}, {"$match": bson.M{"RepairStartTime": bson.M{"$gte": starttimeT, "$lt": endtimeT}}}}).All(&results)
			for _, result := range results {
				// fmt.Println(result)
				ast := result["AbnormalStartTime"].(time.Time)
				rst := result["RepairStartTime"].(time.Time)
				// fmt.Println("   ", "AST:", ast, " RST:", rst)
				mttd := rst.Sub(ast).Seconds()
				// fmt.Println("       ", "MTTD:", mttd)
				sumOfSecond = sumOfSecond + mttd
			}
		} else if computeType == "MTTR" {
			eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"RepairStartTime": bson.M{"$gte": starttimeT, "$lt": endtimeT}}}, {"$match": bson.M{"CompleteTime": bson.M{"$gte": starttimeT, "$lt": endtimeT}}}}).All(&results)
			// fmt.Println("Results:", results)
			for _, result := range results {
				rst := result["RepairStartTime"].(time.Time)
				rct := result["CompleteTime"].(time.Time)
				// fmt.Println(" RST:", rst, " RCT:", rct)
				mttr := rct.Sub(rst).Seconds()
				// fmt.Println("MTTR:", mttr)
				sumOfSecond = sumOfSecond + mttr
			}
		} else if computeType == "MTBF" {
			eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"ProcessingStatusCode": 5}}, {"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}}).All(&results)
			// fmt.Println("EventHistCollectionResults:", results)
			var sumOfFailureSeconds float64
			for _, result := range results {
				// fmt.Println("  AbnormalStartTime ->", result["AbnormalStartTime"])
				ast := result["AbnormalStartTime"].(time.Time)
				var lastEvent map[string]interface{}
				eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"EventCode": 1}}, {"$match": bson.M{"ProcessingStatusCode": 5}}, {"$match": bson.M{"MachineID": result["MachineID"]}}, {"$match": bson.M{"CompleteTime": bson.M{"$lte": ast}}}, {"$sort": bson.M{"CompleteTime": -1}}}).One(&lastEvent)
				// fmt.Println("  LastEventCompleteTime ->", lastEvent["CompleteTime"])
				if len(lastEvent) != 0 {
					rct := lastEvent["CompleteTime"].(time.Time)
					// fmt.Println("AST:", ast, " RCT:", rct)
					failureSeconds := ast.Sub(rct).Seconds()
					// fmt.Println("    FailureSeconds:", failureSeconds)
					var machineRawDataHistResults []map[string]interface{}
					machineRawDataHistCollection.Pipe([]bson.M{{"MachineID": result["MachineID"]}, {"$match": bson.M{"Timestamp": bson.M{"$gte": rct, "$lte": ast}}}, {"$sort": bson.M{"Timestamp": 1}}}).All(&machineRawDataHistResults)
					// fmt.Println("    MachineRawDataHistResults:", machineRawDataHistResults)
					var sumOfOffTimeSeconds float64
					for machineRawDataHistResultIndex, machineRawDataHistResult := range machineRawDataHistResults {
						// fmt.Println("      Timestamp:", machineRawDataHistResult["Timestamp"], "StatusLay1Value:", machineRawDataHistResult["StatusLay1Value"])
						if machineRawDataHistResult["StatusLay1Value"] == 4000 {
							if machineRawDataHistResultIndex != len(machineRawDataHistResults)-1 {
								// fmt.Println("        Next -> Timestamp:", machineRawDataHistResults[machineRawDataHistResultIndex+1]["Timestamp"], "StatusLay1Value:", machineRawDataHistResults[machineRawDataHistResultIndex+1]["StatusLay1Value"])
								offTimeSeconds := machineRawDataHistResults[machineRawDataHistResultIndex+1]["Timestamp"].(time.Time).Sub(machineRawDataHistResult["Timestamp"].(time.Time)).Seconds()
								// fmt.Println("        OffTimeSeconds:", offTimeSeconds)
								sumOfOffTimeSeconds += offTimeSeconds
							} else {
								// fmt.Println("        Next -> Timestamp:", ast)
								offTimeSeconds := ast.Sub(machineRawDataHistResult["Timestamp"].(time.Time)).Seconds()
								// fmt.Println("        OffTimeSeconds:", offTimeSeconds)
								sumOfOffTimeSeconds += offTimeSeconds
							}
						}
					}
					// fmt.Println("    SumOfOffTimeSeconds:", sumOfOffTimeSeconds)
					betweenFailureSeconds := failureSeconds - sumOfOffTimeSeconds
					// fmt.Println("  BetweenFailureSeconds :", betweenFailureSeconds)
					// fmt.Println()
					sumOfFailureSeconds = sumOfFailureSeconds + betweenFailureSeconds
				}
			}
			sumOfFailureHours := sumOfFailureSeconds / 3600.0
			sumOfSecond = sumOfFailureHours
			// fmt.Println("SumOfSecond:", sumOfSecond, " NumberOfEvents:", len(results))
		}

		if sumOfSecond == 0 {
			// datapoint = []interface{}{0.0, starttimeT.Unix() * 1000}
			datapoint = []interface{}{nil, starttimeT.Unix() * 1000}
		} else {
			datapoint = []interface{}{(sumOfSecond / float64(len(results))), starttimeT.Unix() * 1000}
		}

		datapoints = append(datapoints, datapoint)

		if (month - 1) == 0 {
			month = 12
			year = year - 1
		} else {
			month = month - 1
		}
	}

	for i, j := 0, len(datapoints)-1; i < j; i, j = i+1, j-1 {
		datapoints[i], datapoints[j] = datapoints[j], datapoints[i]
	}

	grafanaData := map[string]interface{}{
		"target":     computeType,
		"datapoints": datapoints,
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func factoryMapTimeseries(groupID string) []map[string]interface{} {
	grafnaResponseArray := []map[string]interface{}{}
	machineRawDataCollection := config.DB.C(config.MachineRawData)
	tpcListCollection := config.DB.C(config.TPCList)
	eventLatestCollection := config.DB.C(config.EventLatest)

	var tpcListResults []map[string]interface{}
	tpcListCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&tpcListResults)

	for _, tpcListResult := range tpcListResults {
		var datapoints []interface{}
		var datapoint []interface{}

		var eventLatestResults []map[string]interface{}
		eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"TPCID": tpcListResult["TPCID"]}}, {"$match": bson.M{"EventCode": 2}}}).All(&eventLatestResults)

		datapoint = []interface{}{len(eventLatestResults), time.Now().In(config.TaipeiTimeZone).Unix() * 1000}
		datapoints = append(datapoints, datapoint)
		grafanaData := map[string]interface{}{
			"target":     tpcListResult["TPCName"],
			"datapoints": datapoints,
		}
		grafnaResponseArray = append(grafnaResponseArray, grafanaData)
	}

	var machineRawDataResults []map[string]interface{}
	machineRawDataCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&machineRawDataResults)

	for indexOfMachineRawDataResult := 0; indexOfMachineRawDataResult < len(machineRawDataResults); indexOfMachineRawDataResult++ {
		var datapoints []interface{}
		var datapoint []interface{}
		var result map[string]interface{}
		result = machineRawDataResults[indexOfMachineRawDataResult]
		// fmt.Println("Result: ", result)
		if result["ManualEvent"].(int) > 0 {
			datapoint = []interface{}{2000, result["Timestamp"].(time.Time).Unix() * 1000}
		} else {
			datapoint = []interface{}{result["StatusLay1Value"], result["Timestamp"].(time.Time).Unix() * 1000}
		}
		// fmt.Println(row)
		datapoints = append(datapoints, datapoint)
		grafanaData := map[string]interface{}{
			"target":     result["MachineName"],
			"datapoints": datapoints,
		}
		grafnaResponseArray = append(grafnaResponseArray, grafanaData)
	}
	return grafnaResponseArray
}

//Test ...
// func Test(w http.ResponseWriter, r *http.Request) {
