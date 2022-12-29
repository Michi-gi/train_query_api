package main

import (
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/PuerkitoBio/goquery"
)


type Line struct {
	RailName	string	`json:"RailName"`
	Direction	string	`json:"Direction"`
	Source	string	`json:"Source"`
	RailId	string	`json:"RailId"`
	ServiceDayCode	string	`json:"ServiceDayCode"`
	RailTArget	string	`json:"RailTarget"`
}

type StationDetail struct {
	StationId string `json:"StationId"`
	StationInfo struct {
		RailGroup []Line `json:"RailGroup"`
	} `json:"StationInfo"`
}

type Station struct {
	Props	struct {
		PageProps	struct {
			FeatureWithRail	struct {
				Feature	struct {
					TransitSearchInfo	struct {
						Detail	StationDetail	`json:"Detail"`
					}	`json:"TransitSearchInfo"`
					Name	string	`json:"Name"`
				}	`json:"Feature"`
			}	`json:"FeaturewithRail"`
		}	`json:"PageProps"`
	}	`json:"props"`
}

type StationReturn struct {
	Id	string	`json:"id"`
	Name	string	`json:"name"`
	Lines	[]Line	`json:"lines"`
}

func StationController(w http.ResponseWriter, r *http.Request) {
	stationId := chi.URLParam(r, "stationId")

	url := fmt.Sprintf("https://transit.yahoo.co.jp/station/%s", stationId)
	doc, _ := goquery.NewDocument(url)

	data := doc.Find("#__NEXT_DATA__").Text()
	var jData Station
	json.Unmarshal([]byte(data), &jData)

	detail := jData.Props.PageProps.FeatureWithRail.Feature.TransitSearchInfo.Detail
	stationIdRet := detail.StationId
	stationName := jData.Props.PageProps.FeatureWithRail.Feature.Name
	lines := detail.StationInfo.RailGroup

	ret := StationReturn{Id: stationIdRet, Name: stationName, Lines: lines}
	body, _ := json.Marshal(ret)
	w.Write(body)
}
