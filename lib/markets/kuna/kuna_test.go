package kuna

import (
	"testing"
	"github.com/alexshemesh/claptrap/lib/logs"
	"io"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"path"
)

func Test_NewKunaClient(t *testing.T){
	log := logs.NewLogger("test kuna")
	kunaClient := NewKunaClient(*log)
	if kunaClient == nil {
		t.Errorf("Cannot create Kuna client")
	}
}

func TestKunaClient_GetOrdersBook(t *testing.T) {

	pathUrl := "/api/v2/order_book?market=btcuah"
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		pathToCheck := pathUrl
		if r.URL.Path == pathToCheck && r.Method == "GET" {
			data, _ := ioutil.ReadFile(path.Join("kuna_testdata", "orderbooks.json"))
			io.WriteString(w, string(data))
		}else{
			io.WriteString(w, `500 Wrong request path`)
			r.Response.StatusCode = 500
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	kunaClient := NewKunaClient(*logs.NewLogger("test kuna client"))
	kunaClient.ServerURL = ts.URL

	book,err := kunaClient.GetOrdersBook()
	if err != nil {
		t.Error(err)
	}
	if len( book.Asks ) == 0 {
		t.Errorf("Cannot get asks")
	}
	if len( book.Bids ) == 0 {
		t.Errorf("Cannot get bids")
	}
}