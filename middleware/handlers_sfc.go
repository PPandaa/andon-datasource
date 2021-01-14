package middleware

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang/glog"
	"github.com/imroc/req"
	"github.com/logrusorgru/aurora"
	. "github.com/logrusorgru/aurora"
)

//docker login -u any99147 -p 54P@ssw0rd && ./build_dev.sh

var (
	SfcWorkOrderInfo = "SfcWorkOrderInfo"
	SfcWorkOrderList = "SfcWorkOrderList"
)

//deprecated
func findWorkOrderList() map[string]interface{} {
	// session := createMongoSession()
	// username, password, database := getDBInfo()
	// db := session.DB(database)
	// db.Login(username, password)

	// collection := "iii.sfc.workorder_list"

	// var Results []map[string]interface{}
	// var Rows [][]interface{}

	// err := db.C(collection).Find(nil).All(&Results)
	// if err != nil {
	// 	log.Println(err)
	// }
	trigger := func(i interface{}) ([]byte, error) {
		url := "https://andon-daemon-compute-ifactoryandondev-eks005.sa.wise-paas.com/workorder_list"
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

	var Results []map[string]interface{}
	var Rows [][]interface{}
	Results = JsonAryToMap(res)

	for _, result := range Results {
		// var result map[string]interface{}
		var row []interface{}
		row = append(row, result["WorkOrderId"])
		row = append(row, result["Station"])
		row = append(row, result["Machine"])

		row = append(row, result["Good"])
		row = append(row, result["NonGood"])
		row = append(row, result["Quantity"])

		row = append(row, result["Reporter"])
		row = append(row, result["CreateAt"])
		Rows = append(Rows, row)
	}

	columns := []map[string]string{
		{"text": "WorkOrderId", "type": "string"},
		{"text": "Station", "type": "string"},
		{"text": "Machine", "type": "string"},

		{"text": "Good", "type": "string"},
		{"text": "NonGood", "type": "string"},
		{"text": "Quantity", "type": "string"},

		{"text": "Reporter", "type": "string"},
		{"text": "CreateAt", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    Rows,
		"type":    "table",
	}

	return grafanaData
}

func findWorkOrderInfo() map[string]interface{} {
	trigger := func(i interface{}) ([]byte, error) {
		url := "https://andon-daemon-compute-ifactoryandondev-eks005.sa.wise-paas.com/workorder/info"
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

	var Results []map[string]interface{}
	var Rows [][]interface{}
	Results = JsonAryToMap(res)
	for _, result := range Results {
		// var result map[string]interface{}
		var row []interface{}
		row = append(row, result["WorkOrderId"])
		row = append(row, result["Station"])

		row = append(row, result["ProductID"])
		row = append(row, result["Product"])

		row = append(row, result["Quantity"])
		row = append(row, result["Status"])
		row = append(row, result["Sum"])
		row = append(row, result["SumGood"])
		row = append(row, result["SumNonGood"])
		row = append(row, result["GoodRate"])

		row = append(row, result["PlanStartDate"])
		row = append(row, result["DeliverAt"])
		row = append(row, result["CreateAt"])
		Rows = append(Rows, row)
	}

	columns := []map[string]string{
		{"text": "WorkOrderId", "type": "string"},
		{"text": "Station", "type": "string"},

		{"text": "ProductID", "type": "string"},
		{"text": "Product", "type": "string"},

		{"text": "Quantity", "type": "string"},
		{"text": "Status", "type": "string"},
		{"text": "Sum", "type": "string"},
		{"text": "SumGood", "type": "string"},
		{"text": "SumNonGood", "type": "string"},
		{"text": "GoodRate", "type": "string"},

		{"text": "PlanStartDate", "type": "string"},
		{"text": "DeliverAt", "type": "string"},
		{"text": "CreateAt", "type": "string"},
	}

	grafanaData := map[string]interface{}{
		"columns": columns,
		"rows":    Rows,
		"type":    "table",
	}

	return grafanaData
}

func Runtest() {
	findWorkOrderInfo()
}

//-------------------------------------

type MyError struct {
	ErrDesc string
	ErrMsg  string
	Code    int
}

//只要實現了 error接口, 就可以當成是 error 型態的物件
//可以使用 return &MyError{err, 0, ""}
func (e *MyError) Error() string {
	// return fmt.Sprintf("radius %0.2f: %s", e.Code, e.Message)
	return e.ErrMsg
}

func DoAPI(apiType string, url string, param interface{}) ([]byte, error) {

	// #fix for krunshan
	req.EnableInsecureTLS(true)

	switch param.(type) {
	case nil:
		// fmt.Println(aurora.Yellow("no param"))
	case req.Param:
		// fmt.Println(aurora.Yellow("this is req.Param"))
	// case *req.bodyJson:
	// 	yellow("this is *req.bodyJson")
	default:
		// fmt.Println(aurora.Yellow("this is *req.bodyJson"))
		// yellow(reflect.TypeOf(param))
	}

	fmt.Println(aurora.Blue(fmt.Sprintf("%v %v ", apiType, url)))

	header := req.Header{
		"Accept": "application/json",
		// "Authorization": Token,
	}

	// green := color.New(color.FgGreen).SprintFunc()
	// glog.Infof("param: %+v", green(param))

	var (
		r   *req.Resp
		err error
	)

	if apiType == "GET" {
		r, err = req.Get(url, header, param)
	} else if apiType == "POST" {
		r, err = req.Post(url, header, param)
	} else if apiType == "PATCH" {
		r, err = req.Patch(url, header, param)
	} else if apiType == "PUT" {
		r, err = req.Put(url, header, param)
	} else if apiType == "DELETE" {
		r, err = req.Delete(url, header, param)
	} else {
		panic("apiType invalid")
	}

	if err != nil {
		apiErr := fmt.Errorf("API Err: %s", err.Error())
		glog.Error(Cerr(apiErr))
	}

	rCode := r.Response().StatusCode
	if rCode != 200 && rCode != 201 {
		respErr := fmt.Errorf(string(r.Bytes())) //對方api返回的錯誤訊息
		fmt.Println(aurora.Red(fmt.Errorf("%s FAIL! code=%d, resp=%s ", apiType, rCode, respErr)))
		//# New Design
		myErr := &MyError{
			ErrMsg: respErr.Error(),
			Code:   rCode,
		}
		return nil, myErr
		//how to panic and recover
	}

	// method1
	// r.ToJSON(&foo)       // response => struct/map
	// log.Printf("%+v", r) // print info (try it, you may surprise)
	// method2
	res, err := r.ToBytes()
	if err != nil {
		glog.Error(Cerr(err))
	}

	resStr := string(res)
	resStr = TruncateString(resStr, 1000)
	fmt.Println(aurora.Green(fmt.Sprintf("%s SUCCESS! code=%d, resp=%s ", apiType, rCode, resStr)))
	return res, nil
}

func Cerr(s interface{}) interface{} {
	ps := Sprintf("err: %s", s)
	return (Red(funcName() + ps))
}

//切斷超過上限的字串
func TruncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "...(TOO LONG)"
	}
	// fmt.Println(bnoden)
	return bnoden
}

func funcName() string {
	fileName, _, funcName := "???", 0, "???"
	pc, fileName, _, ok := runtime.Caller(2) //Caller(skip int) 是要提升的堆栈帧数，0-当前函数，1-上一层函数，....
	//只取呼叫來源的簡要名稱
	if ok {
		funcName = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcName = filepath.Ext(funcName)            // .foo
		funcName = strings.TrimPrefix(funcName, ".") // foo

		fileName = filepath.Base(fileName) // /full/path/basename.go => basename.go
	}
	return funcName + " "
}

//將 arrayJson 轉為 []map[string]interface
func JsonAryToMap(myjson []byte) []map[string]interface{} {

	//Marshal the json to a map
	var aryMFace []map[string]interface{}
	err := json.Unmarshal(myjson, &aryMFace)
	if err != nil {
		glog.Error(err)
	}
	return aryMFace

	// a := mmsiface["dtInstance"].(map[string]interface{})["feature"].(map[string]interface{})["monitor"]
	// fmt.Println(a)
}
