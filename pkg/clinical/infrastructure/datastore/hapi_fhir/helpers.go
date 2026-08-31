package hapifhir

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/empower-clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/empower-clinical/pkg/clinical/domain"
)

// This is used to add field meta to any resource being created, which needs to be validated against a SGHI FHIR profile.
func AddRequiredFieldToPayload(payload map[string]interface{}, resourceType string) map[string]interface{} {
	baseURL := serverutils.MustGetEnvVar("HAPI_FHIR_BASE_URL")
	payload["resourceType"] = resourceType
	path := fmt.Sprintf("%s/StructureDefinition/%s", baseURL, strings.ToLower(resourceType))

	result, err := ProcessMetaData(payload, path)
	if err != nil {
		fmt.Printf("Error processing meta data: %v\n", err)
		return nil
	}

	return result
}

func ProcessMetaData(payload map[string]interface{}, path string) (map[string]interface{}, error) {
	meta, ok := payload["meta"].(map[string]interface{})
	if !ok {
		payload["meta"] = map[string]interface{}{
			"profile": []string{path},
		}
	}

	if ok {
		// add profile data to the meta, not to override existing data
		meta["profile"] = []string{path}

		tags, ok := meta["tag"].([]interface{})
		if !ok || len(tags) == 0 {
			return nil, fmt.Errorf("tag field not found or is not an array, or array is empty")
		}

		firstTag, ok := tags[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("first tag is not a map")
		}

		code, ok := firstTag["code"].(string)
		if !ok {
			return nil, fmt.Errorf("code field not found or is not a string")
		}

		if payload["identifier"] == nil {
			payload["identifier"] = createIdentifier(code)
		}
	}

	if _, found := payload["text"].(map[string]interface{}); !found {
		status := domain.NarrativeStatusEnumGenerated.String()
		payload["text"] = utils.NarrativeGenerator("Generated text", &status)
	}

	return payload, nil
}

func createIdentifier(code string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"use":    "official",
			"system": helpers.CodeSystem("default-identifier-codesystem"),
			"value":  uuid.New().String(),
			"type": map[string]interface{}{
				"coding": []map[string]interface{}{
					{
						"system":  helpers.CodeSystem("default-identifier-codesystem"),
						"code":    "default-id",
						"display": "Default Resource Identifier",
					},
				},
			},
			"assigner": map[string]interface{}{
				"reference": fmt.Sprintf("Organization/%s", code),
				"display":   "Organization",
			},
		},
	}
}
