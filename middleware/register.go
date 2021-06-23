package middleware

import (
	"log"
	"os"
)

//docker login -u any99147 -p 54P@ssw0rd && ./build_dev.sh

var Url string

func Start() {
	Url = os.Getenv("API_URL")
	log.Println("url=", Url)
}

var Metrics map[string]func(orderId string, station string, timeFrom string, group string) interface{}

func init() {
	Metrics = make(map[string]func(orderId string, station string, timeFrom string, group string) interface{})
	setMetrics()
}

/*
	// grafana simple json
	//wo
	apiv1.GET("/grafana/table/Wo", v1.GetTableWorkOrder)
	apiv1.GET("/grafana/table/Wo/Count", v1.GetTableWorkOrderCount)
	apiv1.GET("/grafana/table/CompletedWo", v1.GetCompletedWO)
	apiv1.GET("/grafana/table/CompletedWo/Count", v1.GetCompletedWOCount)
	apiv1.GET("/grafana/table/ExecutionWo", v1.GetExecutionWO)
	apiv1.GET("/grafana/table/ExecutionWo/Count", v1.GetExecutionWOCount)
	apiv1.GET("/grafana/table/IdleWo", v1.GetTop10IdleWO)
	apiv1.GET("/grafana/table/DefecRateWOProcess", v1.GetDefecRateWOProcess)
	//mo
	apiv1.GET("/grafana/table/CompletedMo", v1.GetCompletedMo)
	apiv1.GET("/grafana/table/CompletedMo/Count", v1.GetCompletedMoCount)
	//wolist
	apiv1.GET("/grafana/table/OperationSpot", v1.GetOperationSpot)
*/

//todo: 抽出orderId string, station string, timeFrom string為struct

func setMetrics() {
	Metrics["Wo"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return Wo(orderId, station, timeFrom, group)
	}
	Metrics["WoCount"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return WoCount(orderId, station, timeFrom, group)
	}
	Metrics["CompletedWo"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return CompletedWo(orderId, station, timeFrom, group)
	}
	Metrics["CompletedWoCount"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return CompletedWoCount(orderId, station, timeFrom, group)
	}
	Metrics["ExecutionWo"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return ExecutionWo(orderId, station, timeFrom, group)
	}
	Metrics["ExecutionWoCount"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return ExecutionWoCount(orderId, station, timeFrom, group)
	}
	Metrics["IdleWo"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return IdleWo(orderId, station, timeFrom, group)
	}
	Metrics["DefecRateWOProcess"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return DefecRateWOProcess(orderId, station, timeFrom, group)
	}
	Metrics["CompletedMo"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return CompletedMo(orderId, station, timeFrom, group)
	}
	Metrics["CompletedMoCount"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return CompletedMoCount(orderId, station, timeFrom, group)
	}
	Metrics["OperationSpot"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return OperationSpot(orderId, station, timeFrom, group)
	}
	Metrics["Wolist"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return Wolist(orderId, station, timeFrom, group)
	}
	Metrics["OpenWoCount"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return OpenWoCount(orderId, station, timeFrom, group)
	}
	Metrics["Counts"] = func(orderId string, station string, timeFrom string, group string) interface{} {
		return GetCounts(orderId, station, timeFrom, group)
	}
}

func getMetrics() (ss []string) {
	for k, _ := range Metrics {
		ss = append(ss, k)
	}
	return
}

func doFuncByMetric(target, orderId, station string, timeFrom string, group string) interface{} {
	for k, funcs := range Metrics {
		if k == target {
			return funcs(orderId, station, timeFrom, group)
		}
	}
	return nil
}
