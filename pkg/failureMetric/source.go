package failureMetric

import (
	"DataSource/config"
	"sort"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func V1_MeanTimeCompute(computeType string, groupID string, fromTime time.Time, toTime time.Time) map[string]interface{} {
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

func V1_AbnormalReasonRank(groupID string, fromTime time.Time, toTime time.Time) map[string]interface{} {
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
