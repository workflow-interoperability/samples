package worker

import (
	"encoding/json"
	"log"

	"github.com/workflow-interoperability/samples/worker/services"
	"github.com/workflow-interoperability/samples/worker/types"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
)

// PlaceOrderWorker place order
func PlaceOrderWorker(client worker.JobClient, job entities.Job) {
	log.Println("Start to place order...")
	jobKey := job.GetKey()

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

	// change blockchain instance state
	reqData1 := types.ChangeSubscriberInformation{
		ProcessID: payload["processID"].(string),
		SubscriberInformation: types.SubscriberInformation{
			Roles: []string{"seller"},
		},
	}
	data, err := json.Marshal(&payload)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	reqData2 := types.ChangeProcessData{
		IsApplicationRelatedDataChanged: true,
		ProcessRelatedData:              string(data),
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
	jsonReqData2, err := json.Marshal(&reqData2)
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
	err = services.BlockchainTransaction("http://127.0.0.1:3000/api/CangeProcessData", string(jsonReqData2))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	err = services.BlockchainTransaction("http://127.0.0.1:3000/api/ChangeCondition", string(jsonReqData3))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}

	log.Println("Place holder success!!!")
	request.Send()
}
