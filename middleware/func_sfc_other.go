package middleware

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/imroc/req"
)

func GetCounts(orderId, station, timeFrom, group string) map[string]interface{} {
	trigger := func(i interface{}) ([]byte, error) {
		url := Url + "/grafana/counts"
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

//"/grafana/table/workorders/info"
func trigger(api string, bodyParameter interface{}) (grafanaData map[string]interface{}, err error) {
	url := Url + api
	//convert object to json
	param := req.BodyJSON(&bodyParameter)
	//res就是打api成功拿到的response, 如果打失敗則拿到err
	res, err := DoAPI("GET", url, param)
	if err != nil {
		return nil, err
	}

	PrintParameter("trigger response:", string(res))

	err = json.Unmarshal(res, &grafanaData)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return grafanaData, nil
}
