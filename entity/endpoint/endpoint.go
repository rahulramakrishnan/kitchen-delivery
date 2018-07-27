package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kitchen-delivery/entity/exception"
	"github.com/pkg/errors"
)

// FormData holds generic http post request form data for our endpoints.
// We iterate over key, value pairs in order to hydrate a strictly typed request struct.
type FormData map[string][]string

// FieldsToExtract holds options to extract request parameters.
type FieldsToExtract struct {
	RequiredFields []string
	OptionalFields []string
}

// ExtractRequest converts list of required fields to the request object.
// If the request does not contain all of the required fields it returns an error.
// Normalization of parameters happens through converting the request to JSON
// and then deserializing the JSON string into a strictly typed request.
// For more complex request parameter management, do not use this function.
// ExtractRequest is meant for generic HTTP post requests.
func ExtractRequest(formData FormData, fieldsToExtract FieldsToExtract, request interface{}) error {
	// Verify form data contains all of the required fields.
	missingFields := getMissingFields(formData, fieldsToExtract.RequiredFields)
	if len(missingFields) > 0 {
		return errors.Wrap(
			exception.ErrInvalidInput, fmt.Sprintf("missing fields: %+v", missingFields),
		)
	}

	// Extract out required fields in order and generate JSON string to unmarshal into request object.
	var keyValuePairs []string // contains strings of key,value pairs from formData that we join on later
	for _, field := range fieldsToExtract.RequiredFields {
		valueStrs, _ := formData[field]                                       // formData => map[name:["Cheese Pizza"] temp:["hot"] "shelfLife":[300] "decayRate": 0.45]
		keyValuePairStr := fmt.Sprintf("\"%s\": \"%s\"", field, valueStrs[0]) // '"name": "Cheeze Pizza"'
		keyValuePairs = append(keyValuePairs, keyValuePairStr)                // ['"name": "Cheeze Pizza"', '"temp": "hot"', ... ]
	}

	// Extract optional fields if they exist. This is a best effort.
	for _, field := range fieldsToExtract.OptionalFields {
		valueStrs, ok := formData[field]
		if !ok {
			// If optional field does not exist we continue.
			continue
		}
		keyValuePairStr := fmt.Sprintf("\"%s\": \"%s\"", field, valueStrs[0])
		keyValuePairs = append(keyValuePairs, keyValuePairStr)
	}

	// ['"name": "Cheeze Pizza"', '"temp": "hot"', ... ] => '{"name: Cheeze Pizza", "temp: hot"}'
	jsonStr := fmt.Sprintf("{%s}", strings.Join(keyValuePairs, ", "))

	// Deserialize complete JSON string into address of request struct.
	err := json.Unmarshal([]byte(jsonStr), request)
	if err != nil {
		log.Printf("failed to unmarshal JSON | formData: %+v, JSON: %s, type: %T", formData, jsonStr, request)
		return err
	}

	return nil
}

// getMissingFields returns missing fields in an HTTP request given a list of required fields.
func getMissingFields(request map[string][]string, requiredFields []string) []string {
	var missingFields []string

	for _, field := range requiredFields {
		values, ok := request[field]
		if !ok || len(values) == 0 {
			missingFields = append(missingFields, field)
		}
	}

	return missingFields
}
