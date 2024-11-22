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
	// Remove quotes around the time string
	s := string(b)
	s = s[1 : len(s)-1] // Remove the double quotes
	// Parse time
	parsedTime, err := time.Parse("2006-01-02T15:04Z07:00", s)
	if err != nil {
		return err
	}

	ct.Time = parsedTime.UTC()
	fmt.Println(ct.Time)
	
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// Format the time into the desired layout
	formattedTime := fmt.Sprintf(`"%s"`, ct.Format(time.RFC3339))
	return []byte(formattedTime), nil
}
var nilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) IsSet() bool {
    return ct.UnixNano() != nilTime
}

func (c CustomTime) Value() (driver.Value, error) {
    return driver.Value(c.Time.UTC()), nil
}

func (c *CustomTime) Scan(src interface{}) error {
    switch t := src.(type) {
    case time.Time:
        c.Time = t
        return nil
    default:
        return fmt.Errorf("column type not supported")
    }
}

type Birthdate struct {
	time.Time
}

const birthdateLayout = "2006-01-02"

func (bd *Birthdate) UnmarshalJSON(b []byte) error {
	// Remove quotes around the time string
	s := string(b)
	s = s[1 : len(s)-1] // Remove the double quotes
	// Parse time
	parsedTime, err := time.Parse(birthdateLayout, s)
	if err != nil {
		return err
	}
	bd.Time = parsedTime
	return nil
}

func (bd Birthdate) MarshalJSON() ([]byte, error) {
	// Format the time into the desired layout
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
        return fmt.Errorf("column type not supported")
    }
}
