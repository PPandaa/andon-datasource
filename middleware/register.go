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

var Metrics map[string]func(orderId string, station string) interface{}

func init() {
	Metrics = make(map[string]func(orderId string, station string) interface{})
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

func setMetrics() {
	Metrics["Wo"] = func(a string, b string) interface{} { return Wo(a, b) }
	Metrics["WoCount"] = func(a string, b string) interface{} { return WoCount(a, b) }
	Metrics["CompletedWo"] = func(a string, b string) interface{} { return CompletedWo(a, b) }
	Metrics["CompletedWoCount"] = func(a string, b string) interface{} { return CompletedWoCount(a, b) }
	Metrics["ExecutionWo"] = func(a string, b string) interface{} { return ExecutionWo(a, b) }
	Metrics["ExecutionWoCount"] = func(a string, b string) interface{} { return ExecutionWoCount(a, b) }
	Metrics["IdleWo"] = func(a string, b string) interface{} { return IdleWo(a, b) }
	Metrics["DefecRateWOProcess"] = func(a string, b string) interface{} { return DefecRateWOProcess(a, b) }

	Metrics["CompletedMo"] = func(a string, b string) interface{} { return CompletedMo(a, b) }
	Metrics["CompletedMoCount"] = func(a string, b string) interface{} { return CompletedMoCount(a, b) }

	Metrics["OperationSpot"] = func(a string, b string) interface{} { return OperationSpot(a, b) }

	Metrics["Counts"] = func(a string, b string) interface{} { return GetCounts(a, b) }
}

func getMetrics() (ss []string) {
	for k, _ := range Metrics {
		ss = append(ss, k)
	}
	return
}

func doFuncByMetric(target, orderId, station string) interface{} {
	for k, funcs := range Metrics {
		if k == target {
			return funcs(orderId, station)
		}
	}
	return nil
}
