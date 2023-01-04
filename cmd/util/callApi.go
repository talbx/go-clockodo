package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const apiRoot string = "https://my.clockodo.com/api/"

func CallApi[R TimeEntriesResponse | ClockResponse](endpoint string, strukt *R) {
	client := &http.Client{}
	var config GoClockodoConfig
	ReadConfig(&config)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", apiRoot, endpoint), nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(req)
	req.Header.Add("X-ClockodoApiUser", config.ApiUser)
	req.Header.Add("X-ClockodoApiKey", config.ApiKey)
	req.Header.Add("X-Clockodo-External-Application", "clockodo-cli")
	resp, err := client.Do(req)
	fmt.Print(resp.Status)
	if err != nil {
		log.Panic(err)
	}

	json.NewDecoder(resp.Body).Decode(strukt)

}