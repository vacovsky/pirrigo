package data

import (
	"encoding/json"

	"github.com/vacovsky/pirrigo/src/logging"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func jsonifySQLResults(input *gorm.DB) []string {
	var result = []string{}
	r, err := json.Marshal(input.Value)
	if err != nil {
		logging.Service().LogError("Problem parsing SQL results.",
			zap.String("error", err.Error()))
	}
	result = append(result, string(r))
	return result
}
