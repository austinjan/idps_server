package servers

import (
	"encoding/json"
	"fmt"
	mongodb "github.com/austinjan/idps_server/servers/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

func requestData() {
	db := mongodb.GetDB()
	resp, err := http.Get("http://localhost:3001/webapi/api/alarm")
	if err != nil {
		fmt.Println("Get error!!", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Get error!!", err.Error())
	}

	var bodyBson []map[string]interface{}

	if err := json.Unmarshal(body, &bodyBson); err != nil {
		log.Println("Request body can not parse!", err)
	}

	fmt.Println(bodyBson[0])

	db.SaveTagPosition(bodyBson[0])
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
