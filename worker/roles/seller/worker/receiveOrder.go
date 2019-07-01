package worker

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/workflow-interoperability/samples/worker/services"
	"github.com/workflow-interoperability/samples/worker/types"
	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
)

// ReceiveOrderWorker receive order
func ReceiveOrderWorker(client worker.JobClient, job entities.Job) {
	processID := "seller"
	iesmid := "1"
	jobKey := job.GetKey()
	log.Println("Start receive order " + strconv.Itoa(int(jobKey)))
	payload, err := job.GetVariablesAsMap()
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(payload)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}

	// use sleep to mock produce
	time.Sleep(5 * time.Second)
	id := services.GenerateXID()
	newPIIS := types.PIIS{
		ID: id,
		From: types.FromToData{
			ProcessID:         processID,
			ProcessInstanceID: payload["processInstanceID"].(string),
			IESMID:            iesmid,
		},
		To: types.FromToData{
			ProcessID:         "user",
			ProcessInstanceID: payload["fromProcessInstanceID"].(string),
			IESMID:            "1",
		},
		IMID: "-1",
		SubscriberInformation: types.SubscriberInformation{
			Roles: []string{},
			ID:    "user",
		},
	}
	pPIIS := types.PublishPIIS{newPIIS}
	body, err := json.Marshal(&pPIIS)
	if err != nil {
		log.Println(err)
		return
	}
	err = services.BlockchainTransaction("http://127.0.0.1:3001/api/PublishPIIS", string(body))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	log.Println("Publish PIIS success")
	request.Send()
}
