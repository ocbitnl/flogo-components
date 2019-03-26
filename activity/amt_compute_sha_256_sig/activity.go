package amt_compute_sha_256_sig

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strconv"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("amt_compute_sha_256_sig")

const (
	ivKey = "key"
	ivSecret = "secret"

	ovHex = "hex"
	ovBase64 = "base64"
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

	key := context.GetInput(ivKey).(string)
	secret := context.GetInput(ivSecret).(string)

	s := key + secret + strconv.FormatInt(time.Now().Unix(), 10)
	hash := sha256.New()
	hash.Write([]byte(s))

	hex := hex.EncodeToString(hash.Sum(nil))
	log.Debugf("hex = %s", hex)

	base64 := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	log.Debugf("base64 = %s", base64)

	context.SetOutput(ovHex, hex)
	context.SetOutput(ovBase64, base64)

	return true, nil
}
