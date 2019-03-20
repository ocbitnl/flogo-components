package ping

import (
	"encoding/json"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("mashling-ping-activity")

const (
	ivCode = "code"
	ivData = "data"
)

// ReplyActivity is an Activity that is used to reply via the trigger
// inputs : {method,uri,params}
// outputs: {result}
type ReplyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new ReplyActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ReplyActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ReplyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ReplyActivity) Eval(ctx activity.Context) (done bool, err error) {

	replyCode := ctx.GetInput(ivCode).(int)
	replyData := ctx.GetInput(ivData)

	result := replyData.(map[string]interface{})

	if strings.Compare(result["response"].(string), "Detailed") == 0 {

		dataBytes := []byte(`{
			"response":"success",
			"mashlingCliRevision":"` + PingDataPntr.MashlingCliRev + `",
			"MashlingCliLocalRev":"` + PingDataPntr.MashlingCliLocalRev + `",
			"MashlingCliVersion":"` + PingDataPntr.MashlingCliVersion + `",
			"SchemaVersion":"` + PingDataPntr.SchemaVersion + `",
			"AppVersion":"` + PingDataPntr.AppVersion + `",
			"FlogolibRev":"` + PingDataPntr.FlogolibRev + `",
			"MashlingRev":"` + PingDataPntr.MashlingRev + `",
			"AppDescrption":"` + PingDataPntr.AppDescrption + `"
			}`)

		var replyJSON interface{}
		err = json.Unmarshal(dataBytes, &replyJSON)

		replyData = replyJSON
	}

	log.Debugf("Code :'%d', Data: '%+v'", replyCode, replyData)

	replyHandler := ctx.FlowDetails().ReplyHandler()

	//todo support replying with error

	if replyHandler != nil {

		replyHandler.Reply(replyCode, replyData, nil)
	}

	return true, nil
}

// Moved from /lib/util

var PingDataPntr = &PingDataDet{}

type PingDataDet struct {
	MashlingCliRev      string
	MashlingCliLocalRev string
	MashlingCliVersion  string
	SchemaVersion       string
	AppVersion          string
	FlogolibRev         string
	MashlingRev         string
	AppDescrption       string
}
