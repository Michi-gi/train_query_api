package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi"
)

type StopStation struct {
	Code      string `json:"stationCode"`
	Name      string `json:"stationName"`
	AreaCode  string `json:"areaCode"`
	PrefCode  string `json:"prefCode"`
	Arrival   string `json:"arrivalTime"`
	Departure string `json:"departureTime"`
}

type Train struct {
	Props struct {
		PageProps struct {
			TimetableStationTrainResult struct {
				Timetable struct {
					TrainId  string        `json:"trainId"`
					Name     string        `json:"displayName"`
					Stations []StopStation `json:"stopStation"`
				} `json:"timetable"`
			} `json:"timetableStationTrainResult"`
		} `json:"PageProps"`
	} `json:"props"`
}

type TrainReturn struct {
	Id       string              `json:"id"`
	Name     string              `json:"name"`
	LineName string              `json:"lineName"`
	Stations []StopStationReturn `json:"stations"`
}

type StopStationReturn struct {
	Code      string      `json:"stationCode"`
	Name      string      `json:"stationName"`
	AreaCode  string      `json:"areaCode"`
	PrefCode  string      `json:"prefCode"`
	Arrival   *HourMinute `json:"arrivalTime,omitempty"`
	Departure *HourMinute `json:"departureTime,omitempty"`
}

func TrainController(w http.ResponseWriter, r *http.Request) {
	stationId := chi.URLParam(r, "stationId")
	tableId := chi.URLParam(r, "tableId")
	trainId := chi.URLParam(r, "trainId")

	url := fmt.Sprintf("https://transit.yahoo.co.jp/timetable/%s/%s/%s", stationId, tableId, trainId)

	response, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(response.Body)

	data := doc.Find("#__NEXT_DATA__").Text()
	var jData Train
	json.Unmarshal([]byte(data), &jData)

	timetable := jData.Props.PageProps.TimetableStationTrainResult.Timetable
	trainId = timetable.TrainId
	trainName := timetable.Name
	stations := timetable.Stations

	lineName := strings.Split(doc.Find("h2").Text(), " ")[0]

	stationReturn := []StopStationReturn{}
	for _, station := range stations {
		stationReturn = append(stationReturn, StopStationReturn{Code: station.Code, Name: station.Name, AreaCode: station.AreaCode, PrefCode: station.PrefCode, Arrival: str2HourMinute(station.Arrival), Departure: str2HourMinute(station.Departure)})
	}

	ret := TrainReturn{Id: trainId, Name: trainName, LineName: lineName, Stations: stationReturn}

	body, _ := json.Marshal(ret)
	w.Write(body)
}
