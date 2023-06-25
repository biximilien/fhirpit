package resources

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/code_system_go_proto"
)

// CodeSystemHandler handles requests for the CodeSystem resource
func CodeSystemHandler(w http.ResponseWriter, r *http.Request) {
	SetFHIRHeaders(w, r)

	code_system := &code_system_go_proto.CodeSystem{
		Name: &datatypes_go_proto.String{
			Value: "Test Code System",
		},
		Url: &datatypes_go_proto.Uri{
			Value: "http://test.com",
		},
	}
	data, err := json.Marshal(code_system)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(data)
}
