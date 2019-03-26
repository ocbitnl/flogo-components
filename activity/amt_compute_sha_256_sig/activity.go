package amt_compute_sha_256_sig

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strconv"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("amt_compute_sha_256_sig")

const (
	ivApiKey = "apiKey"
	ivSecret = "secret"

	ovXSignature = "xSignature"
)

// Sha256Activity is a stub for your Activity implementation
type Sha256Activity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &Sha256Activity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *Sha256Activity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *Sha256Activity) Eval(context activity.Context) (done bool, err error) {

	apiKey := context.GetInput(ivApiKey).(string)
	secret := context.GetInput(ivSecret).(string)

	s := apiKey + secret + strconv.FormatInt(time.Now().Unix(), 10)
	hash := sha256.New()
	hash.Write([]byte(s))

	xSignature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	log.Debugf("x-signature = %s", xSignature)

	context.SetOutput(ovXSignature, xSignature)

	return true, nil
}
