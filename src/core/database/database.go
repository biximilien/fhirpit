package database

import (
	"dev.bonfhir/fhirpit/src/core"
)

var (
	client DatabaseClient
)

type DatabaseClient interface {
	GetSnomedDescription(conceptId string) ([]core.SnomedDescription, error)
	PutSnomedDescription(record core.SnomedDescription) error
	GetAdapter() interface{}
	SetAdapter(adapter interface{})
	Close()
}
