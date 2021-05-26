package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/tidwall/gjson"
)

//TestConnection ...
func TestConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("/")
	w.Header().Set("Server", "Grafana Datasource Server")
	w.WriteHeader(200)
	msg := "This is DataSource for iFactory/Sfc"
	w.Write([]byte(msg))
}

var (
	SfcWorkOrderDetail            = "SfcWorkOrderDetail"
	SfcWorkOrderList              = "SfcWorkOrderList"
	SfcStatsStation               = "SfcStatsStation"
	SfcCounts                     = "SfcCounts"
	SfcSwitchingPanelWorkorderIds = "SfcSwitchingPanelWorkorderIds"
)

//Search ...
func Search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/search")
	requestBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Body: ", string(requestBody))
	metrics := []string{SfcWorkOrderDetail, SfcWorkOrderList, SfcStatsStation, SfcCounts, SfcSwitchingPanelWorkorderIds}
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
		target := requestBody.Get("targets").GetIndex(indexOfTargets).Get("metrics").MustString()

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
		targetStation := requestBody.Get("targets").GetIndex(indexOfTargets).Get("stationId").MustString()
		if targetStation != "" {
			//get scopedVars OrderId
			scopedVarsJson, _ := requestBody.Get("scopedVars").MarshalJSON()
			scopedVarsJsonValue := gjson.GetBytes(scopedVarsJson, "machineId.value").String() //#注意stationId對應的是machineId
			fmt.Println("station value:", scopedVarsJsonValue)
			station = scopedVarsJsonValue
		}

		TestParameter("station:", station)
		TestParameter("orderId:", orderId)

		dataType := requestBody.Get("targets").GetIndex(indexOfTargets).Get("type").MustString()
		if dataType == "table" {
			switch target {
			case SfcWorkOrderDetail:
				grafnaResponseArray = append(grafnaResponseArray, GetWorkOrderDetail(orderId, station)) //工單狀態
			case SfcWorkOrderList:
				grafnaResponseArray = append(grafnaResponseArray, GetWorkOrderList(orderId, station)) //報工紀錄清單
			case SfcStatsStation:
				grafnaResponseArray = append(grafnaResponseArray, GetTables(orderId, station)) //工單工站狀態
			case SfcCounts:
				grafnaResponseArray = append(grafnaResponseArray, GetCounts(orderId, station))
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
