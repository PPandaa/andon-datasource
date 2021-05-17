package eventWorkOrderOverview

import (
	"DataSource/config"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func V1_EventTable(groupID string, machineID string, startTime time.Time, endTime time.Time) map[string]interface{} {
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
		// fmt.Println("LatestResult: ", result)
		var row []interface{}
		row = append(row, result["EventID"])
		row = append(row, result["EventCode"])
		row = append(row, result["EventName"])
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
		// fmt.Println("LatestRow: ", row)
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
		// fmt.Println("HistResult: ", result)
		var row []interface{}
		row = append(row, result["EventID"])
		row = append(row, result["EventCode"])
		row = append(row, result["EventName"])
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
				row = append(row, result["CompleteTime"].(time.Time).Sub(result["RepairStartTime"].(time.Time)).Seconds())
			} else {
				row = append(row, result["AbnormalStartTime"].(time.Time).Sub(result["AbnormalStartTime"].(time.Time)).Seconds())
				row = append(row, result["AbnormalStartTime"])
				row = append(row, result["CompleteTime"].(time.Time).Sub(result["AbnormalStartTime"].(time.Time)).Seconds())
			}
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
		// fmt.Println("HistRow: ", row)
		rows = append(rows, row)
	}

	// fmt.Println("Rows: ", rows)

	columns := []map[string]string{
		{"text": "EventID", "type": "string"},
		{"text": "EventCode", "type": "number"},
		{"text": "EventName", "type": "string"},
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
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_EventSinglestat(groupID string, machineID string, startTime time.Time, endTime time.Time) map[string]interface{} {

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
