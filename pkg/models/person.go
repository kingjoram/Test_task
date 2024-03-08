package models

type (
	Nationality struct {
		CountryId   string  `json:"country_id"`
		Probability float32 `json:"probability"`
	}
)
