package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/talbx/go-clockodo/cmd/intercept"
)

const apiRoot string = "https://my.clockodo.com/api/"

func CallApi[R TimeEntriesResponse | ClockResponse | Customer](endpoint string, strukt *R) int {
	return getAndDelete("GET", endpoint, strukt)
}

type StartPayload struct {
	CustomerId  int    `json:"customers_id"`
	ServiceId   int    `json:"services_id"`
	ProjectId   int    `json:"projects_id"`
	Description string `json:"text"`
}

func getAndDelete[R TimeEntriesResponse | ClockResponse | Customer, P StartPayload](requestMethod string, endpoint string, strukt *R) int {
	client := &http.Client{}

	config := intercept.ClockodoConfig
	req, err := http.NewRequest(requestMethod, fmt.Sprintf("%s%s", apiRoot, endpoint), nil)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("X-ClockodoApiUser", config.ApiUser)
	req.Header.Add("X-ClockodoApiKey", config.ApiKey)
	req.Header.Add("X-Clockodo-External-Application", "clockodo-cli")
	//res, _ := httputil.DumpRequest(req, true)
	//fmt.Println(string(res))
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	if nil != strukt {
		json.NewDecoder(resp.Body).Decode(strukt)
	}

	return resp.StatusCode
}
