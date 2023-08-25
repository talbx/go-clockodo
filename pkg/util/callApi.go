package util

import (
	"encoding/json"
	"fmt"
	"github.com/talbx/go-clockodo/pkg/intercept"
	"github.com/talbx/go-clockodo/pkg/model"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const apiRoot string = "https://my.clockodo.com/api/"

func CallApi[R model.TimeEntriesResponse | model.ClockResponse | model.Customer](endpoint string, strukt *R) (int, error) {
	s := time.Now()
	ep := strings.Split(endpoint, "?")[0]
	slog.Debug(fmt.Sprintf("[CallApi] calling clocko:do API %v", ep))
	response, err := getAndDelete("GET", endpoint, strukt)
	e := time.Now()
	duration := e.Sub(s).Milliseconds()
	slog.Debug(fmt.Sprintf("[CallApi] done calling clocko:do API %v in %vms", ep, duration))
	return response, err
}

func getAndDelete[R model.TimeEntriesResponse | model.ClockResponse | model.Customer](requestMethod string, endpoint string, strukt *R) (int, error) {
	client := &http.Client{}

	config := intercept.ClockodoConfig
	req, err := http.NewRequest(requestMethod, fmt.Sprintf("%s%s", apiRoot, endpoint), nil)
	if err != nil {
		slog.Error(fmt.Sprint(err))
	}
	req.Header.Add("X-ClockodoApiUser", config.ApiUser)
	req.Header.Add("X-ClockodoApiKey", config.ApiKey)
	req.Header.Add("X-Clockodo-External-Application", "clockodo-cli")
	//res, _ := httputil.DumpRequest(req, true)
	//fmt.Println(string(res))
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(fmt.Sprint(err))
	}

	if nil != strukt {
		err = json.NewDecoder(resp.Body).Decode(strukt)
	}

	return resp.StatusCode, err
}
