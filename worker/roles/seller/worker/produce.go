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

// Produce produce products
func Produce(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()
	log.Println("Start produce for order " + strconv.Itoa(int(jobKey)))
	payload, err := job.GetPayloadAsMap()
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).PayloadFromMap(payload)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}

	log.Println("Produce ", payload["name"])
	// use sleep to mock produce
	time.Sleep(5 * time.Second)

	id, has := payload["processID"].(string)
	if !has {
		log.Println("can not get process id from payload")
		services.FailJob(client, job)
		return
	}
	reqData3 := types.ChangeCondition{
		ProcessID: id,
		Condition: "finished",
	}

	jsonReqData3, err := json.Marshal(&reqData3)
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}

	err = services.BlockchainTransaction("http://127.0.0.1:3001/api/ChangeCondition", string(jsonReqData3))
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}
	log.Println("Produce for order " + strconv.Itoa(int(jobKey)) + " finished")
	request.Send()
}
