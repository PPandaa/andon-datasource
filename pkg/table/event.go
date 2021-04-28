package table

import (
	"DataSource/config"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func EventLatest(groupID string, machineID string) map[string]interface{} {
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

func EventList(groupID string, startTime time.Time, endTime time.Time) map[string]interface{} {
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

func EventSinglestat(groupID string, machineID string) map[string]interface{} {
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

func EventHistTable(groupID string, machineID string, startTime time.Time, endTime time.Time) map[string]interface{} {
	rows := []interface{}{}
	eventLatestcollection := config.DB.C(config.EventLatest)
	var eventLatestResults []map[string]interface{}
	if machineID == "" {
		eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventLatestResults)
	} else {
		eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventLatestResults)
	}
	for indexOfResult := 0; indexOfResult < len(eventLatestResults); indexOfResult++ {
		var result map[string]interface{}
		result = eventLatestResults[indexOfResult]
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
		if result["ProcessingStatusCode"].(int) >= 3 {
			if result["RepairStartTime"] != nil {
				row = append(row, result["RepairStartTime"].(time.Time).Sub(result["AbnormalStartTime"].(time.Time)).Seconds())
				row = append(row, result["RepairStartTime"])
			} else {
				row = append(row, result["AbnormalStartTime"].(time.Time).Sub(result["AbnormalStartTime"].(time.Time)).Seconds())
				row = append(row, result["AbnormalStartTime"])
			}
			if result["ProcessingStatusCode"].(int) == 5 {
				row = append(row, result["CompleteTime"].(time.Time).Sub(result["RepairStartTime"].(time.Time)).Seconds())
				row = append(row, result["CompleteTime"])
			} else {
				row = append(row, nil)
				row = append(row, nil)
			}
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

	eventHistCollection := config.DB.C(config.EventHist)

	var eventHistResults []map[string]interface{}
	if machineID == "" {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventHistResults)
	} else {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventHistResults)
	}
	for indexOfResult := 0; indexOfResult < len(eventHistResults); indexOfResult++ {
		var result map[string]interface{}
		result = eventHistResults[indexOfResult]
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

func EventHistSinglestat(groupID string, machineID string, startTime time.Time, endTime time.Time) map[string]interface{} {

	eventLatestcollection := config.DB.C(config.EventLatest)
	var eventLatestResults []map[string]interface{}
	if machineID == "" {
		eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventLatestResults)
	} else {
		eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventLatestResults)
	}

	eventHistCollection := config.DB.C(config.EventHist)
	var eventHistResults []map[string]interface{}
	if machineID == "" {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventHistResults)
	} else {
		eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"MachineID": machineID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": startTime, "$lt": endTime}}}}).All(&eventHistResults)
	}

	columns := []map[string]string{
		{"text": "histEventCount", "type": "number"},
	}

	rows := []interface{}{
		[]int{len(eventLatestResults) + len(eventHistResults)},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}
