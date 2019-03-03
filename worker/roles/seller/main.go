package main

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/workflow-interoperability/samples/worker/roles/seller/worker"
	"github.com/workflow-interoperability/samples/worker/services"
	"github.com/workflow-interoperability/samples/worker/types"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
)

const brokerAddr = "127.0.0.1:26500"

func main() {
	client, err := zbc.NewZBClient(brokerAddr)
	if err != nil {
		panic(err)
	}

	stopChan := make(chan bool, 0)

	// define worker
	receiveOrderWorker := client.NewJobWorker().JobType("placeOrder").Handler(worker.ReceiveOrderWorker).Open()
	defer receiveOrderWorker.Close()
	go receiveOrderWorker.AwaitClose()

	produce := client.NewJobWorker().JobType("placeOrder").Handler(worker.Produce).Open()
	defer produce.Close()
	go produce.AwaitClose()

	// listen to blockchain event
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:3000", Path: ""}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			// check message type and handle
			var structMsg map[string]interface{}
			err = json.Unmarshal(msg, &structMsg)
			if err != nil {
				log.Println(err)
				return
			}
			switch structMsg["$class"].(string) {
			case "org.sysu.wf.ConditionChangedEvent":
				createSellerWorkflowInstance(structMsg["processID"].(string), client)
			}
		}
	}()

	<-stopChan
}

func createSellerWorkflowInstance(processID string, client zbc.ZBClient) {
	// get process data
	processData, err := services.GetProcessInstance("http://127.0.0.1:3000/api/ProcessInstance/" + processID)
	if err != nil {
		log.Println(err)
		return
	}
	// publish blockchain asset
	id := services.GenerateXID()
	var data map[string]interface{}
	if processData.ProcessRelatedData != "" {
		err = json.Unmarshal([]byte(processData.ProcessRelatedData), &data)
		if err != nil {
			log.Println(err)
			return
		}
		delete(data, "processID")
	}
	aData, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
		return
	}
	newProcessInstance := types.Publish{
		ProcessID:          id,
		ProcessRelatedData: string(aData),
	}
	body, err := json.Marshal(&newProcessInstance)
	if err != nil {
		log.Println(err)
		return
	}
	err = services.BlockchainTransaction("http://127.0.0.1:3001/api/Publish", string(body))
	if err != nil {
		log.Println(err)
		return
	}

	// add workflow instance
	request, err := client.NewCreateInstanceCommand().BPMNProcessId("receiveOrder").LatestVersion().PayloadFromMap(data)
	if err != nil {
		log.Println(err)
		return
	}
	msg, err := request.Send()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(msg.String())
}
