package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi"
)

type Line struct {
	RailName     string `json:"railName"`
	Direction    string `json:"direction"`
	Source       string `json:"source"`
	GroupId      string `json:"groupId"`
	DriveDayKind string `json:"driveDayKind"`
}

type RouteInfo struct {
	RailName  string `json:"railName"`
	RailGroup []Line `json:"railGroup"`
}

type Station struct {
	Props struct {
		PageProps struct {
			DirectionDetail struct {
				StationName   string `json:"stationName"`
				DirectionItem struct {
					RouteInfos []RouteInfo `json:"routeInfos"`
				} `json:"directionItem"`
			} `json:"directionDetail"`
		} `json:"PageProps"`
	} `json:"props"`
}

type StationReturn struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Lines []Line `json:"lines"`
}

func StationController(w http.ResponseWriter, r *http.Request) {
	stationId := chi.URLParam(r, "stationId")

	url := fmt.Sprintf("https://transit.yahoo.co.jp/station/%s", stationId)
	response, _ := http.Get(url)
	doc, _ := goquery.NewDocumentFromReader(response.Body)

	data := doc.Find("#__NEXT_DATA__").Text()
	var jData Station
	json.Unmarshal([]byte(data), &jData)

	// detail := jData.Props.PageProps.DirectionDetail.Feature.TransitSearchInfo.Detail
	// stationIdRet := detail.StationId
	stationName := jData.Props.PageProps.DirectionDetail.StationName
	routes := jData.Props.PageProps.DirectionDetail.DirectionItem.RouteInfos
	var lines []Line
	for _, route := range routes {
		for _, line := range route.RailGroup {
			line.RailName = route.RailName
			lines = append(lines, line)
		}
	}

	ret := StationReturn{Id: stationId, Name: stationName, Lines: lines}
	body, _ := json.Marshal(ret)
	w.Write(body)
}
