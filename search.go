package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)


type SearchResult struct {
	Props	struct {
		PageProps	struct {
			StationFeatures []struct {
				Stations []struct {
					Id	string	`json:"Id"`
					Name	string	`json:"Name"`
					Railway	string	`json:"Railway"`
				}	`json:"Station"`
			}	`json:"stationFeatures"`
		}	`json:"PageProps"`
	}	`json:"props"`
}

type SearchReturn struct {
	Id	string	`json:"id"`
	Name	string	`json:"name"`
	Railways	[]string `json:"railways"`
}

func SearchController(w http.ResponseWriter, r *http.Request) {
	searchWords := r.URL.Query()["q"]
	searchWord := ""
	if len(searchWords) > 0 {
		searchWord = searchWords[0]
	}

	url := fmt.Sprintf("https://transit.yahoo.co.jp/timetable/search?q=%s", url.QueryEscape(searchWord))
	doc, _ := goquery.NewDocument(url)

	data := doc.Find("#__NEXT_DATA__").Text()
	var jData SearchResult
	json.Unmarshal([]byte(data), &jData)

	queryResult := jData.Props.PageProps.StationFeatures
	resultMap := make(map[string]*SearchReturn)
	for _, feature := range queryResult {
		for _, station := range feature.Stations {
			if _, exist := resultMap[station.Id]; !exist {
				resultMap[station.Id] = &SearchReturn{Id: station.Id, Name: station.Name, Railways: []string{}}
			}
			if !exists(resultMap[station.Id].Railways, station.Railway) {
				resultMap[station.Id].Railways = append(resultMap[station.Id].Railways, station.Railway)
}
		}
	}

	body, _ := json.Marshal(resultMap)
	w.Write(body)
}

func exists(slice []string, target string)bool {
	for _, ele := range slice {
		if ele == target {
			return true
		}
	}
	return false
}
