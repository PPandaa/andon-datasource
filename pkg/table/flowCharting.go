package table

import (
	"DataSource/config"

	"gopkg.in/mgo.v2/bson"
)

func MachinesFlowCharting(groupID string, refID string) map[string]interface{} {
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
			if machineRawDataResult["StatusLay1Value"] == nil {
				row = append(row, 4000)
			} else {
				row = append(row, machineRawDataResult["StatusLay1Value"])
			}
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

func TpcFlowCharting(groupID string, refID string) map[string]interface{} {
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
