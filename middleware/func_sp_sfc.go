package middleware

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/imroc/req"
	// . "github.com/logrusorgru/aurora"
)

//switchingPanel
func getWorkorderIds() []map[string]interface{} {
	trigger := func(i interface{}) ([]byte, error) {
		url := apiUrl + "/grafana/switchingPanel/workorders/id"
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

	var grafanaData []map[string]interface{}
	err := json.Unmarshal(res, &grafanaData)
	if err != nil {
		glog.Error(err)
	}

	return grafanaData
}
