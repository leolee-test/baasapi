package cron

import (
	"net/http"
	"fmt"
	"log"
	"bytes"
	"strings"
	"time"
	"io/ioutil"
	"encoding/json"

	"github.com/baasapi/baasapi/api"
)

type jsonResponse struct {
    Success    bool    `json:"success"`
	Secret     string  `json:"secret"`
	Message    string  `json:"message"`
	Namespace  string  `json:"namespace"`
	Token      string  `json:"token"`
	Result     map[string]interface{} `json:"result"`
}

type CHLPayload struct {
	Allchannels            []baasapi.CHL
	Currchannels           []baasapi.CHL
}

type jsonChannelsResponse struct {
	Channels    []jsonChannelID    `json:"channels"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}
type jsoninstCCsResponse struct {
	Success    bool    `json:"success"`
	Message    string    `json:"message"`
	Result     []queryinstCC    `json:"result"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}
type queryinstCC struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Version    string    `json:"version"`
	
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}
type jsonCCsResponse struct {
	Success    bool    `json:"success"`
	Message    string    `json:"message"`
	Result     []queryCC    `json:"result"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}
type queryCC struct {
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	Path       string    `json:"path"`
	
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}

type jsonChannelID struct {
	Channel_id    string    `json:"channel_id"`
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
	//\"peers\": [\"peer0org2.demo-test.baas.com\",\"peer1org2.demo-test.baas.com\"],
	//\"peers\": [\"peer0org1.demo-test.baas.com\",\"peer0org2.demo-test.baas.com\"],
	//\"fcn\":\"move\",
	//\"args\":[\"a\",\"b\",\"10\"]
}

type enrollauthenticatePayload struct {
	Username   string  `json:"username"`
	OrgName    string  `json:"orgName"`
	Password    string  `json:"password"`
	//{"username": "Jim", "orgName": "Org1"}
	//"peers": ["peer0org2.demo-test.baas.com","peer1org2.demo-test.baas.com"]
}

// Baask8sJobRunner is used to run a Baask8sJob
type Baask8sJobRunner struct {
	schedule *baasapi.Schedule
	context  *Baask8sJobContext
}

// Baask8sJobContext represents the context of execution of a Baask8sJob
type Baask8sJobContext struct {
	baask8sService baasapi.Baask8sService
	cafilesmanager baasapi.CAFilesManager
	//baask8ster     baasapi.Baask8ster
}

// NewBaask8sJobContext returns a new context that can be used to execute a Baask8sJob
//func NewBaask8sJobContext(baask8sService baasapi.Baask8sService, baask8ster baasapi.Baask8ster) *Baask8sJobContext {
func NewBaask8sJobContext(baask8sService baasapi.Baask8sService, cafilesmanager baasapi.CAFilesManager) *Baask8sJobContext {
	return &Baask8sJobContext{
		baask8sService: baask8sService,
		cafilesmanager: cafilesmanager,
		//baask8ster:     baask8ster,
	}
}

// NewBaask8sJobRunner returns a new runner that can be scheduled
func NewBaask8sJobRunner(schedule *baasapi.Schedule, context *Baask8sJobContext) *Baask8sJobRunner {
	return &Baask8sJobRunner{
		schedule: schedule,
		context:  context,
	}
}

// GetSchedule returns the schedule associated to the runner
func (runner *Baask8sJobRunner) GetSchedule() *baasapi.Schedule {
	return runner.schedule
}

func baask8sSyncError(err error) bool {
	if err != nil {
		log.Printf("background job error (baask8s synchronization). Unable to synchronize baask8ss (err=%s)\n", err)
		return true
	}
	return false
}

// Run triggers the execution of the schedule.
// It will iterate through all the baask8ss available in the database to
// create a snapshot of each one of them.
// As a snapshot can be a long process, to avoid any concurrency issue we
// retrieve the latest version of the baask8s right after a snapshot.
func (runner *Baask8sJobRunner) Run() {
	go func() {

		log.Printf("background schedule running for baask8s.  (=%s)\n", "test...")

		baask8ss, err := runner.context.baask8sService.Baask8ss()
		if baask8sSyncError(err) {
			return
		}

		//log.Printf("(baask8s channel=%s) \n", baask8ss)

			for index, _ := range baask8ss {

				//for _, baask8sID := range runner.schedule.Baask8sJob.Baask8ss {

				//	log.Printf("backgroud schedule running for %s\n", baask8sID)
				//	baask8s, err := runner.context.baask8sService.Baask8s(baask8sID)
				//	if err != nil {
				//		log.Printf("scheduled job error (script execution). Unable to retrieve information about baask8s (id=%d) (err=%s)\n", baask8sID, err)
				//		return
				//	}

				//	log.Printf("backgroud schedule running for %s\n", baask8s.NetworkName)
				//	if (!runner.baask8s_sync_channels(*baask8s)) {
				//		continue
				//	} else {
				//	    if (!runner.baask8s_sync_instantiated_ccs(*baask8s)) {
				//			continue
				//		} else {
				//			if (!runner.baask8s_sync_installed_ccs(*baask8s)) {
				//				continue
				//			}
				//		}
				//	}
			
					//targets = append(targets, baask8s)
				//}

				log.Printf("backgroud schedule running for %s\n", baask8ss[index].ID)
				log.Printf("backgroud schedule running for %s\n", baask8ss[index].NetworkName)

				if (!runner.baask8s_sync_channels(baask8ss[index])) {
					continue
				} else {
				    if (!runner.baask8s_sync_instantiated_ccs(baask8ss[index])) {
						continue
					} else {
						if (!runner.baask8s_sync_installed_ccs(baask8ss[index])) {
							continue
						}
					}
				}

			}

	}()
}

func (runner *Baask8sJobRunner) baask8s_sync_channels(baask8s baasapi.Baask8s) (flag bool) {
	//var client http.Client
	client := &http.Client{Timeout: time.Second * 10}

	//var 
	//jsonValue, _ := json.Marshal({})
	//fmt.Println(bytes.NewBuffer(jsonValue))
	//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
	

	sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	for oindex, _ := range baask8s.MSPs {

		if (baask8s.MSPs[oindex].Role == 1){

		

		for oindex2, _ := range baask8s.MSPs[oindex].ORGs {

			orgname := baask8s.MSPs[oindex].ORGs[oindex2].ORGName
			fmt.Println("outside")
			fmt.Println(orgname)
			otoken := runner.baask8s_sync_get_tokenbyOrg(sdk_url, orgname)
			if (otoken == "") {
				fmt.Println("error accessing to sdk")
				return false
			}

			for oindex3, _ := range baask8s.MSPs[oindex].ORGs[oindex2].Peers {

			

		

	

	//sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	//req, err := http.NewRequest("GET", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn=query&args=%5B%22a%22%5D", nil)
	//req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn="+payload.Fcn+"&args="+payload.Args[0], nil)
	//req, err := http.NewRequest("GET", sdk_url+"/channels/?peer="+payload.Peer, bytes.NewBuffer(jsonValue))
	
	peername := strings.Split(strings.Split(baask8s.MSPs[oindex].ORGs[oindex2].Peers[oindex3],"@")[1],".")[0]+baask8s.NetworkID[0:13]+"-"+orgname
	req, err := http.NewRequest("GET", sdk_url+"/channels?peer="+peername, nil)
	req.Header.Add("Authorization" , "Bearer " + otoken)
	req.Header.Set("Content-Type", "application/json")

	if err == nil {
	resp, err := client.Do(req)

	if err == nil {

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 { // OK
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		//bodyString := string(bodyBytes)
		var responseObject jsonChannelsResponse
		json.Unmarshal(bodyBytes, &responseObject)

		fmt.Println(responseObject.Channels)

		for index, _ := range responseObject.Channels {

			fmt.Println("test")
			fmt.Println(responseObject.Channels[index].Channel_id)

			myorgs := []baasapi.MSPORG{}

			//baask8sID := handler.Baask8sService.GetNextIdentifier()
			//time.Now().Unix(),
			//var CreatedAt = time.Now().Format("2006-01-02 15:04:05")
			channel := baasapi.CHL{
				CHLName:          responseObject.Channels[index].Channel_id,
				CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
				ORGs:             myorgs,
			}
			flag := 0
			for index2, _ := range baask8s.CHLs {
				
				if (baask8s.CHLs[index2].CHLName == responseObject.Channels[index].Channel_id) {
					flag = 1

					fmt.Println("found match channel name")
					fmt.Println(baask8s.CHLs[index2].ORGs)
					
					flag2 := 0
					for index3, _ := range baask8s.CHLs[index2].ORGs {

						flag1 := 0
						log.Printf(baask8s.CHLs[index2].ORGs[index3].ORGName)

						
						if (baask8s.CHLs[index2].ORGs[index3].ORGName == baask8s.MSPs[oindex].ORGs[oindex2].ORGName) {
                            flag2 = 1 
							log.Printf("peer list")
							log.Printf(orgname)

							//baask8s.MSPs[oindex].ORGs[oindex2].ORGName

							//log.Printf(baask8s.CHLs[index2].ORGs[index3].Peers[0])
						for index4, _ := range baask8s.CHLs[index2].ORGs[index3].Peers {

							log.Printf(baask8s.CHLs[index2].ORGs[index3].Peers[index4])

						
						if (baask8s.CHLs[index2].ORGs[index3].Peers[index4] == peername) {
							flag1 = 1
							//baask8s.CHLs[index].ORGs[index2].Peers = append(baask8s.CHLs[index].ORGs[index2].Peers, payload.Peers...)
						} 
						}
						if (flag1 == 0){

							log.Printf("(need to insert baask8s peer into channel=%s) \n", peername)
							//myorg.ORGName = baask8s.MSPs[oindex].ORGs[oindex2].ORGName
							//myorg.Anchor = ""
							//channel.ORGs = append(channel.ORGs, myorg)
							baask8s.CHLs[index2].ORGs[index3].Peers = append(baask8s.CHLs[index2].ORGs[index3].Peers, peername)
							//baask8s.CHLs[index].ORGs = append(baask8s.CHLs[index].ORGs, myorgs)
							fmt.Println(baask8s.CHLs)
							err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
							if err != nil {
								fmt.Println(err)
								//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
							}
						}
						} 


						//}
						}
						if (flag2 == 0){ 
							myorg := baasapi.MSPORG{}
				
							//myorg.Peers = nil
							myorg.ORGName = orgname
							myorg.Anchor = ""
							myorg.Peers = append(myorg.Peers,peername)
	
							fmt.Println("adding another org"+orgname)
			
							baask8s.CHLs[index2].ORGs = append(baask8s.CHLs[index2].ORGs, myorg)
							err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
							if err != nil {
								fmt.Println(err)
								//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
							}
	
	
					}


					}

			}
			if flag == 0 {
				log.Printf("(need to insert baask8s channel=%s) \n", responseObject.Channels[index].Channel_id)



				myorg := baasapi.MSPORG{}
				
				//myorg.Peers = nil
				myorg.ORGName = orgname
				myorg.Anchor = ""
				myorg.Peers = append(myorg.Peers,peername)

				channel.ORGs = append(channel.ORGs, myorg)
			
			
				
				fmt.Println(baask8s.CHLs)
				//mychns := CHLPayload{}
				
				//mychns.Currchannels = baask8s.CHLs
				baask8s.CHLs = append(baask8s.CHLs, channel)
				//mychns.Allchannels = baask8s.CHLs

				err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
				if err != nil {
					fmt.Println(err)
					//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
				}

				runner.update_connection_profile(baask8s)
			}
			//log.Printf("(baask8s channel=%s) \n", baask8ss[index].CHLs[index2].CHLName)
			
		}

		//fmt.Println(responseObject)
		//fmt.Println(err)
		//return response.JSON(w, responseObject)



	}
	}
	}
			}
		}
		}


	}
	return true
	//return nil
}

func (runner *Baask8sJobRunner) baask8s_sync_instantiated_ccs(baask8s baasapi.Baask8s) (flag bool) {
	//var client http.Client
	client := &http.Client{Timeout: time.Second * 10}

	//var 
	fmt.Println("sync instantiated cc....")

	//jsonValue, _ := json.Marshal({})
	//fmt.Println(bytes.NewBuffer(jsonValue))
	//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
	

	sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	

	for cindex, _ := range baask8s.CHLs {

		var channelname = baask8s.CHLs[cindex].CHLName

	for oindex, _ := range baask8s.MSPs {

		if (baask8s.MSPs[oindex].Role == 1){

		

		for oindex2, _ := range baask8s.MSPs[oindex].ORGs {

			orgname := baask8s.MSPs[oindex].ORGs[oindex2].ORGName
			fmt.Println("outside cc")
			fmt.Println(orgname)
			otoken := runner.baask8s_sync_get_tokenbyOrg(sdk_url, orgname)
			if (otoken == "") {
				fmt.Println("error accessing to sdk")
				return false
			}

			for oindex3, _ := range baask8s.MSPs[oindex].ORGs[oindex2].Peers {

			

		

	

	//sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	//req, err := http.NewRequest("GET", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn=query&args=%5B%22a%22%5D", nil)
	//req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn="+payload.Fcn+"&args="+payload.Args[0], nil)
	//req, err := http.NewRequest("GET", sdk_url+"/channels/?peer="+payload.Peer, bytes.NewBuffer(jsonValue))
	
	peername := strings.Split(strings.Split(baask8s.MSPs[oindex].ORGs[oindex2].Peers[oindex3],"@")[1],".")[0]+baask8s.NetworkID[0:13]+"-"+orgname
	req, err := http.NewRequest("GET", sdk_url+"/channels/"+channelname+"/chaincodes?peer="+peername, nil)
	req.Header.Add("Authorization" , "Bearer " + otoken)
	req.Header.Set("Content-Type", "application/json")

	if err == nil {
	resp, err := client.Do(req)

	if err == nil {

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 { // OK
		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil { fmt.Println(err) }
		//bodyString := string(bodyBytes)
		var responseObject jsoninstCCsResponse
		json.Unmarshal(bodyBytes, &responseObject)

		fmt.Println(responseObject.Result)

		for index, _ := range responseObject.Result {

			//fmt.Println("test")
			fmt.Println(responseObject.Result[index].Name)
			fmt.Println(responseObject.Result[index].Version)
			fmt.Println(responseObject.Result[index].Path)

			//myorgs := []baasapi.MSPORG{}

			//ccs := []baasapi.CC{}

			//CC struct {
			//	//ID              CCID      `json:"id"`
			//	ID              int        `json:"id"`
			//	CCName          string     `json:"chaincodeName"`
			//	CHLName         string     `json:"chlName"`
			//	Version         string     `json:"chaincodeVersion"`
			//	EndorsementPolicyparm        interface{}     `json:"endorsementPolicyparm"`
			//	InstallORGs     []MSPORG   `json:"installORGs"`
			//	InstantiateORGs []MSPORG   `json:"instantiateORGs"`
			//	ChaincodeType   string     `json:"chaincodeType"`
			//}

			//baask8sID := handler.Baask8sService.GetNextIdentifier()
			cc := baasapi.CC{
				ID:          len(baask8s.CCs)+1,
				CCName:     responseObject.Result[index].Name,
				CHLName:    channelname,	
				Version:    responseObject.Result[index].Version,
				Path:        responseObject.Result[index].Path,
			}
			fmt.Println(cc)
			fmt.Println("above is cc")
			flag := 0
			if(responseObject.Result[index].Path == "") { continue }
			for index2, _ := range baask8s.CCs {
				//fmt.Println(baask8s.CCs[index2].CCName)
				//fmt.Println(baask8s.CCs[index2].Version)
				//fmt.Println(baask8s.CCs[index2].Path)
				//fmt.Println(orgname)
				//fmt.Println(peername)
				//fmt.Println(channelname)

				//if(baask8s.CCs[index2].CHLName != "") {
				
				if ((baask8s.CCs[index2].CCName == responseObject.Result[index].Name) && 
					((baask8s.CCs[index2].CHLName == channelname)) &&
				   //((baask8s.CCs[index2].CHLName == "") || (baask8s.CCs[index2].CHLName == channelname)) &&
				   (baask8s.CCs[index2].Version == responseObject.Result[index].Version) ) {
				   //(baask8s.CCs[index2].Path == responseObject.Result[index].Path)) {
					//baask8s.CCs[index2].CHLName = channelname

					//if (baask8s.CCs[index2].CHLName == "") || (baask8s.CCs[index2].CHLName == channelname) {
					flag = 1

					fmt.Println("found match cc name")
					flag2 := 0
					for index3, _ := range baask8s.CCs[index2].InstantiateORGs {
						//MSPORG struct {
						//	ORGName   string     `json:"ORGName"`
						//	Anchor    string     `json:"Anchor"`
						//	Peers     []string   `json:"Peers"`
						//}
						flag3 := 0
						if (baask8s.CCs[index2].InstantiateORGs[index3].ORGName == orgname) {
							flag2 = 1
							for index4, _ := range baask8s.CCs[index2].InstantiateORGs[index3].Peers {
								if (baask8s.CCs[index2].InstantiateORGs[index3].Peers[index4] == peername) {
									flag3 = 1
								}
							}
							if flag3 == 0 {
								//baask8s.CCs[index2].CHLName = channelname
								baask8s.CCs[index2].InstantiateORGs[index3].Peers = append(baask8s.CCs[index2].InstantiateORGs[index3].Peers,peername)
								fmt.Println(baask8s.CCs)
								fmt.Println("something to add")
								
								err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
								if err != nil {
									fmt.Println(err)
									//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
								}
								
							}
						}
					}
					if flag2 == 0 {
						myorg := baasapi.MSPORG{}
				
						//myorg.Peers = nil
						myorg.ORGName = orgname
						myorg.Anchor = ""
						myorg.Peers = append(myorg.Peers,peername)
						//baask8s.CCs[index2].CHLName = channelname
						baask8s.CCs[index2].InstantiateORGs = append(baask8s.CCs[index2].InstantiateORGs, myorg)
						fmt.Println(baask8s.CCs)
						fmt.Println("something to add flag2")
						
						err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
						if err != nil {
							fmt.Println(err)
							//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
						}
						
					}
					//fmt.Println(baask8s.CHLs[index2].ORGs)

					//}
					



				} else {
					
				}
				//}

			}
			if (flag == 0) {
				fmt.Println(baask8s.CCs)
				//mychns := CHLPayload{}
				
				//mychns.Currchannels = baask8s.CHLs
				baask8s.CCs = append(baask8s.CCs, cc)
				fmt.Println(baask8s.CCs)
				//mychns.Allchannels = baask8s.CHLs
				

				//fmt.Println(baask8s.CCs)
				
				err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
				if err != nil {
					fmt.Println(err)
					//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
				}
				

				//runner.update_connection_profile(baask8s)
				fmt.Println("cc not match")
			}

			//log.Printf("(baask8s channel=%s) \n", baask8ss[index].CHLs[index2].CHLName)
			
		}

		//fmt.Println(responseObject)
		//fmt.Println(err)
		//return response.JSON(w, responseObject)



	}
	}
	}
			}
		}
		}


	}
	}
	//return nil
	return true
}

func (runner *Baask8sJobRunner) baask8s_sync_installed_ccs(baask8s baasapi.Baask8s) (flag bool) {
	//var client http.Client
	client := &http.Client{Timeout: time.Second * 10}

	//var 
	fmt.Println("sync cc....")

	//jsonValue, _ := json.Marshal({})
	//fmt.Println(bytes.NewBuffer(jsonValue))
	//req, err := http.NewRequest("POST", "http://11.11.11.120:30500/channels/mychannel/peers", bytes.NewBuffer(jsonValue))
	

	sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	for oindex, _ := range baask8s.MSPs {

		if (baask8s.MSPs[oindex].Role == 1){

		

		for oindex2, _ := range baask8s.MSPs[oindex].ORGs {

			orgname := baask8s.MSPs[oindex].ORGs[oindex2].ORGName
			//fmt.Println("outside cc")
			//fmt.Println(orgname)
			otoken := runner.baask8s_sync_get_tokenbyOrg(sdk_url, orgname)
			if (otoken == "") {
				fmt.Println("error accessing to sdk")
				return false
			}

			for oindex3, _ := range baask8s.MSPs[oindex].ORGs[oindex2].Peers {

			

		

	

	//sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	//req, err := http.NewRequest("GET", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn=query&args=%5B%22a%22%5D", nil)
	//req, err := http.NewRequest("POST", sdk_url+"/channels/"+channelname+"/chaincodes/"+payload.ChaincodeName+"?peer="+payload.Peers[0]+"&fcn="+payload.Fcn+"&args="+payload.Args[0], nil)
	//req, err := http.NewRequest("GET", sdk_url+"/channels/?peer="+payload.Peer, bytes.NewBuffer(jsonValue))
	
	peername := strings.Split(strings.Split(baask8s.MSPs[oindex].ORGs[oindex2].Peers[oindex3],"@")[1],".")[0]+baask8s.NetworkID[0:13]+"-"+orgname
	req, err := http.NewRequest("GET", sdk_url+"/chaincodes?peer="+peername, nil)
	req.Header.Add("Authorization" , "Bearer " + otoken)
	req.Header.Set("Content-Type", "application/json")

	if err == nil {
	resp, err := client.Do(req)

	if err == nil {

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 { // OK
		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil { fmt.Println(err) }
		//bodyString := string(bodyBytes)
		var responseObject jsonCCsResponse
		json.Unmarshal(bodyBytes, &responseObject)

		fmt.Println(responseObject.Result)

		for index, _ := range responseObject.Result {

			//fmt.Println("test")
			//fmt.Println(responseObject.Result[index].Name)
			//fmt.Println(responseObject.Result[index].Version)
			//fmt.Println(responseObject.Result[index].Path)

			//myorgs := []baasapi.MSPORG{}

			//ccs := []baasapi.CC{}

			//CC struct {
			//	//ID              CCID      `json:"id"`
			//	ID              int        `json:"id"`
			//	CCName          string     `json:"chaincodeName"`
			//	CHLName         string     `json:"chlName"`
			//	Version         string     `json:"chaincodeVersion"`
			//	EndorsementPolicyparm        interface{}     `json:"endorsementPolicyparm"`
			//	InstallORGs     []MSPORG   `json:"installORGs"`
			//	InstantiateORGs []MSPORG   `json:"instantiateORGs"`
			//	ChaincodeType   string     `json:"chaincodeType"`
			//}

			//baask8sID := handler.Baask8sService.GetNextIdentifier()
			cc := baasapi.CC{
				ID:          len(baask8s.CCs)+1,
				CCName:     responseObject.Result[index].Name,
				CHLName:     "",
				Version:    responseObject.Result[index].Version,
				Path:        responseObject.Result[index].Path,
			}
			flag := 0
			for index2, _ := range baask8s.CCs {
				//fmt.Println(baask8s.CCs[index2].CCName)
				//fmt.Println(baask8s.CCs[index2].Version)
				//fmt.Println(baask8s.CCs[index2].Path)
				//fmt.Println(orgname)
				//fmt.Println(peername)


				
				if ((baask8s.CCs[index2].CCName == responseObject.Result[index].Name) && 
				   (baask8s.CCs[index2].Version == responseObject.Result[index].Version) ) {
				   //(baask8s.CCs[index2].Path == responseObject.Result[index].Path)) {
					flag = 1

					//fmt.Println("found match cc name")
					flag2 := 0
					for index3, _ := range baask8s.CCs[index2].InstallORGs {
						//MSPORG struct {
						//	ORGName   string     `json:"ORGName"`
						//	Anchor    string     `json:"Anchor"`
						//	Peers     []string   `json:"Peers"`
						//}
						flag3 := 0
						if (baask8s.CCs[index2].InstallORGs[index3].ORGName == orgname) {
							flag2 = 1
							for index4, _ := range baask8s.CCs[index2].InstallORGs[index3].Peers {
								if (baask8s.CCs[index2].InstallORGs[index3].Peers[index4] == peername) {
									flag3 = 1
								}
							}
							if flag3 == 0 {
								baask8s.CCs[index2].InstallORGs[index3].Peers = append(baask8s.CCs[index2].InstallORGs[index3].Peers,peername)
								err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
								if err != nil {
									fmt.Println(err)
									//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
								}
							}
						}
					}
					if flag2 == 0 {
						myorg := baasapi.MSPORG{}
				
						//myorg.Peers = nil
						myorg.ORGName = orgname
						myorg.Anchor = ""
						myorg.Peers = append(myorg.Peers,peername)
						baask8s.CCs[index2].InstallORGs = append(baask8s.CCs[index2].InstallORGs, myorg)
						err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
						if err != nil {
							fmt.Println(err)
							//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
						}
					}
					//fmt.Println(baask8s.CHLs[index2].ORGs)
					



				} else {
					
				}

			}
			if (flag == 0) {
				//fmt.Println(baask8s.CCs)
				//mychns := CHLPayload{}
				
				//mychns.Currchannels = baask8s.CHLs
				baask8s.CCs = append(baask8s.CCs, cc)
				//mychns.Allchannels = baask8s.CHLs

				err = runner.context.baask8sService.UpdateBaask8s(baask8s.ID, &baask8s)
				if err != nil {
					fmt.Println(err)
					//return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist baask8s changes inside the database", err}
				}

				//runner.update_connection_profile(baask8s)
				//fmt.Println("cc not match")
			}

			//log.Printf("(baask8s channel=%s) \n", baask8ss[index].CHLs[index2].CHLName)
			
		}

		//fmt.Println(responseObject)
		//fmt.Println(err)
		//return response.JSON(w, responseObject)



	}
	}
	}
			}
		}
		}


	}
	return true
	//return nil
}

func (runner *Baask8sJobRunner) update_connection_profile(baask8s baasapi.Baask8s) {

	mychns := CHLPayload{}
	
	mychns.Currchannels = baask8s.CHLs
	//baask8s.CHLs = append(baask8s.CHLs, channel)
	mychns.Allchannels = baask8s.CHLs

    //if len(mychns.Currchannels) == 0 {
	//	channel.CHLName = "default"
	//	mychns.Currchannels = append(mychns.Currchannels, channel)
	//}

	//fmt.Println(mychns)
	//log.Printf("(baask8s=%s) mychns type) \n", reflect.TypeOf(mychns))
    //response, err := http.Get("http://11.11.11.120:30500/users")
    //if err != nil {
    //    fmt.Printf("The HTTP request failed with error %s\n", err)
    //} else {
    //    data, _ := ioutil.ReadAll(response.Body)
    //    fmt.Println(string(data))
	//}
	//var data={};
    var jsonData []byte
	jsonData, err := json.Marshal(mychns)
    if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonData)
	
	//--extra-vars {"Allchannels":[{"Id":0,"CHLName":"tst4c","CreatedAt":1567312246,"ORGs":[]}]}

	var ansible_env = "mode=connection_apply env=" +baask8s.Namespace+ " deploy_type=k8s "
	var ansible_extra = string(jsonData)
	//ansible_env = ansible_env + " --extra-vars '{"allchannels":[{"Id":0,"CHLName":"tst4c","CreatedAt":1567312246,"ORGs":[]}]}'"
	//ansible_env = ansible_env + "\" --extra-vars '{\"allchannels\":[{\"Id\":0,\"CHLName\":\"tst4c\",\"CreatedAt\":1567312246,\"ORGs\":[]}]}'"
	var ansible_config = "/data/k8s/ansible/operatefabric.yml"
	err = runner.context.cafilesmanager.Deploy(baask8s.Owner, baask8s.Namespace, ansible_extra, ansible_env, ansible_config, true)
	if err != nil {
		fmt.Println(err)
	}

}




func (runner *Baask8sJobRunner) baask8s_sync_get_tokenbyOrg(sdk_url string, orgname string) (token string) {
	//sdk_url := "http://fabricsdk"+baask8s.NetworkID[0:13]+"."+baask8s.Namespace+":4000"
	//fmt.Printf(reflect.TypeOf(jsonValue))
	var epayload enrollauthenticatePayload

	epayload.Username = "admin"
	epayload.Password = "adminpw"
	epayload.OrgName = orgname

	jsonValue, _ := json.Marshal(epayload)
    jsonresponse, err := http.Post(sdk_url+"/enrollusers", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return ""
        //fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(jsonresponse.Body)
		//fmt.Println(string(data))
		var responseObject jsonResponse
		json.Unmarshal(data, &responseObject)

		//fmt.Println("first try.   join new channel..")
		fmt.Println(responseObject)

		//fmt.Println("2nd try...")



		//if err := json.NewDecoder(jsonresponse.Body).Decode(&responseObject); err != nil {
		//	log.Println(err)
		//}

		fmt.Println(responseObject.Token)

		//fmt.Println(responseObject)


		return responseObject.Token
	}

}