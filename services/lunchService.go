package services

import (
	"encoding/json"
	"ibuYemekBotu/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func GetLunchList(day string) string {
	log.Println("Getting lunch list of " + day)

	var url = os.Getenv("IbuYemekApiUrl")
	log.Println("Api Url: ", url+"/day/"+day)
	resp, err := http.Get(url + "/day/" + day)
	if err != nil {
		log.Println("Getting Lunch From Api Error: ", err)
	}

	lunch, err := parseResponseBody(resp)

	emptyLunch := models.Lunch{Corba: "", AnaYemek: "", YardimciAnaYemek: "", YanYemek1: "", YanYemek2: ""}
	if lunch == emptyLunch {
		return ""
	}

	lunchString := "Çorba: " + lunch.Corba + "\n" +
		"Ana Yemek: " + lunch.AnaYemek + "\n" +
		"İkinci Yemek: " + lunch.YardimciAnaYemek + "\n" +
		"Yan Yemek: " + lunch.YanYemek1 + "\n" +
		"Yan Yemek: " + lunch.YanYemek2

	if day == "today" {
		t := time.Now()
		lunchString = t.Format("02/01/2006") + "\n" + lunchString
		return lunchString
	} else {
		t := time.Now()
		t = t.AddDate(0, 0, 1)
		lunchString = t.Format("02/01/2006") + "\n" + lunchString
		return lunchString
	}
}

func parseResponseBody(resp *http.Response) (models.Lunch, error) {
	var lunchResponse models.LunchResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Reading Body Error: ", err)
		return models.Lunch{}, err
	}

	err = json.Unmarshal(body, &lunchResponse)

	return lunchResponse.Data, err

}
