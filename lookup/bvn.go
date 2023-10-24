package lookup

import (
	"github.com/prodbyola/mono"
	"github.com/prodbyola/mono/internal"
)

type BvnAvailableMethod struct {
	Method string `json:"method"`
	Hint   string `json:"hint"`
}

type BvnVerificationMethod struct {
}

// BvnInitiationResponse represents the response structure for BVN initiation.
// It contains information about the status of the initiation, a message, available methods,
// and a session ID.
type BvnInitiationResponse struct {
	Status    string               `json:"status"`
	Message   string               `json:"message"`
	Methods   []BvnAvailableMethod `json:"methods"`
	SessionID string               `json:"session_id"`
}

// parse is a method of the BvnInitiationResponse struct that populates the struct's fields
// by parsing data from a map[string]interface{}.
//
// Parameters:
//   - d: A map[string]interface{} containing data to populate the BvnInitiationResponse fields.
//     It is expected to have keys "status," "message," "data," and "session_id."
//
// The parse method extracts and processes the data to set the fields of the BvnInitiationResponse struct.
// It converts the "data" field to BvnAvailableMethod slice and handles data type assertions
// to ensure data integrity.
//
// Example usage:
//
//	response := BvnInitiationResponse{}
//	data := map[string]interface{}{"status": "success", "message": "Initiated", ...}
//	response.parse(data)
func (r *BvnInitiationResponse) parse(d map[string]interface{}) {
	message := d["message"].(string) // message should always be a string
	status := d["status"]
	data := d["data"]

	if data != nil {
		data := data.(map[string]interface{})
		sid := data["session_id"]

		m := data["methods"].([]interface{})
		var methods []BvnAvailableMethod

		// Let's convert a `[]interface{}` to []BvnVerificationMethod
		for i := 0; i < len(m); i++ {
			mtds := m[i].(map[string]interface{})
			nm := BvnAvailableMethod{
				Method: mtds["method"].(string),
				Hint:   mtds["hint"].(string),
			}

			methods = append(methods, nm)
		}

		if sid != nil {
			r.SessionID = sid.(string)
		}

		r.Methods = methods
	}

	r.Message = message
	if status != nil {
		r.Status = status.(string)
	}
}

func (r *BvnInitiationResponse) IsSuccessful() bool {
	return r.Status == "successful"
}

func (r *BvnInitiationResponse) NotSuccessful() bool {
	return !r.IsSuccessful()
}

// BvnLookUp represents a client for initiating and verifying BVN (Bank Verification Number) requests.
type BvnLookUp struct {
	apiKey string
}

// Creates a new BvnLookUp instance with the provided API key.
//
// Parameters:
//   - apiKey: A string representing the API key to be used for BVN lookup requests.
//
// Returns:
//   - BvnLookUp: A BvnLookUp instance with the specified API key.
func NewBvnLookUp(apiKey string) BvnLookUp {
	return BvnLookUp{apiKey}
}

// Initiates a BVN lookup request with the specified BVN.
//
// Parameters:
//   - bvn: A string representing the BVN to be initiated.
//
// Returns:
//   - BvnInitiationResponse: The response from the BVN initiation request.
//
// Check Step 1 on https://docs.mono.co/docs/bvn-lookup-integration-guide for more info.
func (l *BvnLookUp) Initiate(bvn string) BvnInitiationResponse {
	url := mono.BASE_URL + "lookup/bvn/initiate"
	var resp BvnInitiationResponse

	data := map[string]string{
		"bvn": bvn,
	}

	d, err := internal.MakeRequest[map[string]interface{}](internal.RequestConfig{
		Url:    url,
		Method: "POST",
		Data:   data,
		ApiKey: l.apiKey,
	})

	if err != nil {
		resp.Message = err.Error()
		return resp
	}

	if d != nil {
		resp.parse(d)
	}

	return resp
}

// Verifies a BVN with the specified verification method, session ID, and alternate phone number.
//
// Parameters:
//   - method: A BvnVerificationMethodType representing the verification method to be used.
//   - session_id: A string representing the session ID for the verification request.
//   - alt_phone: An interface{} that can be a string representing the alternate phone number when using 'alternate_phone' method.
//
// Returns:
//   - BvnInitiationResponse: The response from the BVN verification request.
//
// Check Step 2 on https://docs.mono.co/docs/bvn-lookup-integration-guide for more info.
func (l *BvnLookUp) Verify(method mono.BvnVerificationMethodType, session_id string, alt_phone interface{}) BvnInitiationResponse {
	url := mono.BASE_URL + "lookup/bvn/verify"
	var resp BvnInitiationResponse

	data := map[string]string{
		"method": method.String(),
	}

	if method == mono.AlternatePhoneMethod {
		if str, ok := alt_phone.(string); ok {
			data["phone_number"] = str
		} else {
			resp.Message = "You need to provided phone number when using 'alternate_phone' method"
			return resp
		}
	}

	d, err := internal.MakeRequest[map[string]interface{}](internal.RequestConfig{
		Url:    url,
		Method: "POST",
		Data:   data,
		ApiKey: l.apiKey,
		Headers: map[string]string{
			"x-session-id": session_id,
		},
	})

	if err != nil {
		resp.Message = err.Error()
		return resp
	}

	if d != nil {
		resp.parse(d)
	}

	if resp.IsSuccessful() {
		resp.SessionID = session_id
	}

	return resp
}
