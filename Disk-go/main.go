package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//TODO: APIC username / password
var username string = "xxxxxxxxx"
var pwd string = "pwdpwdpwdpwd"

//TODO: your apic Server web api url (change 127.0.0.1 )
var loginEndpoint string = "https://127.0.0.1/api/class/aaaLogin.json"
var eqptStorageEndpoint string = "https://127.0.0.1/api/class/eqptStorage.json"

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

type EqptStorage struct {
	TotalCount string `json:"totalCount"`
	Imdata     []struct {
		EqptStorage struct {
			Attributes struct {
				Available       string    `json:"available"`
				Blocks          string    `json:"blocks"`
				CapUtilized     string    `json:"capUtilized"`
				ChildAction     string    `json:"childAction"`
				Device          string    `json:"device"`
				Dn              string    `json:"dn"`
				FailReason      string    `json:"failReason"`
				FileSystem      string    `json:"fileSystem"`
				FirmwareVersion string    `json:"firmwareVersion"`
				LcOwn           string    `json:"lcOwn"`
				MediaWearout    string    `json:"mediaWearout"`
				ModTs           time.Time `json:"modTs"`
				Model           string    `json:"model"`
				MonPolDn        string    `json:"monPolDn"`
				Mount           string    `json:"mount"`
				Name            string    `json:"name"`
				NameAlias       string    `json:"nameAlias"`
				OperSt          string    `json:"operSt"`
				Serial          string    `json:"serial"`
				Status          string    `json:"status"`
				Used            string    `json:"used"`
			} `json:"attributes"`
		} `json:"eqptStorage"`
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

// output Pretty json
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
	body, _ := ioutil.ReadAll(resp.Body)

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
	url = eqptStorageEndpoint

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
	// with body get Disk Capacity
	eqptBodyString := string(body)
	var eqptResult EqptStorage
	json.Unmarshal([]byte(eqptBodyString), &eqptResult)
	var count, _ = strconv.Atoi(eqptResult.TotalCount) // Convert the total obtained to int
	mapData := make(map[string]string)                 // Use map to organize data
	// Dynamically put the captured disk path and Disk Capacity into the map
	for i := 0; i < count-1; i++ {
		mapData[eqptResult.Imdata[i].EqptStorage.Attributes.Dn] = eqptResult.Imdata[i].EqptStorage.Attributes.CapUtilized
		//fmt.Println(mapData[eqptResult.Imdata[i].EqptStorage.Attributes.Dn]) // debug
		//fmt.Println(eqptResult.Imdata[i].EqptStorage.Attributes.CapUtilized) // debug
	}

	// ==============================
	// Start crawling the required
	// TODO: Here you need to modify it to the disk path you want to grab
	// Node1: data / data2
	node1Data, ok := mapData["topology/pod-1/node-1/sys/ch/p-[/data]-f-[/dev/mapper/vg_ifc0_ssd-data]"]
	if !ok {
		log.Fatal("It node1Data1 should be true")
	}
	node1Data2, ok := mapData["topology/pod-1/node-1/sys/ch/p-[/data2]-f-[/dev/mapper/vg_ifc0-data2]"]
	if !ok {
		log.Fatal("It node1Data1 should be true")
	}
	// Node2: data / data2
	node2Data, ok := mapData["topology/pod-1/node-2/sys/ch/p-[/data]-f-[/dev/mapper/vg_ifc0_ssd-data]"]
	if !ok {
		log.Fatal("It node1Data1 should be true")
	}
	node2Data2, ok := mapData["topology/pod-1/node-2/sys/ch/p-[/data2]-f-[/dev/mapper/vg_ifc0-data2]"]
	if !ok {
		log.Fatal("It node1Data1 should be true")
	}
	// Node3: data / data2
	node3Data, ok := mapData["topology/pod-1/node-3/sys/ch/p-[/data]-f-[/dev/mapper/vg_ifc0_ssd-data]"]
	if !ok {
		log.Fatal("It node1Data1 should be true")
	}
	node3Data2, ok := mapData["topology/pod-1/node-3/sys/ch/p-[/data2]-f-[/dev/mapper/vg_ifc0-data2]"]
	if !ok {
		log.Fatal("It node1Data1 should be true")
	}

	// ==============================
	// Output format data to a PRTG
	var prtgOutputString = `{"prtg": {"result": [{"channel": "Node1/Data DiskCapacity","value": %s},{"channel": "Node1/Data2 DiskCapacity","value": %s},{"channel": "Node2/Data DiskCapacity","value": %s},{"channel": "Node2/Data2 DiskCapacity","value": %s},{"channel": "Node3/Data DiskCapacity","value": %s},{"channel": "Node3/Data2 DiskCapacity","value": %s}]}}`
	fmt.Printf(jsonPrettyPrint(prtgOutputString)+"\n", node1Data, node1Data2, node2Data, node2Data2, node3Data, node3Data2)
}
