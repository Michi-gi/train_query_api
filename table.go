package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi"
)

type HourData struct {
	Hour         string `json:"hour"`
	MinTimeTable []struct {
		Minute        string `json:"minute"`
		KindId        string `json:"kindId"`
		DestinationId string `json:"destinationId"`
		TrainId       string `json:"trainId"`
		TrainName     string `json:"trainName"`
		FirstStation  string `json:"firstStation"`
		Extra         string `json:"extraTrain"`
		VendorTrainId string `json:"vendorTrainId"`
	} `json:"minTimeTable"`
}

type Table struct {
	Props struct {
		PageProps struct {
			TimetableItem struct {
				Master struct {
					Kind        []IdName `json:"kind"`
					Destination []IdName `json:"destination"`
				} `json:"master"`
				DirectionName  string     `json:"directionName"`
				RailName       string     `json:"railName"`
				HourTimeTables []HourData `json:"hourTimeTable"`
			} `json:"timetableItem"`
		} `json:"PageProps"`
	} `json:"props"`
}

type TrainInTable struct {
	Time           string `json:"time"`
	Kind           string `json:"kind"`
	Dest           string `json:"destination"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	IsFirstStation bool   `json:"isFirstStation"`
	IsExtra        bool   `json:"isExtra"`
	VendorNumber   string `json:"vendorNumber"`
}

type TableReturn struct {
	StationId    string            `json:"stationId"`
	TableId      string            `json:"tableId"`
	Direction    string            `json:"direction"`
	LineName     string            `json:"lineName"`
	Table        []TrainInTable    `json:"table"`
	DayOfWeekMap map[string]string `json:"dayOfWeekMap"`
}

func TableController(w http.ResponseWriter, r *http.Request) {
	stationId := chi.URLParam(r, "stationId")
	tableId := chi.URLParam(r, "tableId")
	kind := r.URL.Query()["kind"]

	url := fmt.Sprintf("https://transit.yahoo.co.jp/timetable/%s/%s", stationId, tableId)
	if len(kind) > 0 {
		url += fmt.Sprintf("?kind=%s", kind[0])
	}

	response, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(response.Body)

	dayOfWeekMap := make(map[string]string)
	dayOfWeekEles := doc.Find(".navDayOfWeek > li")
	dayOfWeekEles.Each(func(index int, s *goquery.Selection) {
		linkEle := s.Find("span")
		var kindNum string
		if linkEle.Length() == 0 {
			linkEle = s.Find("a")
			kindref, _ := linkEle.Attr("href")
			kindNum = kindref[len(kindref)-1:]
		} else if len(kind) > 0 {
			kindNum = kind[0]
		} else {
			kindNum = "-1"
		}
		dayOfWeekMap[kindNum] = linkEle.Text()
	})
	data := doc.Find("#__NEXT_DATA__").Text()
	var jData Table
	json.Unmarshal([]byte(data), &jData)

	tableData := jData.Props.PageProps.TimetableItem
	master := tableData.Master

	kindMap := make(map[string]string)
	for _, kind := range master.Kind {
		kindMap[kind.Id] = kind.Name
	}

	destMap := make(map[string]string)
	for _, dest := range master.Destination {
		destMap[dest.Id] = dest.Name
	}

	direction := tableData.DirectionName
	line := tableData.RailName

	var trains []TrainInTable
	for _, hourData := range tableData.HourTimeTables {
		for _, train := range hourData.MinTimeTable {
			var t TrainInTable
			hour, _ := strconv.Atoi(hourData.Hour)
			minute, _ := strconv.Atoi(train.Minute)
			t.Time = fmt.Sprintf("%02d:%02d", hour, minute)
			t.Kind = kindMap[train.KindId]
			t.Dest = destMap[train.DestinationId]
			t.Id = train.TrainId
			t.Name = train.TrainName
			t.IsFirstStation = (train.FirstStation == "true")
			t.IsExtra = (train.Extra == "true")
			t.VendorNumber = train.VendorTrainId

			trains = append(trains, t)
		}
	}

	ret := TableReturn{StationId: stationId, TableId: tableId, Direction: direction, LineName: line, Table: trains, DayOfWeekMap: dayOfWeekMap}

	body, _ := json.Marshal(ret)
	w.Write(body)
}
