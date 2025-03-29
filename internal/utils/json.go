package helpers

import (
	"encoding/json"
	"fmt"
)

func PrettyJSON(req map[string]interface{}) string {
	jsonStr, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(jsonStr)
}
