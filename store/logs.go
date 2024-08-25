package store

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/zarbanio/market-maker-keeper/pkg/utils"
)

func (p postgres) CreateLog(ctx context.Context, b []byte) (int, error) {
	var entry map[string]interface{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		return 0, err
	}

	level := utils.ConvertToString(entry["level"])
	if level == "" {
		return 0, errors.New("log level didn't found")
	}

	if level == "debug" {
		return len(b), nil
	}

	message := utils.ConvertToString(entry["message"])
	if message == "" {
		return 0, errors.New("log message didn't found")
	}

	cycleId, ok := entry["cycleId"]
	if !ok {
		cycleId = 0
	}

	fieldsJSON, err := json.Marshal(entry)
	if err != nil {
		return 0, err
	}
	stmt := `INSERT INTO logs (cycle_id, level, message, fields) VALUES ($1, $2, $3, $4)`

	_, err = p.conn.Exec(ctx, stmt, cycleId, level, message, string(fieldsJSON))
	if err != nil {
		return 0, err
	}

	return len(b), nil
}
