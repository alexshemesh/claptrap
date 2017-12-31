package quoinex

import (
	"github.com/alexshemesh/claptrap/lib/logs"
	"net/url"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"
	"github.com/alexshemesh/claptrap/lib/contracts"
	httphelper "github.com/alexshemesh/claptrap/lib/http"
	"strings"

	"net/http"
)

type QuoineApi struct{
	log logs.Logger
	apikey string
	secret string
	ApiUrl string
}



// getSha256 creates a sha256 hash for given []byte
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func createSignature(urlPath string, values url.Values, secret []byte) string {
	// See https://www.kraken.com/help/api#general-usage for more information
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macSum := getHMacSha512(append([]byte(urlPath), shaSum...), secret)
	return base64.StdEncoding.EncodeToString(macSum)
}




func NewQuoineClient(logsPar logs.Logger, apiKeyPar string, apiSecretPar string) (retVal *QuoineApi){
	retVal = &QuoineApi{
		log: logsPar,
		apikey: apiKeyPar,
		secret: apiSecretPar,
		ApiUrl: "https://api.quoine.com",
	}
	return retVal
}

func (this QuoineApi) doRequest(reqURL string, values url.Values, headers map[string]string) ([]byte, error) {

	// Create request
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #1 (%s)", err.Error())
	}


	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Execute request
	resp, err := this.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #2 (%s)", err.Error())
	}
	defer resp.Body.Close()

	// Read request
	body, err := ioutil.ReadAll(resp.Body)
}

func (this QuoineApi)performGetRequest(urlPath string, values url.Values){

	reqURL := fmt.Sprintf("%s%s", APIURL, urlPath)
	secret, _ := base64.StdEncoding.DecodeString(this.secret)
	values.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))


	// Create signature
	signature := createSignature(urlPath, values, secret)
	// Add Key and signature to request headers
	headers := map[string]string{
		"API-Key":  this.apikey,
		"API-Sign": signature,
		"X-Quoine-API-Version":	"2",
	}


	resp, err := this.doRequest(reqURL, values, headers)
}

func (this QuoineApi)parseBalance()(retVal contracts.Balance, err error){

	return retVal,err
}
func (this QuoineApi)GetBalance(){
	urlPath := "/accounts/balance"
	secret, _ := base64.StdEncoding.DecodeString(this.secret)
	var values url.Values
	values.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	signature := createSignature(urlPath, values, secret)
	headers := map[string]string{
		"API-Key":  this.apikey,
		"API-Sign": signature,
		"X-Quoine-API-Version":	"2",
	}

	reqURL := fmt.Sprintf("%s%s", this.ApiUrl, urlPath)


	httpExecutor := httphelper.NewHttpExecutor()
	resp, err := httpExecutor.Get().ExecuteWithReader(reqURL, headers, strings.NewReader(values.Encode()))

}