package worker

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/workflow-interoperability/samples/worker/services"
	"github.com/workflow-interoperability/samples/worker/types"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
)

// PlaceOrderWorker place order
func PlaceOrderWorker(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()
	log.Println("Start place order " + strconv.Itoa(int(jobKey)))

	payload, err := job.GetPayloadAsMap()
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	payload["name"] = "book"
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).PayloadFromMap(payload)
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
	newProcessInstance := types.Publish{
		ProcessID:              id,
		ProcessRelatedData:     string(aData),
		ApplicationRelatedData: []string{},
		SubscriberInformation: types.SubscriberInformation{
			Roles: []string{},
			ID:    "",
		},
	}
	body, err := json.Marshal(&newProcessInstance)
	if err != nil {
		log.Println(err)
		return
	}
	err = services.BlockchainTransaction("http://127.0.0.1:3000/api/Publish", string(body))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	payload["processID"] = id
	log.Println("Publish process success")

	// change blockchain instance state
	reqData1 := types.ChangeSubscriberInformation{
		ProcessID: payload["processID"].(string),
		SubscriberInformation: types.SubscriberInformation{
			Roles: []string{"seller"},
		},
	}
	reqData3 := types.ChangeCondition{
		ProcessID: payload["processID"].(string),
		Condition: "ordered",
	}

	jsonReqData1, err := json.Marshal(&reqData1)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	jsonReqData3, err := json.Marshal(&reqData3)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}

	err = services.BlockchainTransaction("http://127.0.0.1:3000/api/ChangeSubscriberInformation", string(jsonReqData1))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	log.Println("Change subscriber success")
	err = services.BlockchainTransaction("http://127.0.0.1:3000/api/ChangeCondition", string(jsonReqData3))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	log.Println("Change condition to ordered success")

	log.Println("Place order " + strconv.Itoa(int(jobKey)) + "success")
	request.Send()
}
