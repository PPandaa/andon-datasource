package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/imroc/req"

	// . "github.com/logrusorgru/aurora"
	"github.com/tidwall/gjson"
)

func GetCounts(orderId, station string) map[string]interface{} {
	trigger := func(i interface{}) ([]byte, error) {
		url := apiUrl + "/grafana/counts"
		if station != "" {
			url = url + "?station=" + station
		}
		//convert object to json
		param := req.BodyJSON(&i)
		//res就是打api成功拿到的response, 如果打失敗則拿到err
		res, err := DoAPI("GET", url, param)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	res, _ := trigger(nil)

	var grafanaData map[string]interface{}
	err := json.Unmarshal(res, &grafanaData)
	if err != nil {
		glog.Error(err)
	}

	return grafanaData
}

// 工單工站狀態 SfcStatsStation (目前未使用)
func GetTables(orderId, station string) map[string]interface{} {
	PrintParameter(orderId, station)
	trigger := func(i interface{}) ([]byte, error) {
		url := apiUrl + "/grafana/tables?groupBy=station"
		if orderId != "" {
			url = url + "&workorderId=" + orderId
		}
		//convert object to json
		param := req.BodyJSON(&i)
		//res就是打api成功拿到的response, 如果打失敗則拿到err
		res, err := DoAPI("GET", url, param)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	res, _ := trigger(nil)

	var Results []map[string]interface{}
	var Rows [][]interface{}
	Results = JsonAryToMap(res)

	for _, result := range Results {
		// var result map[string]interface{}
		var row []interface{}
		row = append(row, result["WorkOrderId"])
		row = append(row, result["StationName"])

		row = append(row, result["GoodQty"])
		row = append(row, result["NonGoodQty"])
		row = append(row, result["CompletedQty"])
		row = append(row, result["ToBeCompletedQty"])
		row = append(row, result["Quantity"])

		row = append(row, result["RealCompletedRate"])
		row = append(row, result["EstiCompletedRate"])

		row = append(row, result["Status"])

		Rows = append(Rows, row)
	}

	if len(Rows) == 0 {
		row := []interface{}{}
		Rows = append(Rows, row)
	}

	columns := []map[string]string{
		{"text": "WorkOrderId", "type": "string"},
		{"text": "StationName", "type": "string"},

		{"text": "GoodQty", "type": "string"},
		{"text": "NonGoodQty", "type": "string"},
		{"text": "CompletedQty", "type": "string"},
		{"text": "ToBeCompletedQty", "type": "string"},
		{"text": "Quantity", "type": "string"},

		{"text": "RealCompletedRate", "type": "string"},
		{"text": "EstiCompletedRate", "type": "string"},

		{"text": "Status", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    Rows,
		"type":    "table",
	}

	return grafanaData
}

//報工單紀錄
func GetWorkOrderList(orderId, station string) map[string]interface{} {
	fmt.Println("-------------GetWorkOrderList-------------")
	PrintParameter(orderId, station)
	trigger := func(i interface{}) ([]byte, error) {
		url := apiUrl + "/workorders"
		if orderId != "" {
			url = url + "?workorderId=" + orderId
		}
		//convert object to json
		param := req.BodyJSON(&i)
		//res就是打api成功拿到的response, 如果打失敗則拿到err
		res, err := DoAPI("GET", url, param)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	res, _ := trigger(nil)

	var Results []map[string]interface{}
	Rows := [][]interface{}{}

	//method2
	gj := gjson.GetBytes(res, "#.WorkOrderList")
	for _, wols := range gj.Array() {
		// for _, wol := range wols.Array() {
		// }
		Results = JsonAryToMap([]byte(wols.Raw))
		for _, result := range Results {
			var row []interface{}
			row = append(row, result["WorkOrderId"])
			row = append(row, result["SubWorkOrderId"])
			row = append(row, result["Reporter"])
			row = append(row, result["StationName"])
			row = append(row, result["CompletedQty"])
			row = append(row, result["NonGoodQty"])
			row = append(row, result["CreateAt"])
			Rows = append(Rows, row)
		}
	}

	if len(Rows) == 0 {
		row := []interface{}{}
		Rows = append(Rows, row)
	}

	columns := []map[string]string{
		{"text": "WorkOrderId", "type": "string"},
		{"text": "SubWorkOrderId", "type": "string"},
		{"text": "Reporter", "type": "string"},
		{"text": "StationName", "type": "string"},
		{"text": "CompletedQty", "type": "string"},
		{"text": "NonGoodQty", "type": "string"},
		{"text": "CreateAt", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    Rows,
		"type":    "table",
	}
	b, _ := json.Marshal(grafanaData)
	fmt.Println(string(b))
	fmt.Println(Rows)

	return grafanaData
}

//工單狀態
func GetWorkOrderDetail(orderId, station string) map[string]interface{} {
	PrintParameter(orderId, station)
	trigger := func(i interface{}) ([]byte, error) {
		url := apiUrl + "/workorders"
		if station != "" {
			url = url + "&station=" + station
		}
		// url := apiUrl + "/grafana/table/workorderDetail" 	//做到一半

		//convert object to json
		param := req.BodyJSON(&i)
		res, err := DoAPI("GET", url, param)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	res, _ := trigger(nil)

	//做到一半
	// fmt.Println(string(res))
	// var m map[string]interface{}
	// err := json.Unmarshal(res, &m)
	// if err != nil {
	// 	log.Error(err)
	// }

	// return m

	var Results []map[string]interface{}
	var Rows [][]interface{}
	Results = JsonAryToMap(res)

	for _, result := range Results {
		// var result map[string]interface{}
		var row []interface{}
		row = append(row, result["WorkOrderId"])
		// row = append(row, result["Station"])

		row = append(row, result["Product"].(map[string]interface{})["ProductId"])
		row = append(row, result["Product"].(map[string]interface{})["ProductName"])
		row = append(row, result["Product"].(map[string]interface{})["StationType"])

		row = append(row, result["CompletedQty"])
		row = append(row, result["NonGoodQty"])
		row = append(row, result["GoodQty"])
		row = append(row, result["GoodQtyRate"])
		row = append(row, result["GoodProductQty"])
		row = append(row, result["Quantity"])

		row = append(row, result["Status"])

		row = append(row, result["PlanStartDate"])
		row = append(row, result["DeliverAt"])
		row = append(row, result["CreateAt"])
		Rows = append(Rows, row)
	}

	if len(Rows) == 0 {
		row := []interface{}{}
		Rows = append(Rows, row)
	}

	columns := []map[string]string{
		{"text": "WorkOrderId", "type": "string"},
		// {"text": "Station", "type": "string"},

		{"text": "ProductId", "type": "string"},
		{"text": "ProductName", "type": "string"},
		{"text": "StationType", "type": "string"},

		{"text": "CompletedQty", "type": "string"},
		{"text": "NonGoodQty", "type": "string"},
		{"text": "GoodQty", "type": "string"},
		{"text": "GoodQtyRate", "type": "string"},
		{"text": "GoodProductQty", "type": "string"},
		{"text": "Quantity", "type": "string"},

		{"text": "Status", "type": "string"},

		{"text": "PlanStartDate", "type": "string"},
		{"text": "DeliverAt", "type": "string"},
		{"text": "CreateAt", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    Rows,
		"type":    "table",
	}

	return grafanaData
}
