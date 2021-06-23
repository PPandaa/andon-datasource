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

//Search ...
func Search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/search")
	requestBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Body: ", string(requestBody))
	metrics := getMetrics()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(metrics)
}

//Query ...
func Query(w http.ResponseWriter, r *http.Request) {
	// var grafnaResponseArray []map[string]interface{}
	var grafnaResponseArray []interface{}

	// fmt.Println("/query")
	requestBody, _ := simplejson.NewFromReader(r.Body)
	// fmt.Println("Body: ", requestBody)
	// fmt.Println(requestBody.Get("targets").MustArray())
	for indexOfTargets := 0; indexOfTargets < len(requestBody.Get("targets").MustArray()); indexOfTargets++ {
		target := requestBody.Get("targets").GetIndex(indexOfTargets).Get("metrics").MustString()

		var station string
		var orderId string
		var timeFrom string
		var group string

		//orderId
		targetOrderId := requestBody.Get("targets").GetIndex(indexOfTargets).Get("orderId").MustString()
		if targetOrderId != "" {
			//get scopedVars OrderId
			scopedVarsJson, _ := requestBody.Get("scopedVars").MarshalJSON()
			scopedVarsJsonValue := gjson.GetBytes(scopedVarsJson, "orderId.value").String()
			// fmt.Println("orderId value:", scopedVarsJsonValue)
			orderId = scopedVarsJsonValue
		}
		//station
		targetStation := requestBody.Get("targets").GetIndex(indexOfTargets).Get("stationId").MustString()
		if targetStation != "" {
			//get scopedVars OrderId
			scopedVarsJson, _ := requestBody.Get("scopedVars").MarshalJSON()
			scopedVarsJsonValue := gjson.GetBytes(scopedVarsJson, "machineId.value").String() //#注意stationId對應的是machineId
			// fmt.Println("station value:", scopedVarsJsonValue)
			station = scopedVarsJsonValue
		}
		// time
		func() {
			// from 不會放在target裡面
			scopedVarsJson, _ := requestBody.Get("scopedVars").MarshalJSON()
			scopedVarsJsonValue := gjson.GetBytes(scopedVarsJson, "__from.value").String()
			// fmt.Println("time from value:", scopedVarsJsonValue)
			timeFrom = scopedVarsJsonValue
		}()
		//station
		targetGroup := requestBody.Get("targets").GetIndex(indexOfTargets).Get("group").MustString()
		if targetGroup != "" {
			//get scopedVars OrderId
			scopedVarsJson, _ := requestBody.Get("scopedVars").MarshalJSON()
			scopedVarsJsonValue := gjson.GetBytes(scopedVarsJson, "machineId.value").String() //#注意stationId對應的是machineId
			// fmt.Println("station value:", scopedVarsJsonValue)
			group = scopedVarsJsonValue
		}

		PrintParameterCyan("station:", station, "orderId:", orderId, "timeFrom:", timeFrom, "group:", group)

		dataType := requestBody.Get("targets").GetIndex(indexOfTargets).Get("type").MustString()

		if dataType == "table" {
			fmt.Println("target:", target)
			grafnaResponseArray = append(grafnaResponseArray, doFuncByMetric(target, orderId, station, timeFrom, group))
		}
	}
	// jsonStr, _ := json.Marshal(grafnaResponseArray)
	// fmt.Println(string(jsonStr))
	// fmt.Println(grafnaResponseArray)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(grafnaResponseArray)
}
