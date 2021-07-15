package meraki

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
	"stellar.af/meraki-reboot/types"
	"stellar.af/meraki-reboot/util"
)

const MERAKI_API_URL string = "https://api.meraki.com"

var httpClient *http.Client
var emptyMap map[string]string = make(map[string]string)
var emptyQuery map[string]interface{} = make(map[string]interface{})

func init() {
	httpClient = util.CreateHTTPClient()
}

func merakiRawRequest(m string, fu string) (gjson.Result, error) {
	empty := gjson.Result{}
	token := util.GetEnv("MERAKI_API_KEY")

	req, _ := http.NewRequest(m, fu, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-cisco-meraki-api-key", token)

	res, err := httpClient.Do(req)

	if err != nil {
		return empty, fmt.Errorf("Request to '%s' failed:\n%s", fu, err.Error())
	}

	if res.StatusCode > 399 {
		return empty, fmt.Errorf("Error requesting data from '%s' - Status %s", fu, res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return empty, fmt.Errorf("Unable to parse response from '%s'\n%s", fu, err.Error())
	}
	return gjson.Parse(string(body)), nil
}

func MerakiRequest(m string, p string, q types.QueryParams) (gjson.Result, error) {
	u := util.BuildUrl(MERAKI_API_URL, p, q)
	return merakiRawRequest(m, u.String())
}
