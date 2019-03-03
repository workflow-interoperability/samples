package worker

import (
	"log"

	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
)

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}
