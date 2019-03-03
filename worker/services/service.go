package services

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/zeebe-io/zeebe/clients/go/entities"
	"github.com/zeebe-io/zeebe/clients/go/worker"
)

// FailJob fail a job
func FailJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
	client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send()
}

// BlockchainTransaction contact blockchain
func BlockchainTransaction(url, body string) error {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(string(resbody))
	}
	return nil
}
