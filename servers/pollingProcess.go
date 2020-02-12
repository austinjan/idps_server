package servers

import (
	"io/ioutil"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	mongodb "github.com/austinjan/idps_server/servers/mongo"
	"fmt"
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
	ID      string `json:"_id" bson:"_id"`
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
	fmt.Printf("%s", body)
	var bsonData bson.M
	bson.Unmarshal(body,&bsonData)
	db.SaveTagPosition(bsonData) 
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
