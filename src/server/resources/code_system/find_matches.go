package code_system

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"dev.bonfhir/fhirpit/src/core"
	"dev.bonfhir/fhirpit/src/core/database"
	"dev.bonfhir/fhirpit/src/server/resources"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/parameters_go_proto"
)

type CodeSystemFindMatchesOperation struct {
	Response CodeSystemFindMatchesOperationResponse
	Records  []core.SnomedDescription
	client   database.DatabaseClient
}

type CodeSystemFindMatchesOperationResponse struct {
	ResourceType string                          `json:"resourceType,omitempty"`
	Parameters   *parameters_go_proto.Parameters `json:"parameter,omitempty"`
}

func NewCodeSystemFindMatchesOperation(client database.DatabaseClient) *CodeSystemFindMatchesOperation {
	return &CodeSystemFindMatchesOperation{
		Response: CodeSystemFindMatchesOperationResponse{},
		client:   client,
	}
}

// CodeSystemHandler handles requests for the CodeSystem resource
func (op CodeSystemFindMatchesOperation) CodeSystemFindMatchesOperationHandler(w http.ResponseWriter, r *http.Request) {
	resources.SetFHIRHeaders(w, r)

	// code, err := strconv.Atoi(r.URL.Query().Get("code"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	code := r.URL.Query().Get("code")
	system := r.URL.Query().Get("system")
	// version := r.URL.Query().Get("version")
	// property := r.URL.Query().Get("property")

	params := []*parameters_go_proto.Parameters_Parameter{}

	if system == "http://snomed.info/sct" {
		descriptions, err := op.client.GetSnomedDescription(code)
		if err != nil {
			log.Println(err)
		}
		for _, description := range descriptions {
			params = append(params, &parameters_go_proto.Parameters_Parameter{

				Name: &datatypes_go_proto.String{
					Value: "http://snomed.info/sct",
				},

				Part: []*parameters_go_proto.Parameters_Parameter{
					{
						Name: &datatypes_go_proto.String{
							Value: description.Term,
						},
						Value: &parameters_go_proto.Parameters_Parameter_ValueX{},
					},
				},
			})
		}
	}

	op.Response.ResourceType = "Parameters"

	// prepare response
	op.Response.Parameters = &parameters_go_proto.Parameters{
		Parameter: params,
	}

	data, err := json.Marshal(op.Response)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(data)
}
