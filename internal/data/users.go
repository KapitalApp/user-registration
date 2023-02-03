package data

import (
	"github.com/google/uuid"
	"time"
)

type MetaField struct {
	Key       string
	Namespace string
	Value     string
	Type      string
}

type RangeNumber struct {
	UpLimit   int64
	DownLimit int64
}

type User struct {
	ID                     uuid.UUID
	Email                  string
	FirstName              string
	LastName               string
	ProvinceCode           string
	CountryCode            string
	AdministrativeDivision string
	AgeRange               RangeNumber
	FamilyNumberRange      RangeNumber
	CreatedAt              time.Time
	Meta                   []MetaField
}
