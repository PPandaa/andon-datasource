package monitoringCenter

import (
	"DataSource/config"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func V1_Panel1Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	// nowTime := time.Now().In(config.TaipeiTimeZone)
	// starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	// starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	// endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	// endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)

	// var eventLatestCallRepairResults []map[string]interface{}
	// eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$match": bson.M{"EventCode": 2}}}).All(&eventLatestCallRepairResults)
	var eventLatestCallRepairResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestCallRepairResults)

	columns := []map[string]string{
		{"text": "total", "type": "number"},
	}

	rows := []interface{}{
		[]int{len(eventLatestCallRepairResults)},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel1Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestCallRepairResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestCallRepairResults)
	rows := []interface{}{}
	for _, eventLatestCallRepairResult := range eventLatestCallRepairResults {
		var row []interface{}
		if eventLatestCallRepairResult["MachineName"] != nil {
			row = []interface{}{eventLatestCallRepairResult["MachineName"], eventLatestCallRepairResult["AbnormalLastingSecond"]}
		} else {
			row = []interface{}{eventLatestCallRepairResult["TPCName"], eventLatestCallRepairResult["AbnormalLastingSecond"]}
		}
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "name", "type": "string"},
		{"text": "abnormalLastingSecond", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel3(groupID string) map[string]interface{} {
	var new, ass, ip, od int
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$group": bson.M{"_id": "$ProcessingStatusCode", "count": bson.M{"$sum": 1}}}}).All(&eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		switch eventLatestResult["_id"].(int) {
		case 0:
			new = eventLatestResult["count"].(int)
			break
		case 1:
			ass = eventLatestResult["count"].(int)
			break
		case 3:
			ip = eventLatestResult["count"].(int)
			break
		case 4:
			od = eventLatestResult["count"].(int)
			break
		}
	}

	columns := []map[string]string{
		{"text": "new", "type": "number"},
		{"text": "ass", "type": "number"},
		{"text": "ip", "type": "number"},
		{"text": "od", "type": "number"},
	}

	rows := []interface{}{
		[]int{new, ass, ip, od},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel5Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventWithPrincipalCount := 0

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["PrincipalID"] != nil {
			eventWithPrincipalCount += 1
		}
	}

	rows := []interface{}{
		[]int{eventWithPrincipalCount},
	}

	columns := []map[string]string{
		{"text": "eventWithPrincipalCount", "type": "number"},
	}

	// rows := []interface{}{
	// 	[]interface{}{"Station-004", 15},
	// 	[]interface{}{"Station-001", 12},
	// 	[]interface{}{"Station-002", 8},
	// 	[]interface{}{"Station-003", 4},
	// 	[]interface{}{"Station-005", 3},
	// }

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel5Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	userListCollection := config.DB.C(config.UserList)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$group": bson.M{"_id": "$PrincipalID", "count": bson.M{"$sum": 1}}}}).All(&eventLatestResults)

	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["_id"] != nil {
			var userListResult map[string]interface{}
			userListCollection.Find(bson.M{"PrincipalID": eventLatestResult["_id"]}).One(&userListResult)
			row := []interface{}{userListResult["PrincipalName"], eventLatestResult["count"]}
			rows = append(rows, row)
		}
	}

	columns := []map[string]string{
		{"text": "principalName", "type": "string"},
		{"text": "times", "type": "number"},
	}

	// rows := []interface{}{
	// 	[]interface{}{"Station-004", 15},
	// 	[]interface{}{"Station-001", 12},
	// 	[]interface{}{"Station-002", 8},
	// 	[]interface{}{"Station-003", 4},
	// 	[]interface{}{"Station-005", 3},
	// }

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel6Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"ProcessingStatusCode": 4}}}).All(&eventLatestResults)

	rows := []interface{}{
		[]int{len(eventLatestResults)},
	}

	columns := []map[string]string{
		{"text": "overdueCount", "type": "number"},
	}

	// rows := []interface{}{
	// 	[]interface{}{"TPC-001", 17},
	// 	[]interface{}{"TPC-002", 7},
	// }

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel6Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"ProcessingStatusCode": 4}}}).All(&eventLatestResults)

	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		var row []interface{}
		if eventLatestResult["MachineName"] != nil {
			row = []interface{}{eventLatestResult["MachineName"], eventLatestResult["AbnormalLastingSecond"]}
		} else {
			row = []interface{}{eventLatestResult["TPCName"], eventLatestResult["AbnormalLastingSecond"]}
		}

		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "machineName", "type": "string"},
		{"text": "abnormalLastingSecond", "type": "number"},
	}

	// rows := []interface{}{
	// 	[]interface{}{"TPC-001", 17},
	// 	[]interface{}{"TPC-002", 7},
	// }

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel8(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventHistCollection := config.DB.C(config.EventHist)
	tpcListCollection := config.DB.C(config.TPCList)

	nowTime := time.Now().In(config.TaipeiTimeZone)
	starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)

	totalTimes := map[string]int{}

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$group": bson.M{"_id": "$TPCID", "count": bson.M{"$sum": 1}}}}).All(&eventLatestResults)
	// fmt.Println("EventLatestResults:", eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["_id"] != nil {
			totalTimes[eventLatestResult["_id"].(string)] += eventLatestResult["count"].(int)
		}
	}

	var eventHistResults []map[string]interface{}
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$group": bson.M{"_id": "$TPCID", "count": bson.M{"$sum": 1}}}}).All(&eventHistResults)
	// fmt.Println("EventHistResults:", eventHistResults)
	for _, eventHistResult := range eventHistResults {
		if eventHistResult["_id"] != nil {
			totalTimes[eventHistResult["_id"].(string)] += eventHistResult["count"].(int)
		}
	}
	// fmt.Println("TotalTimes:", totalTimes)
	rows := []interface{}{}
	for k, v := range totalTimes {
		var tpcListResult map[string]interface{}
		tpcListCollection.Find(bson.M{"TPCID": k}).One(&tpcListResult)
		row := []interface{}{tpcListResult["TPCName"], v}
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "tpcName", "type": "string"},
		{"text": "times", "type": "number"},
	}

	// rows := []interface{}{
	// 	[]interface{}{"Station-004", 15},
	// 	[]interface{}{"Station-001", 12},
	// 	[]interface{}{"Station-002", 8},
	// 	[]interface{}{"Station-003", 4},
	// 	[]interface{}{"Station-005", 3},
	// }

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println("GrafanaData:", grafanaData)
	return grafanaData
}

func V1_Panel9(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventHistCollection := config.DB.C(config.EventHist)
	eventConfigCollection := config.DB.C(config.EventConfig)

	nowTime := time.Now().In(config.TaipeiTimeZone)
	starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)

	eventTotalTimes := map[int]int{}

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$group": bson.M{"_id": "$EventCode", "count": bson.M{"$sum": 1}}}}).All(&eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["_id"] != nil {
			eventTotalTimes[eventLatestResult["_id"].(int)] += eventLatestResult["count"].(int)
		}
	}

	var eventHistResults []map[string]interface{}
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$group": bson.M{"_id": "$EventCode", "count": bson.M{"$sum": 1}}}}).All(&eventHistResults)
	for _, eventHistResult := range eventHistResults {
		if eventHistResult["_id"] != nil {
			eventTotalTimes[eventHistResult["_id"].(int)] += eventHistResult["count"].(int)
		}
	}

	rows := []interface{}{}
	for k, v := range eventTotalTimes {
		var eventConfigResult map[string]interface{}
		eventConfigCollection.Find(bson.M{"EventCode": k}).One(&eventConfigResult)
		row := []interface{}{eventConfigResult["EventName"], v}
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "eventName", "type": "string"},
		{"text": "times", "type": "number"},
	}

	// rows := []interface{}{
	// 	[]interface{}{"TPC-001", 17},
	// 	[]interface{}{"TPC-002", 7},
	// }

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V1_Panel10(groupID string) map[string]interface{} {
	eventLatestcollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)

	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		var row []interface{}
		row = append(row, eventLatestResult["EventID"])
		row = append(row, eventLatestResult["EventCode"])
		row = append(row, eventLatestResult["EventName"])
		row = append(row, eventLatestResult["Type"])
		row = append(row, eventLatestResult["GroupID"])
		row = append(row, eventLatestResult["GroupName"])
		row = append(row, eventLatestResult["MachineID"])
		row = append(row, eventLatestResult["MachineName"])
		row = append(row, eventLatestResult["TPCID"])
		row = append(row, eventLatestResult["TPCName"])
		row = append(row, eventLatestResult["ProcessingStatusCode"])
		row = append(row, eventLatestResult["ProcessingProgress"])
		row = append(row, eventLatestResult["PrincipalID"])
		row = append(row, eventLatestResult["PrincipalName"])
		row = append(row, eventLatestResult["AbnormalLastingSecond"])
		row = append(row, eventLatestResult["AbnormalStartTime"])
		row = append(row, eventLatestResult["PlanRepairTime"])
		row = append(row, eventLatestResult["ShouldRepairTime"])
		row = append(row, eventLatestResult["RepairStartTime"])
		row = append(row, eventLatestResult["CompleteTime"])
		row = append(row, eventLatestResult["AbnormalReason"])
		row = append(row, eventLatestResult["AbnormalSolution"])
		row = append(row, eventLatestResult["AbnormalCode"])
		row = append(row, eventLatestResult["AbnormalPosition"])
		// fmt.Println(row)
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "EventID", "type": "string"},
		{"text": "EventCode", "type": "number"},
		{"text": "EventName", "type": "string"},
		{"text": "Type", "type": "string"},
		{"text": "GroupID", "type": "string"},
		{"text": "GroupName", "type": "string"},
		{"text": "MachineID", "type": "string"},
		{"text": "MachineName", "type": "string"},
		{"text": "TPCID", "type": "string"},
		{"text": "TPCName", "type": "string"},
		{"text": "ProcessingStatusCode", "type": "number"},
		{"text": "ProcessingProgress", "type": "string"},
		{"text": "PrincipalID", "type": "string"},
		{"text": "PrincipalName", "type": "string"},
		{"text": "AbnormalLastingSecond", "type": "number"},
		{"text": "AbnormalStartTime", "type": "time"},
		{"text": "PlanRepairTime", "type": "time"},
		{"text": "ShouldRepairTime", "type": "time"},
		{"text": "RepairStartTime", "type": "time"},
		{"text": "CompleteTime", "type": "time"},
		{"text": "AbnormalReason", "type": "string"},
		{"text": "AbnormalSolution", "type": "string"},
		{"text": "AbnormalCode", "type": "number"},
		{"text": "AbnormalPosition", "type": "string"},
	}

	// columns := []map[string]string{
	// 	{"text": "事件類別", "type": "string"},
	// 	{"text": "類型", "type": "string"},
	// 	{"text": "位置", "type": "string"},
	// 	{"text": "狀態", "type": "number"},
	// 	{"text": "負責人", "type": "string"},
	// 	{"text": "持續時間", "type": "number"},
	// 	{"text": "發生時間", "type": "time"},
	// 	{"text": "排修時間", "type": "time"},
	// 	{"text": "執行時間", "type": "time"},
	// 	{"text": "完成時間", "type": "time"},
	// }
	// formatStartTime1, _ := time.Parse(time.RFC3339, "2021-03-19T18:00:00+08:00")
	// formatStartTime2, _ := time.Parse(time.RFC3339, "2021-03-26T12:00:00+08:00")
	// formatPlanTime1, _ := time.Parse(time.RFC3339, "2021-03-26T12:00:00+08:00")
	// formatRepairTime1, _ := time.Parse(time.RFC3339, "2021-03-26T12:00:00+08:00")
	// formatCompleteTime1, _ := time.Parse(time.RFC3339, "2021-03-26T12:00:00+08:00")
	// rows := []interface{}{
	// 	[]interface{}{"Call Repair", "Manual", "Station-004", 3, "黃伯依", 1200, formatStartTime1, formatPlanTime1, nil, nil},
	// 	[]interface{}{"Call Repair", "Auto", "Station-001", 4, "林宏冰", 259200, formatStartTime2, formatPlanTime1, formatRepairTime1, nil},
	// 	[]interface{}{"Call Repair", "Auto", "Station-002", 1, "郭育廷", 300, formatStartTime2, formatPlanTime1, formatRepairTime1, nil},
	// 	[]interface{}{"Call Repair", "Manual", "TPC-001", 0, "張雯婷", 120, formatStartTime2, nil, nil, nil},
	// 	[]interface{}{"Call Repair", "Manual", "TPC-002", 6, "陳雅惠", nil, formatStartTime2, nil, nil, formatCompleteTime1},
	// }

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel1(groupID string) map[string]interface{} {
	var new, ass, ip, od int
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$group": bson.M{"_id": "$ProcessingStatusCode", "count": bson.M{"$sum": 1}}}}).All(&eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		switch eventLatestResult["_id"].(int) {
		case 0:
			new = eventLatestResult["count"].(int)
			break
		case 1:
			ass = eventLatestResult["count"].(int)
			break
		case 3:
			ip = eventLatestResult["count"].(int)
			break
		case 4:
			od = eventLatestResult["count"].(int)
			break
		}
	}

	columns := []map[string]string{
		{"text": "new", "type": "number"},
		{"text": "ass", "type": "number"},
		{"text": "ip", "type": "number"},
		{"text": "od", "type": "number"},
	}

	rows := []interface{}{
		[]int{new, ass, ip, od},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel2Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	// nowTime := time.Now().In(config.TaipeiTimeZone)
	// starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	// starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	// endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	// endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)

	// var eventLatestCallRepairResults []map[string]interface{}
	// eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$match": bson.M{"EventCode": 2}}}).All(&eventLatestCallRepairResults)
	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)

	columns := []map[string]string{
		{"text": "total", "type": "number"},
	}

	rows := []interface{}{
		[]int{len(eventLatestResults)},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel2Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)
	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		var row []interface{}
		if eventLatestResult["MachineName"] != nil && len(eventLatestResult["MachineName"].(string)) != 0 {
			row = []interface{}{eventLatestResult["MachineName"], eventLatestResult["AbnormalLastingSecond"]}
		} else {
			row = []interface{}{eventLatestResult["TPCName"], eventLatestResult["AbnormalLastingSecond"]}
		}
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "name", "type": "string"},
		{"text": "abnormalLastingSecond", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel3Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"ProcessingStatusCode": 0}}}).All(&eventLatestResults)

	rows := []interface{}{
		[]int{len(eventLatestResults)},
	}

	columns := []map[string]string{
		{"text": "newEventCount", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel3Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"ProcessingStatusCode": 0}}}).All(&eventLatestResults)

	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		var row []interface{}
		if eventLatestResult["MachineName"] != nil && len(eventLatestResult["MachineName"].(string)) != 0 {
			row = []interface{}{eventLatestResult["MachineName"], eventLatestResult["AbnormalLastingSecond"]}
		} else {
			row = []interface{}{eventLatestResult["TPCName"], eventLatestResult["AbnormalLastingSecond"]}
		}
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "name", "type": "string"},
		{"text": "times", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel4Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"ProcessingStatusCode": 4}}}).All(&eventLatestResults)

	rows := []interface{}{
		[]int{len(eventLatestResults)},
	}

	columns := []map[string]string{
		{"text": "overdueCount", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel4Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"ProcessingStatusCode": 4}}}).All(&eventLatestResults)

	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		var row []interface{}
		if eventLatestResult["MachineName"] != nil && len(eventLatestResult["MachineName"].(string)) != 0 {
			row = []interface{}{eventLatestResult["MachineName"], eventLatestResult["AbnormalLastingSecond"]}
		} else {
			row = []interface{}{eventLatestResult["TPCName"], eventLatestResult["AbnormalLastingSecond"]}
		}

		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "machineName", "type": "string"},
		{"text": "abnormalLastingSecond", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel5Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventWithPrincipalCount := 0

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["PrincipalID"] != nil && len(eventLatestResult["PrincipalID"].(string)) != 0 {
			eventWithPrincipalCount += 1
		}
	}

	rows := []interface{}{
		[]int{eventWithPrincipalCount},
	}

	columns := []map[string]string{
		{"text": "eventWithPrincipalCount", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel5Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	userListCollection := config.DB.C(config.UserList)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$group": bson.M{"_id": "$PrincipalID", "count": bson.M{"$sum": 1}}}}).All(&eventLatestResults)

	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["_id"] != nil {
			var userListResult map[string]interface{}
			userListCollection.Find(bson.M{"PrincipalID": eventLatestResult["_id"]}).One(&userListResult)
			row := []interface{}{userListResult["PrincipalName"], eventLatestResult["count"]}
			rows = append(rows, row)
		}
	}

	columns := []map[string]string{
		{"text": "principalName", "type": "string"},
		{"text": "times", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel6Singlestat(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventHistCollection := config.DB.C(config.EventHist)

	nowTime := time.Now().In(config.TaipeiTimeZone)
	starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}}).All(&eventLatestResults)

	var eventHistResults []map[string]interface{}
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}}).All(&eventHistResults)

	rows := []interface{}{
		[]int{len(eventLatestResults) + len(eventHistResults)},
	}

	columns := []map[string]string{
		{"text": "todayEventCount", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel6Table(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventHistCollection := config.DB.C(config.EventHist)
	eventConfigCollection := config.DB.C(config.EventConfig)

	nowTime := time.Now().In(config.TaipeiTimeZone)
	starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)

	eventTotalTimes := map[int]int{}

	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$group": bson.M{"_id": "$EventCode", "count": bson.M{"$sum": 1}}}}).All(&eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["_id"] != nil {
			eventTotalTimes[eventLatestResult["_id"].(int)] += eventLatestResult["count"].(int)
		}
	}

	var eventHistResults []map[string]interface{}
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}, {"$group": bson.M{"_id": "$EventCode", "count": bson.M{"$sum": 1}}}}).All(&eventHistResults)
	for _, eventHistResult := range eventHistResults {
		if eventHistResult["_id"] != nil {
			eventTotalTimes[eventHistResult["_id"].(int)] += eventHistResult["count"].(int)
		}
	}

	rows := []interface{}{}
	for k, v := range eventTotalTimes {
		var eventConfigResult map[string]interface{}
		eventConfigCollection.Find(bson.M{"EventCode": k}).One(&eventConfigResult)
		row := []interface{}{eventConfigResult["EventName"], v}
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "eventName", "type": "string"},
		{"text": "todayTimes", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}

func V2_Panel8(groupID string) map[string]interface{} {
	eventLatestCollection := config.DB.C(config.EventLatest)
	eventHistCollection := config.DB.C(config.EventHist)
	machineRawDataCollection := config.DB.C(config.MachineRawData)

	nowTime := time.Now().In(config.TaipeiTimeZone)
	starttimeFS := fmt.Sprintf("%02d-%02d-%02dT00:00:00+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	starttimeT, _ := time.Parse(time.RFC3339, starttimeFS)
	endtimeFS := fmt.Sprintf("%02d-%02d-%02dT23:59:59+08:00", nowTime.Year(), nowTime.Month(), nowTime.Day())
	endtimeT, _ := time.Parse(time.RFC3339, endtimeFS)

	totalTimes := map[string]int{}
	totalLastingSeconds := map[string]float64{}
	var eventLatestResults []map[string]interface{}
	eventLatestCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}}).All(&eventLatestResults)
	// fmt.Println("EventLatestResults:", eventLatestResults)
	for _, eventLatestResult := range eventLatestResults {
		if eventLatestResult["MachineID"] != nil && len(eventLatestResult["MachineID"].(string)) != 0 {
			totalTimes[eventLatestResult["MachineID"].(string)] += 1
			totalLastingSeconds[eventLatestResult["MachineID"].(string)] += eventLatestResult["AbnormalLastingSecond"].(float64)
		}
	}

	var eventHistResults []map[string]interface{}
	eventHistCollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"AbnormalStartTime": bson.M{"$gte": starttimeT, "$lte": endtimeT}}}}).All(&eventHistResults)
	// fmt.Println("EventHistResults:", eventHistResults)
	for _, eventHistResult := range eventHistResults {
		if eventHistResult["MachineID"] != nil && len(eventHistResult["MachineID"].(string)) != 0 {
			totalTimes[eventHistResult["MachineID"].(string)] += 1
			totalLastingSeconds[eventHistResult["MachineID"].(string)] += eventHistResult["AbnormalLastingSecond"].(float64)
		}
	}
	// fmt.Println("TotalTimes:", totalTimes)
	rows := []interface{}{}
	for k, v := range totalTimes {
		var machineRawDataResult map[string]interface{}
		machineRawDataCollection.Find(bson.M{"MachineID": k}).One(&machineRawDataResult)
		row := []interface{}{machineRawDataResult["MachineName"], totalLastingSeconds[k], v}
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "machineName", "type": "string"},
		{"text": "abnormalLastingSecond", "type": "number"},
		{"text": "times", "type": "number"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println("GrafanaData:", grafanaData)
	return grafanaData
}

func V2_Panel9(groupID string) map[string]interface{} {
	eventLatestcollection := config.DB.C(config.EventLatest)

	var eventLatestResults []map[string]interface{}
	eventLatestcollection.Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}}).All(&eventLatestResults)

	rows := []interface{}{}
	for _, eventLatestResult := range eventLatestResults {
		var row []interface{}
		row = append(row, eventLatestResult["EventID"])
		row = append(row, eventLatestResult["EventCode"])
		row = append(row, eventLatestResult["EventName"])
		row = append(row, eventLatestResult["Type"])
		row = append(row, eventLatestResult["GroupID"])
		row = append(row, eventLatestResult["GroupName"])
		row = append(row, eventLatestResult["MachineID"])
		row = append(row, eventLatestResult["MachineName"])
		row = append(row, eventLatestResult["TPCID"])
		row = append(row, eventLatestResult["TPCName"])
		row = append(row, eventLatestResult["ProcessingStatusCode"])
		row = append(row, eventLatestResult["ProcessingProgress"])
		row = append(row, eventLatestResult["PrincipalID"])
		row = append(row, eventLatestResult["PrincipalName"])
		row = append(row, eventLatestResult["AbnormalLastingSecond"])
		row = append(row, eventLatestResult["AbnormalStartTime"])
		row = append(row, eventLatestResult["PlanRepairTime"])
		row = append(row, eventLatestResult["ShouldRepairTime"])
		row = append(row, eventLatestResult["RepairStartTime"])
		row = append(row, eventLatestResult["CompleteTime"])
		row = append(row, eventLatestResult["AbnormalReason"])
		row = append(row, eventLatestResult["AbnormalSolution"])
		row = append(row, eventLatestResult["AbnormalCode"])
		row = append(row, eventLatestResult["AbnormalPosition"])

		row = append(row, eventLatestResult["Parameters"].(map[string]interface{})["WIP_NO"])
		row = append(row, eventLatestResult["Parameters"].(map[string]interface{})["UNIT_NO"])
		row = append(row, eventLatestResult["Parameters"].(map[string]interface{})["PLAN_QTY"])
		row = append(row, eventLatestResult["Parameters"].(map[string]interface{})["MODEL_NO"])
		row = append(row, eventLatestResult["Parameters"].(map[string]interface{})["LINE_DESC"])
		row = append(row, eventLatestResult["Parameters"].(map[string]interface{})["ITEM_NO"])
		// fmt.Println(row)
		rows = append(rows, row)
	}

	columns := []map[string]string{
		{"text": "EventID", "type": "string"},
		{"text": "EventCode", "type": "number"},
		{"text": "EventName", "type": "string"},
		{"text": "Type", "type": "string"},
		{"text": "GroupID", "type": "string"},
		{"text": "GroupName", "type": "string"},
		{"text": "MachineID", "type": "string"},
		{"text": "MachineName", "type": "string"},
		{"text": "TPCID", "type": "string"},
		{"text": "TPCName", "type": "string"},
		{"text": "ProcessingStatusCode", "type": "number"},
		{"text": "ProcessingProgress", "type": "string"},
		{"text": "PrincipalID", "type": "string"},
		{"text": "PrincipalName", "type": "string"},
		{"text": "AbnormalLastingSecond", "type": "number"},
		{"text": "AbnormalStartTime", "type": "time"},
		{"text": "PlanRepairTime", "type": "time"},
		{"text": "ShouldRepairTime", "type": "time"},
		{"text": "RepairStartTime", "type": "time"},
		{"text": "CompleteTime", "type": "time"},
		{"text": "AbnormalReason", "type": "string"},
		{"text": "AbnormalSolution", "type": "string"},
		{"text": "AbnormalCode", "type": "number"},
		{"text": "AbnormalPosition", "type": "string"},
		{"text": "WIP_NO", "type": "string"},
		{"text": "UNIT_NO", "type": "string"},
		{"text": "PLAN_QTY", "type": "number"},
		{"text": "MODEL_NO", "type": "string"},
		{"text": "LINE_DESC", "type": "string"},
		{"text": "ITEM_NO", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	// fmt.Println(grafanaData)
	return grafanaData
}
