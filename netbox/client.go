package netbox

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
	"stellar.af/netbox-to-nfa/types"
	"stellar.af/netbox-to-nfa/util"
)

var httpClient *http.Client
var emptyMap map[string]string = make(map[string]string)

func init() {
	httpClient = util.CreateHTTPClient()
}

func NetboxRawRequest(m string, fu string) (gjson.Result, error) {
	empty := gjson.Result{}
	token := util.GetEnv("NETBOX_TOKEN")

	req, _ := http.NewRequest(m, fu, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

	res, err := httpClient.Do(req)

	if err != nil {
		return empty, fmt.Errorf("Request to '%s' failed:\n%s", fu, err.Error())
	}

	if res.StatusCode != 200 {
		return empty, fmt.Errorf("Error requesting data from '%s' - Status %s", fu, res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return empty, fmt.Errorf("Unable to parse response from '%s'\n%s", fu, err.Error())
	}
	return gjson.Parse(string(body)), nil
}

func NetboxRequest(m string, p string, q types.QueryParams) (gjson.Result, error) {
	baseURL := util.GetEnv("NETBOX_URL")
	u := util.BuildUrl(baseURL, p, q)
	return NetboxRawRequest(m, u.String())
}
