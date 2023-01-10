package main

import (
	"strconv"
)

type IdName struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type HourMinute struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func str2HourMinute(hhmm string) *HourMinute {
	if hhmm == "" {
		return nil
	}
	var hour int
	var minute int
	if len(hhmm) == 3 {
		hour, _ = strconv.Atoi(hhmm[:1])
		minute, _ = strconv.Atoi(hhmm[1:])
	} else {
		hour, _ = strconv.Atoi(hhmm[:2])
		minute, _ = strconv.Atoi(hhmm[2:])
	}
	return &HourMinute{Hour: hour, Minute: minute}
}
