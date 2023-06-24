package code_system

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dev.bonfhir/fhirpit/src/core"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/parameters_go_proto"
)

type CodeSystemLookupOperation struct {
	Response CodeSystemLookupOperationResponse
	Records  []core.SnomedDescription
}

type CodeSystemLookupOperationResponse struct {
	ResourceType string                          `json:"resourceType,omitempty"`
	Parameters   *parameters_go_proto.Parameters `json:"parameter,omitempty"`
}

// CodeSystemHandler handles requests for the CodeSystem resource
func (codeSystemLookupOperation CodeSystemLookupOperation) CodeSystemLookupOperationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		descriptions := core.GetSnomedDescription(code)
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

	codeSystemLookupOperation.Response.ResourceType = "Parameters"

	// prepare response
	codeSystemLookupOperation.Response.Parameters = &parameters_go_proto.Parameters{
		Parameter: params,
	}

	data, err := json.Marshal(codeSystemLookupOperation.Response)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(data)
}
