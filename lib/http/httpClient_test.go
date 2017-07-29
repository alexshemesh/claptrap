package httpClient

import (
	"testing"
	"github.com/bouk/monkey"
	"io"
	"net/http/httptest"
	"net/http"

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


func Test_NILForHeder(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, `200 OK`)
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()


	obj := NewHttpExecutor()
	_,err := obj.Get().Execute(ts.URL,nil,nil)

	if err != nil {
		t.Error(err)
	}
}

func Test_HttpErrorToRealError(t *testing.T){
		echoHandler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}

		// create test server with handler
		ts := httptest.NewServer(http.HandlerFunc(echoHandler))
		defer ts.Close()


		obj := NewHttpExecutor()
		_,err := obj.Get().Execute(ts.URL,nil,nil)

		if err == nil {
			t.Errorf("Cannot detect HTTP error")
		}

		if err.Error() != "400 Bad Request" {
			t.Errorf("Wrong error code. Suppose to be 400")
		}
}