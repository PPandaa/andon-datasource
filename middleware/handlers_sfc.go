package middleware

import "log"

var (
	sfcWorkOrders = "sfc_workorders"
)

func findWorkOrders() map[string]interface{} {
	session := createMongoSession()
	username, password, database := getDBInfo()
	db := session.DB(database)
	db.Login(username, password)

	collection := "iii.sfc.workoder"

	var Results []map[string]interface{}
	var Rows [][]interface{}

	err := db.C(collection).Find(nil).All(&Results)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("Result: ", result)

	for _, result := range Results {
		// var result map[string]interface{}
		var row []interface{}
		row = append(row, result["WorkOrderId"])
		row = append(row, result["Station"])
		row = append(row, result["Machine"])
		row = append(row, result["Status"])
		row = append(row, result["Good"])
		row = append(row, result["NonGood"])
		row = append(row, result["Quantity"])
		row = append(row, result["Reporter"])
		row = append(row, result["PlanStartDate"])
		Rows = append(Rows, row)
	}

	columns := []map[string]string{
		{"text": "WorkOrderId", "type": "string"},
		{"text": "Station", "type": "string"},
		{"text": "Machine", "type": "string"},
		{"text": "Status", "type": "string"},
		{"text": "Good", "type": "string"},
		{"text": "NonGood", "type": "string"},
		{"text": "Quantity", "type": "string"},
		{"text": "Reporter", "type": "string"},
		{"text": "PlanStartDate", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    Rows,
		"type":    "table",
	}
	closeMongoSession(session)
	return grafanaData
}
