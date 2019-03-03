package main

import (
	"github.com/workflow-interoperability/samples/worker/roles/user/worker"
	"github.com/zeebe-io/zeebe/clients/go/zbc"
)

const brokerAddr = "127.0.0.1:26500"

func main() {
	client, err := zbc.NewZBClient(brokerAddr)
	if err != nil {
		panic(err)
	}
	placeOrderWorker := client.NewJobWorker().JobType("placeOrder").Handler(worker.PlaceOrderWorker).Open()
	defer placeOrderWorker.Close()

	placeOrderWorker.AwaitClose()
}
