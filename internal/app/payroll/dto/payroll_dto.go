package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

type YearMonth string

func (y YearMonth) String() string {
	return string(y)
}

func (y YearMonth) TimeLayout() string {
	return "YYYY-MM"
}

func (y YearMonth) IsValid() bool {
	_, err := time.Parse(y.TimeLayout(), string(y))

	return err == nil
}

func (y *YearMonth) UnmarshalJSON(src []byte) error {
	if len(src) == 0 {
		return fmt.Errorf("empty year month")
	}

	var str string
	if err := json.Unmarshal(src, &str); err != nil {
		return fmt.Errorf("YearMonth json.Unmarshal string failed: %w", err)
	}

	_, err := time.Parse(y.TimeLayout(), str)
	if err != nil {
		return fmt.Errorf("YearMonth time.Parse failed: %w", err)
	}

	*y = YearMonth(str)
	return nil
}

func (y YearMonth) MarshalJSON() ([]byte, error) {
	if !y.IsValid() {
		return nil, fmt.Errorf("invalid YearMonth: %s", y)
	}

	return json.Marshal(y.String())
}

type GeneratePayrollTaskRequest struct {
	RequestID string    `json:"request_id"`
	YearMonth YearMonth `json:"year_month"`
}

type GeneratePayrollTaskResult struct {
	RequestID string `json:"request_id"`
}
