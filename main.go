package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/vefaas"
	"github.com/volcengine/vefaas-golang-runtime/vefaascontext"
)

func main() {
	// Start your vefaas function =D.
	vefaas.Start(handler)
}

// Define your handler function.
func handler(ctx context.Context, r *events.HTTPRequest) (*events.EventResponse, error) {
	fmt.Printf("received new request: %s %s, request id: %s\n", r.HTTPMethod, r.Path, vefaascontext.RequestIdFromContext(ctx))
	fmt.Printf("debug request: header:%v, body:%s\n", r.Headers, r.Body)
	ret := make(map[string]interface{})
	ret["content"] = "Hello From QA YTR Test For gateway!"
	ret["http_method"] = r.HTTPMethod
	ret["http_path"] = r.Path
	query := make(map[string]interface{})
	header := make(map[string]interface{})
	for k, v := range r.QueryStringParameters {
		query[k] = v
	}
	for k, v := range r.Headers {
		header[k] = v
	}
	ret["http_query"] = query
	ret["http_header"] = header
	ret["http_body"] = string(r.Body)
	openapi_body := AccessOpenApi()
	ret["openapi_body"] = openapi_body
	retBody, _ := json.Marshal(ret)
	return &events.EventResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: retBody,
	}, nil
}

func AccessOpenApi() string {

	//url := "https://developer.toutiao.com/api/apps/qrcode"
	host := "dev.douyincloud.gateway.egress.ivolces.com"
	url := fmt.Sprintf("http://%s/api/v2/tags/text/antidirt", host)
	method := "POST"

	//payload := strings.NewReader(`{"access_token": "0801121846765a5a4d2f6b385a68307237534d43397a667865513d3d","appname": "douyin"}`)
	payloadWithoutToken := strings.NewReader(`{"tasks": [{"content": "要检测的文本"}]}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadWithoutToken)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	resheader := res.Header
	log.Printf("ytr test raw request:%+v\n\n", req)
	log.Printf("ytr test resp from openapi:%+v,%+v\n\n", string(body), resheader)
	return string(body)
}
