package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// makeRequest is a generic function for making HTTP requests with the provided parameters.
//
// This function can be used to send HTTP requests with various request data types (R) to a specified URL
// using the given HTTP method, API key, and request data. The response is unmarshaled into the provided
// response data type (R) and returned along with any potential error.
//
// Parameters:
//   - url: A string representing the URL to send the HTTP request to.
//   - method: A string representing the HTTP method (e.g., "GET", "POST") for the request.
//   - data: The request data to be included in the HTTP request body. This can be of any type (R) that can be
//     marshaled into JSON. Pass nil if no data is required.
//   - apiKey: A string representing the API key to be included in the request headers.
//
// Returns:
// - R: The response data type. The HTTP response is unmarshaled into this type.
// - error: An error, if any, encountered during the HTTP request or response processing.
//
// Example usage:
//
//	response, err := makeRequest<MyResponseType>("https://example.com/api", "POST", requestData, "API_KEY")
func MakeRequest[R any](url string, method string, data interface{}, apiKey string) (R, error) {
	var resp R
	client := http.Client{}
	var buff *bytes.Buffer

	if data != nil {
		req_data, err := json.Marshal(data)
		if err != nil {
			return resp, err
		}

		buff = bytes.NewBuffer(req_data)
	}

	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return resp, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("mono-sec-key", apiKey)

	res, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	json.NewDecoder(res.Body).Decode(&resp)
	return resp, nil
}
