package ping

var DataPntr = &DataDet{}

type DataDet struct {
	MashlingCliRev      string
	MashlingCliLocalRev string
	MashlingCliVersion  string
	SchemaVersion       string
	AppVersion          string
	FlogolibRev         string
	MashlingRev         string
	AppDescription      string
}
