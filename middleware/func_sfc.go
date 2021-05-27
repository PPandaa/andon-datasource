package middleware

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

func Wo(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/Wo"
	m, _ := trigger(api, nil)
	return m
}

func WoCount(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/Wo/Count"
	m, _ := trigger(api, nil)
	return m
}

func CompletedWo(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/CompletedWo"
	m, _ := trigger(api, nil)
	return m
}

func CompletedWoCount(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/CompletedWo/Count"
	m, _ := trigger(api, nil)
	return m
}

func ExecutionWoCount(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/ExecutionWo/Count"
	m, _ := trigger(api, nil)
	return m
}
func ExecutionWo(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/ExecutionWo"
	m, _ := trigger(api, nil)
	return m
}

func IdleWo(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/IdleWo"
	m, _ := trigger(api, nil)
	return m
}

func DefecRateWOProcess(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/DefecRateWOProcess"
	m, _ := trigger(api, nil)
	return m
}

func CompletedMo(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/CompletedMo"
	m, _ := trigger(api, nil)
	return m
}
func CompletedMoCount(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/CompletedMo/Count"
	m, _ := trigger(api, nil)
	return m
}
func OperationSpot(orderId, station string) map[string]interface{} {
	// PrintParameter(orderId, station)
	api := "/grafana/table/OperationSpot"
	m, _ := trigger(api, nil)
	return m
}
