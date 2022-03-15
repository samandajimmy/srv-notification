package nclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nlogger/v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
)

var log = nlogger.Get()

type Nclient struct {
	Client     http.Client
	WebhookURL string
}

func NewNucleoClient(url string) *Nclient {
	return &Nclient{
		Client:     http.Client{},
		WebhookURL: url,
	}
}

func (c *Nclient) PostData(header map[string]string, body map[string]interface{}) (*http.Response, error) {
	var result *http.Response

	// Get body request
	payload := getBodyRequest(header, body)

	// Make request http
	endPoint := c.WebhookURL
	request, err := http.NewRequest("POST", endPoint, payload)
	if err != nil {
		log.Errorf("Error when make new request. err: %s", err)
		return result, errx.Trace(err)
	}

	// Set header
	request = setHeaderRequest(request, header)
	log.Debugf("Request header: %s", request.Header)
	log.Debugf("Request body: %s", request.Body)

	// Do request http with client
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Errorf("Error when request client. err: %s", err)
		return result, errx.Trace(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(GetResponseString(resp))
	}

	// Set result
	result = resp

	return result, nil
}

func GetResponseString(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error while reading the response bytes:", err)
		return ""
	}
	log.Debugf(string(body))
	return string(body)
}

func getBodyRequest(header map[string]string, body map[string]interface{}) *bytes.Buffer {
	var payload *bytes.Buffer
	switch header["Content-Type"] {
	case "application/json":
		payload = setBodyApplicationJSON(body)
	case "application/x-www-form-urlencoded":
		payload = setBodyUrlEncoded(body)
	default:
		payload = setBodyApplicationJSON(body)
	}
	return payload
}

func setBodyUrlEncoded(data map[string]interface{}) *bytes.Buffer {
	var param = url.Values{}
	for key, value := range data {
		param.Set(key, nval.ParseStringFallback(value, ""))
	}

	return bytes.NewBufferString(param.Encode())
}

func setBodyApplicationJSON(data map[string]interface{}) *bytes.Buffer {
	// Set param for body request
	jsonValue, _ := json.Marshal(data)

	return bytes.NewBuffer(jsonValue)
}

func setHeaderRequest(request *http.Request, data map[string]string) *http.Request {
	// setHeaderRequest
	for key, value := range data {
		request.Header.Add(key, value)
	}
	return request
}
