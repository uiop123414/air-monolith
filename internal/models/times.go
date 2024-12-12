package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

const timeLayout = "2006-01-02T15:04Z07:00"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	parsedTime, err := time.Parse(timeLayout, s)
	if err != nil {
		return err
	}

	ct.Time = parsedTime

	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formattedTime := fmt.Sprintf(`"%s"`, ct.Format(time.RFC3339))
	return []byte(formattedTime), nil
}

var nilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

func (c CustomTime) Value() (driver.Value, error) {
	return driver.Value(c.Time.Format(time.RFC3339)), nil
}

func (ct *CustomTime) GetTimezone() int {
	_, offset := ct.Zone()
	return -offset / 3600
}

func (c *CustomTime) Scan(src interface{}) error {
	switch t := src.(type) {
	case time.Time:
		c.Time = t
		return nil
	default:
		return ErrColumnNotSupported
	}
}

type Birthdate struct {
	time.Time
}

const birthdateLayout = "2006-01-02"

func (bd *Birthdate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	parsedTime, err := time.Parse(birthdateLayout, s)
	if err != nil {
		return err
	}
	bd.Time = parsedTime
	return nil
}

func (bd Birthdate) MarshalJSON() ([]byte, error) {
	formattedTime := fmt.Sprintf(`"%s"`, bd.Format(birthdateLayout))
	return []byte(formattedTime), nil
}

func (bd *Birthdate) IsSet() bool {
	return bd.UnixNano() != nilTime
}

func (bd Birthdate) Value() (driver.Value, error) {
	return driver.Value(bd.Time.Format(birthdateLayout)), nil
}

func (bd *Birthdate) Scan(src interface{}) error {
	switch t := src.(type) {
	case time.Time:
		bd.Time = t
		return nil
	default:
		return ErrColumnNotSupported
	}
}
