package lookup

import (
	"errors"

	"github.com/prodbyola/mono"
	"github.com/prodbyola/mono/internal"
)

type BvnVerificationMethod struct {
	Method string `json:"method"`
	Hint   string `json:"hint"`
}

type BvnInitiationResponse struct {
	Status    string                  `json:"status"`
	Message   string                  `json:"message"`
	Methods   []BvnVerificationMethod `json:"methods"`
	SessionID string                  `json:"session_id"`
}

func (r *BvnInitiationResponse) parse(d map[string]interface{}) {
	message := d["message"].(string) // message should always be a string
	status := d["status"]
	data := d["data"]

	if data != nil {
		data := data.(map[string]interface{})
		sid := data["session_id"]

		m := data["methods"].([]interface{})
		var methods []BvnVerificationMethod

		// Let's convert a `[]interface{}` to []BvnVerificationMethod
		for i := 0; i < len(m); i++ {
			mtds := m[i].(map[string]interface{})
			nm := BvnVerificationMethod{
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

type BvnLookUp struct {
	apiKey string
}

func NewBvnLookUp(apiKey string) BvnLookUp {
	return BvnLookUp{apiKey}
}

func (l *BvnLookUp) InitiateBvn(bvn string) (BvnInitiationResponse, error) {
	url := mono.BASE_URL + "lookup/bvn/initiate"
	var resp BvnInitiationResponse

	data := map[string]string{
		"bvn": bvn,
	}

	d, err := internal.MakeRequest[map[string]interface{}](url, "POST", data, l.apiKey)
	if err != nil {
		return resp, err
	}

	if d != nil {
		resp.parse(d)

		if !resp.IsSuccessful() {
			return resp, errors.New(resp.Message)
		}
	}

	return resp, nil
}
