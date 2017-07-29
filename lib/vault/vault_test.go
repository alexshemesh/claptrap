package vault

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io"
	"github.com/alexshemesh/claptrap/lib/logs"

)

func Test_vaultIsInitialized_true(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/init" {
			io.WriteString(w, `{"initialized":true}`)
		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	initialized,err := vault.IsInitialized()
	if err != nil {
		t.Error(err)
	}
	if initialized == false {
		t.Errorf("Cannot parse init response")
	}

}

func Test_vaultIsInitialized_false(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/init" {
			io.WriteString(w, `{"initialized":false}`)
		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	initialized,err := vault.IsInitialized()
	if err != nil {
		t.Error(err)
	}
	if initialized == true {
		t.Errorf("Cannot parse init response.Suppose to be false")
	}
}

func Test_vaulIsSealed_true(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/seal-status" {
			io.WriteString(w, `{"sealed":true}`)
		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	sealed,err := vault.IsSealed()
	if err != nil {
		t.Error(err)
	}
	if sealed == false {
		t.Errorf("Cannot parse seal response.Suppose to be true")
	}
}

func Test_vaulInitialize(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/init" && r.Method == "PUT" {
			io.WriteString(w, `{"keys":["12345678123456781234567807582a552e10c713423f43d3ebfcdaf80a18b83f"],"keys_base64":["somekeymashupverylong"],"root_token":"12345678-ffff-ffff-ffff-fe38b7453f03"}`)
		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	keys,token,err := vault.Initialize()
	if err != nil {
		t.Error(err)
	}
	if len(keys) != 1  {
		t.Errorf("Wrong number of keys")
	}

	if token == "" {
		t.Errorf("No root token")
	}
}

func Test_vaulUnseal(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/unseal" && r.Method == "PUT" {
			io.WriteString(w, `200 OK`)
		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	err := vault.Unseal([]string{"fffff-ffff"})
	if err != nil {
		t.Error(err)
	}
}

func Test_vaulSeal_notauthorized(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/seal" && r.Method == "PUT" {
			token := r.Header.Get("X-Vault-Token")
			if token != "" {
				io.WriteString(w, `200 OK`)
			}else{
				w.WriteHeader(http.StatusNotFound)
			}

		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	err := vault.Seal()
	if err == nil {
		t.Error("Cannot detect not authorized request")
	}
}

func Test_vaulSeal(t *testing.T){
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/seal" && r.Method == "PUT" {
			token := r.Header.Get("X-Vault-Token")
			if token != "" {
				io.WriteString(w, `200 OK`)
			}else{
				w.WriteHeader(http.StatusNotFound)
			}

		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	vault.Auth("token")
	err := vault.Seal()
	if err != nil {
		t.Error(err)
	}
}

func Test_vaulSetSecret(t *testing.T){
	path := "path/to/secret"
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		pathToCheck := "/v1/secret/" + path
		if r.URL.Path == pathToCheck && r.Method == "PUT" {
			token := r.Header.Get("X-Vault-Token")
			if token != "" {
				io.WriteString(w, `200 OK`)
			}else{
				w.WriteHeader(http.StatusNotFound)
			}

		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	vault.Auth("token")
	err := vault.SetValue(path,"secretdata")
	if err != nil {
		t.Error(err)
	}
}

func Test_vaulGetSecret(t *testing.T){
	path := "path/to/secret"
	secretValue := `{"data":{"key":"somevalue"}}`
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/secret/" + path && r.Method == "GET" {
			token := r.Header.Get("X-Vault-Token")
			if token != "" {
				io.WriteString(w, secretValue)
			}else{
				w.WriteHeader(http.StatusNotFound)
			}

		}else{
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	vault := NewVaultClient(ts.URL,*logs.NewLogger("root"))
	vault.Auth("token")
	response,err := vault.GetValue(path)
	if err != nil {
		t.Error(err)
	}

	if response != "somevalue" {
		t.Errorf("Cannot retrive secret")
	}
}