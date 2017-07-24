package vault

import (
	"net/url"
	"path"
	"github.com/alexshemesh/claptrap/lib/http"
	"github.com/Jeffail/gabs"

	"fmt"

	"encoding/json"
	"github.com/alexshemesh/claptrap/lib/logs"
)

type DataToSet struct {
	Key string `json:"key"`
}

type VaultClient struct {
	ServerURL string
	token string
	log logs.Logger
}

func NewVaultClient(serverURL string, logPar logs.Logger) (retVal *VaultClient){
	retVal = &VaultClient{ServerURL: serverURL, log: *logPar.SubLogger("vault")}
	return retVal
}

func (this VaultClient)IsInitialized() (retVal bool, err error){
	var u *url.URL
	u, err = url.Parse(this.ServerURL)
	http := httpClient.NewHttpExecutor()
	u.Path = path.Join(u.Path, "v1", "sys", "init")
	var response []byte
	response, err = http.Get().Execute( u.String(), nil, nil )
	if err == nil {
		var jsonParsed *gabs.Container
		jsonParsed, err = gabs.ParseJSON(response)
		if err == nil {
			retVal = jsonParsed.Path("initialized").Data().(bool)
		}
	}
	return retVal,err
}

func (this VaultClient)IsSealed() (retVal bool, err error){
	var u *url.URL
	u, err = url.Parse(this.ServerURL)
	http := httpClient.NewHttpExecutor()
	u.Path = path.Join(u.Path, "v1", "sys", "seal-status")
	var response []byte
	response, err = http.Get().Execute( u.String(), nil, nil )
	if err == nil {
		var jsonParsed *gabs.Container
		jsonParsed, err = gabs.ParseJSON(response)
		if err == nil {
			retVal = jsonParsed.Path("sealed").Data().(bool)
		}
	}
	if err == nil {
		this.log.Log("Vault unsealed successfully")
	}
	return retVal,err
}

func (this VaultClient)Initialize() (keys []string, token string, err error){
	http := httpClient.NewHttpExecutor()
	body :=`{"secret_shares":1,"secret_threshold":1}`

	var u *url.URL
	u, err = url.Parse(this.ServerURL)
	u.Path = path.Join(u.Path, "v1", "sys", "init")
	var response []byte
	response, err = http.Put().Execute( u.String(), nil, []byte(body) )

	if err == nil {
		if response == nil {
			err = fmt.Errorf("No response")
		}else{
			var jsonParsed *gabs.Container
			jsonParsed, err = gabs.ParseJSON(response)

			children, _ := jsonParsed.S("keys").Children()
			for _, child := range (children ) {
				keys = append(keys, child.Data().(string))
			}

			token = jsonParsed.Path("root_token").Data().(string)
		}
	}
	if err == nil {
		this.log.Log("Vault initialized successfully")
	}
	return keys,token, err
}

func (this VaultClient)Unseal(keys []string) ( err error){
	type KeyData struct {
		Key string `json:"key"`
	}

	http := httpClient.NewHttpExecutor()
	var body []byte
	var u *url.URL
	u, err = url.Parse(this.ServerURL)
	u.Path = path.Join(u.Path, "v1", "sys", "unseal")
	for _, key := range( keys ){
		body,err = json.Marshal( KeyData{Key: key} )
		_, err = http.Put().Execute( u.String(), nil, body )
		if err != nil {
			break
		}
	}
	if err == nil {
		this.log.Log("Vault unsealed successfully")
	}

	return err
}

func  (this *VaultClient)Auth(tokenPar string){
	this.token = tokenPar
}

func (this VaultClient)Seal() ( err error){
	if this.token == "" {
		err = fmt.Errorf("Not authorized")
	} else {

		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		u.Path = path.Join(u.Path, "v1", "sys", "seal")

		_, err = http.Put().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.token }, nil )

		if err == nil {
			this.log.Log("Vault sealed successfully")
		}
	}
	return err
}

func (this VaultClient) SetVaue(secretPath string, secretValue string) ( err error){

	dataToSet := DataToSet{secretValue}
	this.log.Debug("Set secret " + secretPath)
	if this.token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		u.Path = path.Join(u.Path, "v1", "secret",secretPath)
		var body []byte
		body,err = json.Marshal(dataToSet)
		_, err = http.Put().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.token }, body )
	}
	return err
}

func (this VaultClient) GetValue(secretPath string) (retVal string, err error){
	if this.token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		u.Path = path.Join(u.Path, "v1", "secret",secretPath)
		var response []byte

		response, err = http.Get().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.token },nil )

		if err == nil {
			var parsedJson *gabs.Container
			parsedJson,err = gabs.ParseJSON(response)

			retVal = parsedJson.Path("data.key").Data().(string)
		}
	}
	return retVal,err
}
