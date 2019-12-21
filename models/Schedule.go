package models

import "time"

type Schedule struct {
	ID        string
	Classroom Classroom
	Day       Day
	StartTime *time.Time
	EndTime   *time.Time
}
