package table

import (
	"DataSource/config"
	"fmt"
	"sort"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func MachinesSinglestat(groupID string) map[string]interface{} {
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

func AbnormalReasonRank(groupID string, fromTime time.Time, toTime time.Time) map[string]interface{} {
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

func LastMonthAbReasonRank(groupID string) map[string]interface{} {
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

func MeanTimeComputeSinglestat(computeType string, groupID string, totalComputeValue float64, fromTime time.Time, toTime time.Time) map[string]interface{} {
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

func MeanTimeComputeV2(computeType string, groupID string, fromTime time.Time, toTime time.Time) (map[string]interface{}, float64) {
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

func MeanTimeComputeV3(computeType string, groupID string, fromTime time.Time, toTime time.Time) map[string]interface{} {
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

func Panel1Singlestat(groupID string) map[string]interface{} {
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

func Panel1Table(groupID string) map[string]interface{} {
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

func Panel3(groupID string) map[string]interface{} {
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

func Panel5Singlestat(groupID string) map[string]interface{} {
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

func Panel5Table(groupID string) map[string]interface{} {
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

func Panel6Singlestat(groupID string) map[string]interface{} {
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

func Panel6Table(groupID string) map[string]interface{} {
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

func Panel8(groupID string) map[string]interface{} {
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

func Panel9(groupID string) map[string]interface{} {
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

func Panel10(groupID string) map[string]interface{} {
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
