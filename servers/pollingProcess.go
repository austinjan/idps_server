package servers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	mongodb "github.com/austinjan/idps_server/servers/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type pollingProcessor struct {
	ticker *time.Ticker
	stop   chan string
}

func newPollingPorcessor() *pollingProcessor {
	rv := &pollingProcessor{
		ticker: time.NewTicker(time.Second * 5),
		stop:   make(chan string),
	}
	go rv.run()
	return rv
}

type response struct {
	Key     string `json:"key" bson:"key"`
	Message string `json:"message" bson:"message"`
}

type tagInfoItem struct {
	Acceleration               interface{}   `json:"acceleration"`
	AccelerationTS             int           `json:"accelerationTS"`
	BatteryAlarm               string        `json:"batteryAlarm"`
	BatteryAlarmTS             int           `json:"batteryAlarmTS"`
	BatteryVoltage             string        `json:"batteryVoltage"`
	BatteryVoltageTS           int           `json:"batteryVoltageTS"`
	ButtonState                string        `json:"buttonState"`
	ButtonStateTS              int           `json:"buttonStateTS"`
	Color                      string        `json:"color"`
	ConfigStatus               string        `json:"configStatus"`
	ConfigStatusTS             interface{}   `json:"configStatusTS"`
	CoordinateSystemID         string        `json:"coordinateSystemId"`
	CoordinateSystemName       string        `json:"coordinateSystemName"`
	DeviceAddress              string        `json:"deviceAddress"`
	DeviceType                 string        `json:"deviceType"`
	Group                      string        `json:"group"`
	ID                         string        `json:"id"`
	IoStates                   []string      `json:"ioStates"`
	IoStatesTS                 int           `json:"ioStatesTS"`
	LastAreaID                 string        `json:"lastAreaId"`
	LastAreaName               string        `json:"lastAreaName"`
	LastAreaTS                 int           `json:"lastAreaTS"`
	LastButton2PressTS         interface{}   `json:"lastButton2PressTS"`
	LastButtonPressTS          interface{}   `json:"lastButtonPressTS"`
	LastPacketTS               int           `json:"lastPacketTS"`
	Name                       string        `json:"name"`
	Rssi                       int           `json:"rssi"`
	RssiCoordinateSystemID     string        `json:"rssiCoordinateSystemId"`
	RssiCoordinateSystemName   string        `json:"rssiCoordinateSystemName"`
	RssiLocator                string        `json:"rssiLocator"`
	RssiLocatorCoords          []float64     `json:"rssiLocatorCoords"`
	RssiTS                     int           `json:"rssiTS"`
	TagState                   string        `json:"tagState"`
	TagStateTS                 int           `json:"tagStateTS"`
	TagStateTransitionStatus   string        `json:"tagStateTransitionStatus"`
	TagStateTransitionStatusTS int           `json:"tagStateTransitionStatusTS"`
	TriggerCount               int           `json:"triggerCount"`
	TriggerCountTS             int           `json:"triggerCountTS"`
	TxPower                    int           `json:"txPower"`
	TxPowerTS                  int           `json:"txPowerTS"`
	TxRate                     int           `json:"txRate"`
	TxRateTS                   int           `json:"txRateTS"`
	Zones                      []interface{} `json:"zones"`
}

type tagInfo struct {
	Code             int           `json:"code"`
	Command          string        `json:"command"`
	Message          string        `json:"message"`
	OutputFormatID   string        `json:"outputFormatId"`
	OutputFormatName string        `json:"outputFormatName"`
	ResponseTS       int           `json:"responseTS"`
	Status           string        `json:"status"`
	Tags             []tagInfoItem `json:"tags"`
	Version          string        `json:"version"`
}

var hostURL = "http://210.201.89.44:8080"
var tagInfoAPI = "/qpe/getTagInfo?version=2"
var tagPositionAPI = "/qpe/getTagPosition?version=2"

func requestTagInfo() (bson.M, error) {
	resp, err := http.Get(hostURL + tagPositionAPI)
	if err != nil {
		fmt.Println("Get error!!", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Get error!!", err.Error())
		return nil, err
	}

	var bodyBson bson.M

	if err := json.Unmarshal(body, &bodyBson); err != nil {
		fmt.Printf("%s", body)
		fmt.Println("Request body can not parse!", err)
		return nil, err
	}
	return bodyBson, nil
}

func requestTagPosition() (bson.M, error) {
	resp, err := http.Get(hostURL + tagInfoAPI)
	if err != nil {
		fmt.Println("Get error!!", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Get error!!", err.Error())
		return nil, err
	}

	var bodyBson bson.M

	if err := json.Unmarshal(body, &bodyBson); err != nil {
		fmt.Println("Request body can not parse!", err)
		return nil, err
	}
	return bodyBson, nil
}

func requestData() {
	db := mongodb.GetDB()
	ret, err := requestTagInfo()
	if err == nil {
		fmt.Println("write tagInfo ")
		db.SaveTagInfo(ret)
	} else {
		fmt.Println("write tagInfo error", err)
	}

	ret, err = requestTagPosition()
	if err == nil {
		fmt.Println("write tagPosition ")
		db.SaveTagPosition(ret)
	} else {
		fmt.Println("write tagPosition error", err)
	}

}

func (p *pollingProcessor) run() {
	for {
		select {
		case <-p.ticker.C:
			requestData()
			//interval
		case v := <-p.stop:
			fmt.Println("Polling process receive stop ", v)
			return
		}
	}
}
