package httpClient

import (
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
	"math/rand"
	"github.com/alexshemesh/claptrap/lib/logs"
	"errors"

)

type HttpClient struct {
	contentType string
	timeOut     int
	reguestID   int
	requestType string
	log         logs.Logger
}


func NewHttpExecutor() HttpExecutor {
	newObj := &HttpClient{contentType: "application/json" }
	return newObj
}

func calcRequestID() int {
	return 1000 + rand.Intn(9999-1000) //4 digit random number
}

func (this HttpClient) generateCopyWithOperation(requestType string) (HttpExecutor) {
	newVal := this
	newVal.contentType = this.contentType
	newVal.requestType = requestType
	newVal.reguestID = calcRequestID()
	newVal.log = *this.log.SubLogger(fmt.Sprintf("%s %d", newVal.requestType, newVal.reguestID))
	return newVal
}

func (this HttpClient) Get() (retVal HttpExecutor) {
	return this.generateCopyWithOperation("GET")
}

func (this HttpClient) Put() (retVal HttpExecutor) {
	return this.generateCopyWithOperation("PUT")
}

func (this HttpClient) Post() (retVal HttpExecutor) {
	return this.generateCopyWithOperation("POST")
}

func (this HttpClient) Delete() (retVal HttpExecutor) {
	return this.generateCopyWithOperation("DELETE")
}

func (this HttpClient) Execute(url string, headers map[string]string, body []byte) (response []byte, err error) {
	requestID := calcRequestID()

	client := &http.Client{}

	if headers == nil {
		headers = make(map[string]string)
	}



	req, err := http.NewRequest(this.requestType, url, bytes.NewBuffer(body))
	if err == nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		this.log.Debug("Set content type to be " + this.contentType)
		req.Header.Set("Content-Type", this.contentType)
		this.log.Debug(fmt.Sprintf("Executing request URL: %s, headers: %s, body: %s", url, headers, body))
		var resp *http.Response
		resp, err = client.Do(req)
		if err == nil {
			if resp != nil {
				this.log.Log(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status))

				if resp.StatusCode >= 400 {
						err = errors.New(resp.Status)
				}
			}
		}

		defer func() {
			if resp != nil {
				resp.Body.Close()
			}
		}()


		if err != nil {
			this.log.Error(fmt.Errorf("[%d] HTTP request failed", requestID))
			this.log.Error(err)
		} else {
			response, _ = ioutil.ReadAll(resp.Body)
		}
	}
	return response, err
}

func (this HttpClient)SetContentType(newContentType string)( HttpClient){
	newVal := this
	newVal.contentType = newContentType
	newVal.reguestID = calcRequestID()
	newVal.log = *this.log.SubLogger(fmt.Sprintf("%s %d", newVal.requestType, newVal.reguestID))
	return newVal
}
