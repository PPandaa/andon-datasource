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

func Wo(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/Wo"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func WoCount(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/Wo/Count"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func CompletedWo(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/CompletedWo"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func CompletedWoCount(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/CompletedWo/Count"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func ExecutionWoCount(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/ExecutionWo/Count"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}
func ExecutionWo(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/ExecutionWo"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func IdleWo(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/IdleWo"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func DefecRateWOProcess(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/DefecRateWOProcess"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func CompletedMo(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/CompletedMo"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}
func CompletedMoCount(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/CompletedMo/Count"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}
func OperationSpot(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/OperationSpot"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

//orderId= workorderId= moId
func Wolist(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/Wolist"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func OpenWoCount(orderId, station, timeFrom, group string) map[string]interface{} {
	api := "/grafana/table/OpenWo/Count"
	api = concateUrl(api, orderId, station, timeFrom, group)
	m, _ := trigger(api, nil)
	return m
}

func concateUrl(api, orderId, station, timeFrom, group string) string {
	api = api + "?"
	if orderId != "" {
		api = api + "workorderId=" + orderId + "&"
	}
	if station != "" {
		api = api + "station=" + station + "&"
	}
	if timeFrom != "" {
		api = api + "timeFrom=" + timeFrom + "&"
	}
	if group != "" {
		api = api + "group=" + group + "&"
	}
	PrintParameter("---------url:", api, "-----------")
	return api
}
