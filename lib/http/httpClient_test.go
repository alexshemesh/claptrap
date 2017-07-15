package httpClient

import (
	"testing"
	"github.com/bouk/monkey"
)

func Test_NewHttpExecutor(t *testing.T){
	obj := NewHttpExecutor()
	if obj == nil {
		t.Errorf("Unable to create new instance of HttpClient")
	}
}

func Test_GetPutPostDelete( t *testing.T) {
	var actualCommand string
	monkey.Patch(HttpClient.Execute, func(obj HttpClient, url string, headers map[string]string, body []byte)(response []byte, err error) {
		actualCommand = obj.requestType
		response = []byte( "200 OK" )
		return response,err
	})

	defer func() {
		monkey.Unpatch(HttpClient.Execute)
	}()

	obj := NewHttpExecutor()

	obj.Get().Execute("someurl", nil, nil)
	if actualCommand != "GET" {
		t.Errorf("Unable to process GET request")
	}

	obj.Put().Execute("someurl", nil, nil)
	if actualCommand != "PUT" {
		t.Errorf("Unable to process PUT request")
	}

	obj.Post().Execute("someurl", nil, nil)
	if actualCommand != "POST" {
		t.Errorf("Unable to process POST request")
	}

	obj.Delete().Execute("someurl", nil, nil)
	if actualCommand != "DELETE" {
		t.Errorf("Unable to process DELETE request")
	}
}

