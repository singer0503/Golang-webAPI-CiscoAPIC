package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//TODO: APIC username / password
var username string = "xxxxxxxxx"
var pwd string = "pwdpwdpwdpwd"

//TODO: your apic Server web api url
var loginEndpoint string = "https://172.16.253.101/api/class/aaaLogin.json"
var cpuEndpoint string = "https://172.16.253.101/api/node/class/procEntity.json"

type LoginResult struct {
	TotalCount string `json:"totalCount"`
	Imdata     []struct {
		AaaLogin struct {
			Attributes struct {
				Token                  string `json:"token"`
				SiteFingerprint        string `json:"siteFingerprint"`
				RefreshTimeoutSeconds  string `json:"refreshTimeoutSeconds"`
				MaximumLifetimeSeconds string `json:"maximumLifetimeSeconds"`
				GuiIdleTimeoutSeconds  string `json:"guiIdleTimeoutSeconds"`
				RestTimeoutSeconds     string `json:"restTimeoutSeconds"`
				CreationTime           string `json:"creationTime"`
				FirstLoginTime         string `json:"firstLoginTime"`
				UserName               string `json:"userName"`
				RemoteUser             string `json:"remoteUser"`
				UnixUserID             string `json:"unixUserId"`
				SessionID              string `json:"sessionId"`
				LastName               string `json:"lastName"`
				FirstName              string `json:"firstName"`
				ChangePassword         string `json:"changePassword"`
				Version                string `json:"version"`
				BuildTime              string `json:"buildTime"`
				Node                   string `json:"node"`
			} `json:"attributes"`
			Children []struct {
				AaaUserDomain struct {
					Attributes struct {
						Name   string `json:"name"`
						RolesR string `json:"rolesR"`
						RolesW string `json:"rolesW"`
					} `json:"attributes"`
					Children []struct {
						AaaReadRoles struct {
							Attributes struct {
							} `json:"attributes"`
						} `json:"aaaReadRoles,omitempty"`
						AaaWriteRoles struct {
							Attributes struct {
							} `json:"attributes"`
							Children []struct {
								Role struct {
									Attributes struct {
										Name string `json:"name"`
									} `json:"attributes"`
								} `json:"role"`
							} `json:"children"`
						} `json:"aaaWriteRoles,omitempty"`
					} `json:"children"`
				} `json:"aaaUserDomain,omitempty"`
				DnDomainMapEntry struct {
					Attributes struct {
						Dn              string `json:"dn"`
						ReadPrivileges  string `json:"readPrivileges"`
						WritePrivileges string `json:"writePrivileges"`
					} `json:"attributes"`
				} `json:"DnDomainMapEntry,omitempty"`
			} `json:"children"`
		} `json:"aaaLogin"`
	} `json:"imdata"`
}

type ProcEntity struct {
	TotalCount string `json:"totalCount"`
	Imdata     []struct {
		ProcEntity struct {
			Attributes struct {
				AdminSt     string    `json:"adminSt"`
				ChildAction string    `json:"childAction"`
				CPUPct      string    `json:"cpuPct"`
				Dn          string    `json:"dn"`
				MaxMemAlloc string    `json:"maxMemAlloc"`
				MemFree     string    `json:"memFree"`
				ModTs       time.Time `json:"modTs"`
				MonPolDn    string    `json:"monPolDn"`
				Name        string    `json:"name"`
				OperErr     string    `json:"operErr"`
				OperSt      string    `json:"operSt"`
				Status      string    `json:"status"`
			} `json:"attributes"`
		} `json:"procEntity"`
	} `json:"imdata"`
}

type PrtgOutput struct {
	Prtg struct {
		Result []struct {
			Channel string `json:"channel"`
			Value   int    `json:"value"`
		} `json:"result"`
	} `json:"prtg"`
}

//  output Pretty json
func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "    ")
	if err != nil {
		return in
	}
	return out.String()
}

func main() {
	// disable security checks globally for all requests of the default client:
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// ======================================================
	// Use the post method to log in
	url := loginEndpoint
	var jsonStr = []byte("{\"aaaUser\":{\"attributes\":{\"name\":\"" + username + "\",\"pwd\":\"" + pwd + "\"}}}")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // Get response body data

	// ======================================================
	// Take the token value out of the body (XXXXXXX...)
	loginBodyString := string(body)
	var loginResult LoginResult
	json.Unmarshal([]byte(loginBodyString), &loginResult)
	token := loginResult.Imdata[0].AaaLogin.Attributes.Token

	// ======================================================
	// The next Get Request Header must bring token
	// Sample ( Cookie : "APIC-cookie=" + XXXXXXX...)
	// If necessary, add ( ; Path=/; Domain=172.16.253.101; Secure; HttpOnly; )
	url = cpuEndpoint

	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "APIC-cookie="+token)
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	// ==============================
	// Take cpuPct (CPU usage rate) out of body
	cpuBodyString := string(body)
	var cpuResult ProcEntity
	json.Unmarshal([]byte(cpuBodyString), &cpuResult)
	var count, _ = strconv.Atoi(cpuResult.TotalCount) // Convert the total obtained to int
	var array = make([]string, count)                 // Generate dynamic array to store data

	// Put the captured CPU value into the dynamic array
	for i, _ := range array {
		array[i] = cpuResult.Imdata[i].ProcEntity.Attributes.CPUPct
		//fmt.Println("array[" + strconv.Itoa(i) + "] = " + array[i]) // debug
	}

	// ==============================
	// Output format data to a PRTG
	var prtgOutputString = `{"prtg": {"result": [{"channel": "Node-1","value": %s},{"channel": "Node-2","value": %s},{"channel": "Node-3","value": %s}]}}`
	fmt.Printf(jsonPrettyPrint(prtgOutputString)+"\n", array[0], array[1], array[2])
}
