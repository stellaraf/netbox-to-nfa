package nfa

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
	"stellar.af/netbox-to-nfa/types"
	"stellar.af/netbox-to-nfa/util"
)

var httpClient *http.Client
var emptyMap types.QueryParams = make(map[string]interface{})
var emptyResult gjson.Result = gjson.Result{}

func init() {
	httpClient = util.CreateInsecureHTTPClient()
}

// NFAAuth authenticates with the NFA API and returns the authenticated client.
func NFAAuth(client *http.Client) (*http.Client, error) {
	baseURL := util.GetEnv("NFA_URL")
	username := util.GetEnv("NFA_USERNAME")
	password := util.GetEnv("NFA_PASSWORD")
	body := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	u := util.BuildUrl(baseURL, "/api/login", emptyMap)
	req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error authenticating to '%s' failed:\n%s", u, err.Error())
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error authenticating to '%s' failed:\n%s", u, res.Status)
	}
	return client, nil
}

// NFARawRequest sends an HTTP request to NFA and returns the parsed GJSON result.
func NFARawRequest(m string, fu string, b *[]byte) (gjson.Result, error) {

	client, err := NFAAuth(httpClient)

	if err != nil {
		return emptyResult, err
	}

	var reqBody bytes.Buffer

	if b != nil {
		reqBody = *bytes.NewBuffer(*b)
	}

	req, _ := http.NewRequest(m, fu, &reqBody)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return emptyResult, fmt.Errorf("Request to '%s' failed:\n%s", fu, err.Error())
	}

	if res.StatusCode > 399 {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return emptyResult, fmt.Errorf("Error sending %s request to '%s' - %s - Detail:\n%s", m, fu, res.Status, string(body))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return emptyResult, fmt.Errorf("Unable to parse response from '%s'\n%s", fu, err.Error())
	}
	return gjson.Parse(string(body)), nil
}

// NFARequest is an abstraction around NFARawRequest for easier code-writing.
func NFARequest(m string, p string, q types.QueryParams, b *[]byte) (gjson.Result, error) {
	baseURL := util.GetEnv("NFA_URL")
	u := util.BuildUrl(baseURL, p, q)
	return NFARawRequest(m, u.String(), b)
}
