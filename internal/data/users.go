package data

import (
	"github.com/google/uuid"
	"time"
)

type MetaField struct {
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
	Value     string `json:"value"`
	Type      string `json:"type"`
}

type RangeNumber struct {
	UpLimit   int64 `json:"up_limit"`
	DownLimit int64 `json:"down_limit"`
}

type User struct {
	ID                     uuid.UUID   `json:"id"`
	Email                  string      `json:"email"`
	FirstName              string      `json:"first_name"`
	LastName               string      `json:"last_name"`
	ProvinceCode           string      `json:"province_code"`
	CountryCode            string      `json:"country_code"`
	AdministrativeDivision string      `json:"administrative_division"`
	AgeRange               RangeNumber `json:"age_range"`
	FamilyNumber           int64       `json:"family_number"`
	CreatedAt              time.Time   `json:"created_at"`
	Meta                   []MetaField `json:"meta"`
}
