package middleware

import (
	"DataSource/config"
	"DataSource/pkg/table"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"gopkg.in/mgo.v2/bson"
)

func dbHealthCheck() {
	err := config.Session.Ping()
	if err != nil {
		fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
		fmt.Println("MongoDB", err, "->", "URL:", config.MongodbURL, " Database:", config.MongodbDatabase)
		config.Session.Refresh()
		fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
		fmt.Println("MongoDB Reconnect ->", " URL:", config.MongodbURL, " Database:", config.MongodbDatabase)
	}
}

func TestConnection(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/")
	w.Header().Set("Server", "Grafana Datasource Server")
	w.WriteHeader(200)
	msg := "This is DataSource for iFactory/Andon"
	w.Write([]byte(msg))
}

func GetGroup(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/group")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	groupTopoCollection := config.DB.C(config.GroupTopo)
	var groupTopoResults []map[string]interface{}
	groupNameArray := []map[string]string{}

	parentID := requestBody.Get("parentID").MustString()
	fmt.Println("ParentID:", parentID)

	if parentID == "ALL" {
		groupTopoCollection.Find(bson.M{}).All(&groupTopoResults)
	} else {
		groupTopoCollection.Find(bson.M{"ParentID": parentID}).All(&groupTopoResults)
	}

	for _, groupTopoResult := range groupTopoResults {
		temp := map[string]string{"text": groupTopoResult["GroupName"].(string), "value": groupTopoResult["GroupID"].(string)}
		groupNameArray = append(groupNameArray, temp)
	}
	// fmt.Println(groupNameArray)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(groupNameArray)
}

func GetMachine(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/machine")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	machineRawDataCollection := config.DB.C(config.MachineRawData)
	var machineRawDataResults []map[string]interface{}
	machineIDArray := []map[string]string{}

	groupID := requestBody.Get("groupID").MustString()
	fmt.Println("GroupID:", groupID)

	if groupID == "ALL" {
		machineRawDataCollection.Find(bson.M{}).All(&machineRawDataResults)
	} else {
		machineRawDataCollection.Find(bson.M{"GroupID": groupID}).All(&machineRawDataResults)
	}

	for _, machineRawDataResult := range machineRawDataResults {
		temp := map[string]string{"text": machineRawDataResult["MachineName"].(string), "value": machineRawDataResult["MachineID"].(string)}
		machineIDArray = append(machineIDArray, temp)
	}
	// fmt.Println(machineIDArray)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(machineIDArray)
}

func GetAssignee(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/assignee")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	eventHistCollection := config.DB.C(config.EventHist)
	var assigneeIDByGroupIDResults []string
	var assigneeIDByMachineIDResults []string
	assigneeIDArray := []map[string]string{}

	filterID := requestBody.Get("filterID").MustString()
	fmt.Println("FilterID:", filterID)

	eventHistCollection.Find(bson.M{"GroupID": filterID}).Distinct("PrincipalID", &assigneeIDByGroupIDResults)
	eventHistCollection.Find(bson.M{"MachineID": filterID}).Distinct("PrincipalID", &assigneeIDByMachineIDResults)

	for _, eventHistByGroupIDResult := range assigneeIDByGroupIDResults {
		var oneResult map[string]interface{}
		eventHistCollection.Find(bson.M{"PrincipalID": eventHistByGroupIDResult}).One(&oneResult)
		temp := map[string]string{"text": oneResult["PrincipalName"].(string), "value": oneResult["PrincipalID"].(string)}
		assigneeIDArray = append(assigneeIDArray, temp)
	}

	for _, eventHistByMachineIDResults := range assigneeIDByMachineIDResults {
		var oneResult map[string]interface{}
		eventHistCollection.Find(bson.M{"PrincipalID": eventHistByMachineIDResults}).One(&oneResult)
		temp := map[string]string{"text": oneResult["PrincipalName"].(string), "value": oneResult["PrincipalID"].(string)}
		assigneeIDArray = append(assigneeIDArray, temp)
	}
	// fmt.Println(assigneeIDArray)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(assigneeIDArray)
}

func Search(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/search")
	requestBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Body: ", string(requestBody))

	metrics := []string{"EventLatest", "EventHist", "EventList", "LastMonthAbReasonRank", "MTTD", "MTTR", "MTBF", "FlowCharting-Machines", "FlowCharting-TPC", "Singlestat-Machines", "Singlestat-Event", "Singlestat-MeanTimeCompute", "Panel1Singlestat", "Panel1Table", "Panel2Singlestat", "Panel2Table", "Panel3AndPanel4", "Panel5Singlestat", "Panel5Table", "Panel6Singlestat", "Panel6Table", "Panel7", "Panel8", "Panel9", "Panel10"}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(metrics)
}

func Query(w http.ResponseWriter, r *http.Request) {
	dbHealthCheck()
	grafnaResponseArray := []map[string]interface{}{}
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("/query")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	fmt.Println("Body: ", requestBody)
	// fmt.Println("Targets: ", requestBody.Get("targets").MustArray())

	for indexOfTargets := 0; indexOfTargets < len(requestBody.Get("targets").MustArray()); indexOfTargets++ {
		refID := requestBody.Get("targets").GetIndex(indexOfTargets).Get("refId").MustString()
		dataType := requestBody.Get("targets").GetIndex(indexOfTargets).Get("type").MustString()
		groupID := requestBody.Get("targets").GetIndex(indexOfTargets).Get("groupID").MustString()
		if strings.HasPrefix(groupID, "$") {
			temp := strings.Split(groupID, "$")
			varName := temp[1]
			groupID = requestBody.Get("scopedVars").Get(varName).Get("value").MustString()
		}
		machineID := requestBody.Get("targets").GetIndex(indexOfTargets).Get("machineID").MustString()
		if strings.HasPrefix(machineID, "$") {
			temp := strings.Split(machineID, "$")
			varName := temp[1]
			machineID = requestBody.Get("scopedVars").Get(varName).Get("value").MustString()
		}
		metrics := requestBody.Get("targets").GetIndex(indexOfTargets).Get("metrics").MustString()
		fromTimeString := requestBody.Get("scopedVars").Get("__from").Get("value").MustString()
		fromUnixTime, _ := strconv.ParseInt(fromTimeString, 10, 64)
		// fmt.Println("FromUnixTime:", fromUnixTime)
		temp := fromUnixTime / 1000
		fromTime := time.Unix(temp, 0)
		fromTime = fromTime.In(config.TaipeiTimeZone)
		toTimeSting := requestBody.Get("scopedVars").Get("__to").Get("value").MustString()
		toUnixTime, _ := strconv.ParseInt(toTimeSting, 10, 64)
		// fmt.Println("ToUnixTime:", toUnixTime)
		temp = toUnixTime / 1000
		toTime := time.Unix(temp, 0)
		toTime = toTime.In(config.TaipeiTimeZone)
		// fmt.Println("  RefID:", refID, "DataType:", dataType, " GroupID:", groupID, " MachineID:", machineID, " Metrics:", metrics, " From:", fromUnixTime, fromTime.Format(time.RFC3339), " To:", toUnixTime, toTime.Format(time.RFC3339))
		if dataType == "table" {
			switch metrics {
			case "EventLatest":
				grafnaResponseArray = append(grafnaResponseArray, table.EventLatest(groupID, machineID))
			case "EventHist":
				grafnaResponseArray = append(grafnaResponseArray, table.EventHist(groupID, machineID, fromTime, toTime))
			case "EventList":
				grafnaResponseArray = append(grafnaResponseArray, table.EventList(groupID, fromTime, toTime))
			case "Singlestat-Event":
				grafnaResponseArray = append(grafnaResponseArray, table.EventSinglestat(groupID, machineID))
			case "Singlestat-Machines":
				grafnaResponseArray = append(grafnaResponseArray, table.MachinesSinglestat(groupID))
			case "LastMonthAbReasonRank":
				grafnaResponseArray = append(grafnaResponseArray, table.AbnormalReasonRank(groupID, fromTime, toTime))
			case "FlowCharting-Machines":
				grafnaResponseArray = append(grafnaResponseArray, table.MachinesFlowCharting(groupID, refID))
			case "FlowCharting-TPC":
				grafnaResponseArray = append(grafnaResponseArray, table.TpcFlowCharting(groupID, refID))
			case "MTTD", "MTTR", "MTBF":
				grafnaResponseArray = append(grafnaResponseArray, table.MeanTimeComputeV3(metrics, groupID, fromTime, toTime))
				// v2
				// grafnaResponseArrayElement, totalComputeValue := v2MeanTimeCompute(metrics, groupID, fromTime, toTime)
				// grafnaResponseArray = append(grafnaResponseArray, grafnaResponseArrayElement)
				// grafnaResponseArrayElement = meanTimeComputeSinglestat(metrics, groupID, totalComputeValue, fromTime, toTime)
				// grafnaResponseArray = append(grafnaResponseArray, grafnaResponseArrayElement)
			case "Panel1Singlestat":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel1Singlestat(groupID))
			case "Panel1Table":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel1Table(groupID))
			case "Panel2Singlestat":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel2Singlestat(groupID))
			case "Panel2Table":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel2Table(groupID))
			case "Panel3AndPanel4":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel3AndPanel4(groupID))
			case "Panel5Singlestat":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel5Singlestat(groupID))
			case "Panel5Table":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel5Table(groupID))
			case "Panel6Singlestat":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel6Singlestat(groupID))
			case "Panel6Table":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel6Table(groupID))
			case "Panel8":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel8(groupID))
			case "Panel9":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel9(groupID))
			case "Panel10":
				grafnaResponseArray = append(grafnaResponseArray, table.Panel10(groupID))
			}

		} else {
			switch metrics {
			// case "MTTD", "MTTR", "MTBF":
			// 	grafnaResponseArray = append(grafnaResponseArray, timeseries.MeanTimeComputeV1(metrics, groupID, fromTime, toTime))
			}
		}
	}

	// jsonStr, _ := json.Marshal(grafnaResponseArray)
	// fmt.Println(string(jsonStr))
	// fmt.Println(grafnaResponseArray)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(grafnaResponseArray)
}

//Test ...
// func Test(w http.ResponseWriter, r *http.Request) {
