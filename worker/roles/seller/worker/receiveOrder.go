package worker

import (
	"log"

	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
)

// ReceiveOrderWorker receive order
func ReceiveOrderWorker(client worker.JobClient, job entities.Job) {
	log.Println("Start to receive order...")
	jobKey := job.GetKey()

}
