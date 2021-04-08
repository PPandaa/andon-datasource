package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bitly/go-simplejson"
	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	envName = "demo.env"
)

func getDBInfo() (string, string, string) {
	err := godotenv.Load(envName)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	database := os.Getenv("MONGODB_DATABASE")
	return username, password, database
}

func createMongoSession() *mgo.Session {
	// load .env file
	err := godotenv.Load(envName)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	addr := os.Getenv("MONGODB_URL")
	// Open the connection
	session, _ := mgo.Dial(addr)
	// return the connection
	return session
}

func closeMongoSession(Session *mgo.Session) {
	Session.Close()
}

//TestConnection ...
func TestConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("/")
	w.Header().Set("Server", "Grafana Datasource Server")
	w.WriteHeader(200)
	msg := "This is DataSource for iFactory/Andon"
	w.Write([]byte(msg))
}

//Search ...
func Search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/search")
	requestBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Body: ", string(requestBody))
	metrics := []string{"AbnormalEventLatestTable", "DO-Singlestat", "AO-Singlestat", SfcWorkOrderDetail, SfcWorkOrderList, SfcStatsStation, SfcCounts, SfcSwitchingPanelWorkorderIds}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(metrics)
}

//Query ...
func Query(w http.ResponseWriter, r *http.Request) {
	var grafnaResponseArray []map[string]interface{}
	fmt.Println("/query")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	// fmt.Println(requestBody.Get("targets").MustArray())
	for indexOfTargets := 0; indexOfTargets < len(requestBody.Get("targets").MustArray()); indexOfTargets++ {
		target := requestBody.Get("targets").GetIndex(indexOfTargets).Get("target").MustString()

		//yoga
		var station string
		var orderId string
		targetOrderId := requestBody.Get("targets").GetIndex(indexOfTargets).Get("orderId").MustString()
		if targetOrderId != "" {
			//get scopedVars OrderId
			scopedVarsJson, _ := requestBody.Get("scopedVars").MarshalJSON()
			scopedVarsJsonValue := gjson.GetBytes(scopedVarsJson, "orderId.value").String()
			fmt.Println("orderId value:", scopedVarsJsonValue)
			orderId = scopedVarsJsonValue
		}
		targetStation := requestBody.Get("station").GetIndex(indexOfTargets).Get("station").MustString()
		if targetStation != "" {
			//get scopedVars OrderId
			scopedVarsJson, _ := requestBody.Get("scopedVars").MarshalJSON()
			scopedVarsJsonValue := gjson.GetBytes(scopedVarsJson, "station.value").String()
			fmt.Println("station value:", scopedVarsJsonValue)
			station = scopedVarsJsonValue
		}

		dataType := requestBody.Get("targets").GetIndex(indexOfTargets).Get("type").MustString()
		if dataType == "table" {
			switch target {
			case "AbnormalEventLatest":
				grafnaResponseArray = append(grafnaResponseArray, abnormalEventLatestTable())
			case "AO-Singlestat":
				grafnaResponseArray = append(grafnaResponseArray, abnormalOverviewSinglestat())
			case "DO-Singlestat":
				grafnaResponseArray = append(grafnaResponseArray, deviceOverviewSinglestat())
			case SfcWorkOrderDetail:
				grafnaResponseArray = append(grafnaResponseArray, GetWorkOrderDetail(orderId, station)) //工單狀態
			case SfcWorkOrderList:
				grafnaResponseArray = append(grafnaResponseArray, GetWorkOrderList(orderId, station)) //報工紀錄清單
			case SfcStatsStation:
				grafnaResponseArray = append(grafnaResponseArray, GetTables()) //工單工站狀態
			case SfcCounts:
				grafnaResponseArray = append(grafnaResponseArray, GetCounts())
			}
		}
	}
	// jsonStr, _ := json.Marshal(grafnaResponseArray)
	// fmt.Println(string(jsonStr))
	// fmt.Println(grafnaResponseArray)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(grafnaResponseArray)
}

func GetWorkorderIds(w http.ResponseWriter, r *http.Request) {
	res := getWorkorderIds()
	json.NewEncoder(w).Encode(res)
}

func abnormalEventLatestTable() map[string]interface{} {
	session := createMongoSession()
	username, password, database := getDBInfo()
	db := session.DB(database)
	db.Login(username, password)
	groupID := "R3JvdXA.X-0tJMYGAgAG-fkZ"
	collection := "iii.dae.EventLatest"
	var results []map[string]interface{}
	var rows []interface{}
	db.C(collection).Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"EventCode": 1}}}).All(&results)
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
		row = append(row, result["AbnormalStartTime"])
		row = append(row, result["ProcessingStatusCode"])
		row = append(row, result["ProcessingProgress"])
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
		{"text": "AbnormalStartTime", "type": "time"},
		{"text": "ProcessingStatusCode", "type": "number"},
		{"text": "ProcessingProgress", "type": "string"},
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
	closeMongoSession(session)
	return grafanaData
}

func abnormalOverviewSinglestat() map[string]interface{} {
	session := createMongoSession()
	username, password, database := getDBInfo()
	db := session.DB(database)
	db.Login(username, password)
	groupID := "R3JvdXA.X-0tJMYGAgAG-fkZ"
	collection := "iii.dae.Statistics"
	var result map[string]interface{}
	var rows []interface{}
	db.C(collection).Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"Dashboard": "AO"}}}).One(&result)
	// fmt.Println("Result: ", result)

	var row []interface{}
	row = append(row, result["AbnormalMachine"])
	row = append(row, result["TodayRepair"])
	row = append(row, result["Overdue"])
	rows = append(rows, row)
	// fmt.Println(rows)
	columns := []map[string]string{
		{"text": "AbnormalMachine", "type": "number"},
		{"text": "TodayRepair", "type": "number"},
		{"text": "Overdue", "type": "number"},
	}
	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    rows,
		"type":    "table",
	}
	closeMongoSession(session)
	return grafanaData
}

func deviceOverviewSinglestat() map[string]interface{} {
	session := createMongoSession()
	username, password, database := getDBInfo()
	db := session.DB(database)
	db.Login(username, password)
	groupID := "R3JvdXA.X-0tJMYGAgAG-fkZ"
	collection := "iii.dae.Statistics"
	var result map[string]interface{}
	var rows []interface{}
	db.C(collection).Pipe([]bson.M{{"$match": bson.M{"GroupID": groupID}}, {"$match": bson.M{"Dashboard": "DO"}}}).One(&result)
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
	closeMongoSession(session)
	return grafanaData
}
