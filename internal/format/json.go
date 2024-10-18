package format

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
)

func JSON(input *Data) string {
	inputRows := input.GetRows()
	inputCols := input.GetColumns()
	jsonData := make([]map[string]interface{}, 0, len(inputRows))

	for _, row := range inputRows {
		jsonRow := make(map[string]interface{}, len(inputCols))

		for colIndex, colValue := range row {
			colName := inputCols[colIndex].(string)
			jsonRow[colName] = colValue
		}
		jsonData = append(jsonData, jsonRow)
	}

	jsonBytes, err := json.MarshalIndent(jsonData, "", strings.Repeat(" ", 4))
	if err != nil {
		log.WithError(err).Error(
			"Failed to marshal JSON",
		)
		return ""
	}

	return string(jsonBytes)
}
