package timeseries

import (
	"DataSource/config"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func MeanTimeComputeV1(computeType string, groupID string) map[string]interface{} {
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

func FactoryMap(groupID string) []map[string]interface{} {
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
