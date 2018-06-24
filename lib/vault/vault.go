package vault

import (
	"net/url"
	"path"
	"github.com/alexshemesh/claptrap/lib/http"
	"github.com/Jeffail/gabs"

	"fmt"

	"encoding/json"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/spf13/viper"

	garbler "github.com/michaelbironneau/garbler/lib"

	"github.com/alexshemesh/claptrap/lib/types"

	"strings"
)

type DataToSet struct {
	Key string `json:"key"`
}

type VaultClient struct {
	ServerURL  string
	Token      string
	UserDomain string
	log        logs.Logger
}

func calcPolicyName(userName string)(retVal string){
	retVal = userName + "_policy"
	return retVal
}

func NewVaultClientInitialized( logPar logs.Logger) (retVal *VaultClient){
	vaultAddr := viper.Get("vault.url").(string)
	vaultToken := viper.Get("vault.Token").(string)

	retVal = &VaultClient{ServerURL: vaultAddr, log: *logPar.SubLogger("vault"), Token: vaultToken}
	return retVal
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

func  (this *VaultClient) SetToken(tokenPar string){
	this.Token = tokenPar
}

func (this VaultClient)Seal() ( err error){
	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {

		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		u.Path = path.Join(u.Path, "v1", "sys", "seal")

		_, err = http.Put().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token}, nil )

		if err == nil {
			this.log.Log("Vault sealed successfully")
		}
	}
	return err
}

func (this VaultClient) SetValue(secretPath string, secretValue string) ( err error){

	dataToSet := DataToSet{secretValue}
	this.log.Debug("Set secret " + secretPath)
	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		u.Path = path.Join(u.Path, "v1", "secret",secretPath)
		var body []byte
		body,err = json.Marshal(dataToSet)
		_, err = http.Put().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token}, body )
	}
	return err
}

func (this VaultClient) GetValue(secretPath string) (retVal string, err error){
	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		if this.UserDomain != "" {
			u.Path = path.Join(u.Path, "v1", "secret",this.UserDomain,secretPath)
		}else{
			u.Path = path.Join(u.Path, "v1", "secret",secretPath)
		}

		var response []byte

		response, err = http.Get().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token},nil )

		if err == nil {
			var parsedJson *gabs.Container
			parsedJson,err = gabs.ParseJSON(response)

			retVal = parsedJson.Path("data.key").Data().(string)
		}
	}
	return retVal,err
}


func (this VaultClient)CreateUser( userName string ) (initialPassword string, err error){
	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		initialPassword, _ = garbler.NewPassword(nil)


	  body :=fmt.Sprintf(`{"password":"%s", "policies":"%s"}`, initialPassword ,calcPolicyName(userName) )
	//	body :=fmt.Sprintf(`{"password":"%s", "policies":"admin,default"}`, initialPassword )


		u.Path = path.Join(u.Path, "v1", "auth","userpass","users",userName)

		_, err = http.Post().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token},[]byte(body) )

	}
	return initialPassword, err
}


func (this VaultClient)CreateDefaultPolicyForUser( userName string ) (  err error){


	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		body := fmt.Sprintf(`{"rules":"{\"path\":{\"auth\/approle\/role\/\":{\"policy\":\"write\"},\"secret\/%s\/*\":{\"policy\":\"read\"}}}"}`,userName)
		/*body := `{
								"rules":"{
										\"path\": {
												\"auth\/approle\/role\/\":{
														\"policy\":\"write\"
												},
												\"secret\/\":{
														\"policy\":\"write\"}
												}
										}"
							}`
*/
		//	body := `{"rules":"path \"postgresql/creds/readonly\" {\n  capabilities = [\"read\"]\n}"}`
	/*	policyTemplate := `{"rules":{
			path "secret/*" {
				capabilities = ["deny"]
			}

			path "sys/*" {
				capabilities = ["deny"]
			}

			path "secret/users/%s/*" {
				capabilities = ["create", "read", "update", "delete", "list"]
			}

			path "secret/settings/*" {
				capabilities = [ "read" ]
			}

			path "auth/Token/lookup-self" {
				capabilities = [ "read" ]
			}
		}`

		body := fmt.Sprintf(policyTemplate,userName )
*/
		u.Path = path.Join(u.Path, "v1", "sys","policy", calcPolicyName(userName) )
		//http = http.SetContentType("text/plain")

		_, err = http.Put().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token},[]byte(body) )

	}
	return err
}


func (this VaultClient)LogInCached(userName string)(retObjSettings types.Settings,retObjAuth types.Auth ,err error){
	var token,tokenName string
	token,err = this.GetValue("tokens/" + userName)
	tokenName,err = this.LookupToken(token)
	retVal := NewVaultClientInitialized(this.log)
	if err == nil {
		this.log.Debug("Looked up token is OK" + tokenName)

		retVal.Token = token
		retVal.UserDomain = userName
	}else{
		this.log.Error(err)
	}
	return retVal ,retVal  ,err
}

func (this VaultClient)RenewToken()(err error){
	return err
}

func (this VaultClient)LogIn( userName string, password string ) ( retObjSettings types.Settings,retObjAuth types.Auth ,err error){
	retVal := this
	http := httpClient.NewHttpExecutor()

	body := fmt.Sprintf(`{"password": "%s"}`, password)

	var u *url.URL
	u, err = url.Parse(this.ServerURL)

	u.Path = path.Join(u.Path, "v1", "auth", "userpass","login",userName )

	var resp []byte
	resp, err = http.Post().Execute( u.String(), nil,[]byte(body) )
	if err == nil {
		if  resp != nil {
			this.log.Debug(string(resp))
			var respData *gabs.Container

			respData,err = gabs.ParseJSON(resp)
			if err == nil {
				//setting up new Token
				retVal.Token = strings.Trim(respData.Path("auth.client_token").String(),"\"")
				retVal.UserDomain = userName
			}
		}
	}
	if err != nil{
		this.log.Error(err)
	}

	return retVal,retVal,err
}

func (this VaultClient) SaveCachedToken(userName string , token string) (err error){
	err = this.SetValue("tokens/" + userName, token)
	return err
}

func (this VaultClient) LookupToken(token string) (retVal string, err error){
	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)
		u.Path = path.Join(u.Path, "v1", "auth","token","lookup")
		body := fmt.Sprintf(`{"token": "%s"}`, token)
		var response []byte
		response, err = http.Get().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token},[]byte(body) )
		this.log.Debug(string(response))
		if err == nil {
			var jsonParsed *gabs.Container
			jsonParsed, err = gabs.ParseJSON(response)
			if err == nil {
				retVal = jsonParsed.Path("data.display_name").Data().(string)
			}
		}
	}
	return retVal, err
}



func (this VaultClient)DeleteUserPolicy(userName string)(err error){
	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)

		u.Path = path.Join(u.Path, "v1", "sys","policy", calcPolicyName(userName) )

		_, err = http.Delete().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token},nil )

	}
	return err
}

func (this VaultClient)DeleteUser( userName string )  (err error){
	if this.Token == "" {
		err = fmt.Errorf("Not authorized")
	} else {
		http := httpClient.NewHttpExecutor()

		var u *url.URL
		u, err = url.Parse(this.ServerURL)

		u.Path = path.Join(u.Path,"v1", "auth","userpass","users",userName)

		_, err = http.Delete().Execute( u.String(), map[string]string{ "X-Vault-Token":" " + this.Token},nil )

		err = this.DeleteUserPolicy(userName)

	}
	return err
}

func (this VaultClient)GetUserName()( string){
	return this.UserDomain
}

func (this VaultClient)GetUserToken()( string){
	return this.Token
}
