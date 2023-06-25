package resources

import "net/http"

const (
	JSON        string = "application/json"
	JSON_FHIR   string = "application/json+fhir"
	FHIR_JSON   string = "application/fhir+json"
	TEXT_JSON   string = "text/json"
	XML         string = "application/xml"
	XML_FHIR    string = "application/xml+fhir"
	FHIR_XML    string = "application/fhir+xml"
	TEXT_XML    string = "text/xml"
	Turtle      string = "application/turtle"
	Turtle_FHIR string = "application/turtle+fhir"
	FHIR_Turtle string = "application/fhir+turtle"
	TEXT_Turtle string = "text/turtle"
)

func SetFHIRHeaders(w http.ResponseWriter, r *http.Request) {

	// 3.2.0.1.10 Content Types and encodings

	// The formal MIME-type for FHIR resources is application/fhir+xml or application/fhir+json. The correct mime type SHALL be used by clients and servers:

	// XML: application/fhir+xml
	// JSON: application/fhir+json
	// RDF: application/fhir+turtle (only the Turtle format is supported)

	// Implementation Notes:

	//    - The content type application/x-www-form-urlencoded (Specification icon) is also accepted for posting search requests.
	//    - If a client provides a generic mime type in the Accept header (application/xml, text/json, or application/json), the server SHOULD respond with the requested mime type, using the XML or JSON formats described in this specification as the best representation for the named mime type (except for binary - see the note on the Binary resource).
	//    - Note: between FHIR DSTU2 and STU3, the correct mime type was changed from application/xml+fhir and application/json+fhir to application/fhir+xml and application/fhir+json. Servers MAY also support the older mime types, and are encouraged to do so to smooth the transition process.
	//    - 406 Not Acceptable is the appropriate response when the Accept header requests a format that the server does not support, and 415 Unsupported Media Type when the client posts a format that is not supported to the server.

	// UTF-8 encoding SHALL be used for FHIR instances. This MAY be specified as a MIME type parameter, but is not required.

	switch content_type := r.Header.Get("Content-Type"); content_type {
	case JSON:
		w.Header().Set("Content-Type", JSON)
	case JSON_FHIR:
		w.Header().Set("Content-Type", JSON_FHIR)
	case FHIR_JSON:
		w.Header().Set("Content-Type", FHIR_JSON)
	case TEXT_JSON:
		w.Header().Set("Content-Type", TEXT_JSON)
	case XML:
		w.Header().Set("Content-Type", XML)
	case XML_FHIR:
		w.Header().Set("Content-Type", XML_FHIR)
	case FHIR_XML:
		w.Header().Set("Content-Type", FHIR_XML)
	case TEXT_XML:
		w.Header().Set("Content-Type", TEXT_XML)
	case Turtle:
		w.Header().Set("Content-Type", Turtle)
	case Turtle_FHIR:
		w.Header().Set("Content-Type", Turtle_FHIR)
	case FHIR_Turtle:
		w.Header().Set("Content-Type", FHIR_Turtle)
	case TEXT_Turtle:
		w.Header().Set("Content-Type", TEXT_Turtle)
	default:
		// default to FHIR JSON format
		w.Header().Set("Content-Type", FHIR_JSON)
	}
}
