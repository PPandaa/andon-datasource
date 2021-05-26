package middleware

import "os"

//docker login -u any99147 -p 54P@ssw0rd && ./build_dev.sh

var apiUrl string

func Start() {
	apiUrl = os.Getenv("API_URL")
}

func TestDatasourceFn() {
	// GetWorkOrderDetail()
	// GetStats()
	// GetWorkOrderList()
}
