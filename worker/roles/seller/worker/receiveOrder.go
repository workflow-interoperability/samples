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

// ReceiveOrderWorker receive order
func ReceiveOrderWorker(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()
	log.Println("Start receive order" + strconv.Itoa(int(jobKey)))
	payload, err := job.GetPayloadAsMap()
	if err != nil {
		log.Println(err)
		services.FailJob(client, job)
		return
	}

	reqData3 := types.ChangeCondition{
		ProcessID: payload["processID"].(string),
		Condition: "ordered",
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
}
