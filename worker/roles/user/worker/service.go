package worker

import (
	"encoding/json"
	"log"

	"github.com/workflow-interoperability/samples/worker/types"

	"github.com/gorilla/websocket"
	"github.com/workflow-interoperability/samples/worker/services"
)

func publishPIIS(piisid, processid, processInstanceid, iermid string, IM types.IM, conn *websocket.Conn) (bool, error) {
	// get piis
	processData, err := services.GetPIIS("http://127.0.0.1:3000/api/PIIS/" + piisid)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if !(processData.To.ProcessID == processid && processData.To.ProcessInstanceID == processInstanceid && processData.To.IESMID == iermid) {
		return false, nil
	}
	// create piis
	id := services.GenerateXID()
	newPIIS := types.PIIS{
		ID:    id,
		From:  IM.Payload.WorkflowRelevantData.From,
		To:    IM.Payload.WorkflowRelevantData.To,
		IMID:  IM.ID,
		Owner: "user",
		SubscriberInformation: types.SubscriberInformation{
			Roles: []string{},
			ID:    "seller",
		},
	}
	pPIIS := types.PublishPIIS{newPIIS}
	body, err := json.Marshal(&pPIIS)
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = services.BlockchainTransaction("http://127.0.0.1:3000/api/PublishPIIS", string(body))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}
