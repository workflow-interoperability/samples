package worker

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/workflow-interoperability/samples/worker/services"
	"github.com/workflow-interoperability/samples/worker/types"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
)

// PlaceOrderWorker place order
func PlaceOrderWorker(client worker.JobClient, job entities.Job) {
	processID := "user"
	IESMID := "1"
	jobKey := job.GetKey()
	log.Println("Start place order " + strconv.Itoa(int(jobKey)))

	payload, err := job.GetVariablesAsMap()
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	payload["name"] = "book"
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(payload)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}

	// create blockchain instance
	id := services.GenerateXID()
	aData, err := json.Marshal(&payload)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	newIM := types.IM{
		ID: id,
		Payload: types.Payload{
			ApplicationData: types.ApplicationData{
				URL: string(aData),
			},
			WorkflowRelevantData: types.WorkflowRelevantData{
				From: types.FromToData{
					ProcessID:         processID,
					ProcessInstanceID: payload["processInstanceID"].(string),
					IESMID:            IESMID,
				},
				To: types.FromToData{
					ProcessID:         "seller",
					ProcessInstanceID: "-1",
					IESMID:            "1",
				},
			},
		},
		SubscriberInformation: types.SubscriberInformation{
			Roles: []string{},
			ID:    "seller",
		},
	}
	pim := types.PublishIM{newIM}
	body, err := json.Marshal(&pim)
	if err != nil {
		log.Println(err)
		return
	}
	err = services.BlockchainTransaction("http://127.0.0.1:3000/api/PublishIM", string(body))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	payload["processID"] = id
	log.Println("Publish IM success")

	// waiting for PIIS from receiver
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:3001", Path: ""}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	for {
		finished := false
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
			services.FailJob(client, job)
			return
		}
		switch structMsg["$class"].(string) {
		case "org.sysu.wf.PIISCreatedEvent":
			if ok, err := publishPIIS(structMsg["id"].(string), newIM, c); err != nil {
				services.FailJob(client, job)
				return
			} else if ok {
				finished = true
				break
			}
		default:
			continue
		}
		if finished {
			log.Println("publish piis success")
			break
		}
	}
	request.Send()
}
