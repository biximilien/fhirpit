package resources

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
)

func PatientHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// create a new patient
	patient := &patient_go_proto.Patient{
		Active: &datatypes_go_proto.Boolean{
			Value: true,
		},
		Name: []*datatypes_go_proto.HumanName{
			{
				Family: &datatypes_go_proto.String{
					Value: "Doe",
				},
				Given: []*datatypes_go_proto.String{
					{
						Value: "John",
					},
				},
			},
		},
		Telecom: []*datatypes_go_proto.ContactPoint{
			{
				Value: &datatypes_go_proto.String{
					Value: "555-555-5555",
				},
			},
		},
		BirthDate: &datatypes_go_proto.Date{
			ValueUs: 1234567890,
		},
	}
	data, err := json.Marshal(patient)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(data)
}
