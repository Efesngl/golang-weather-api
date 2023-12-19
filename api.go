package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Weather struct {
	Forecast struct {
		Forecastday []struct {
			Day struct {
				MaxtempC float64 `json:"maxtemp_c"`
				MintempC float64 `json:"mintemp_c"`
				AvgtempC float64 `json:"avgtemp_c"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	var w Weather
	weatherMap := make(map[string]string)
	r := gin.Default()
	r.GET("/today", func(c *gin.Context) {
		link := "http://api.weatherapi.com/v1/forecast.json?key=fa802e2cc1714f7281c100339231912&q=Izmir&days=1&aqi=no&alerts=no"
		getWeather(&w,link,weatherMap)
		c.JSON(http.StatusOK,weatherMap)
	})
	r.GET("/yesterday", func(c *gin.Context) {
		yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		link := fmt.Sprintf("http://api.weatherapi.com/v1/history.json?key=fa802e2cc1714f7281c100339231912&q=Izmir&dt=%v", yesterday)
		getWeather(&w,link,weatherMap)
		c.JSON(http.StatusOK,weatherMap)
	})
	r.Run()
}
func getWeather(w *Weather,link string, wMap map[string]string)  {
	req, err := http.NewRequest("GET", link, nil)
	if(err!=nil){
		fmt.Println(err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if(err!=nil){
		fmt.Println(err)
	}
	resBody, _ := io.ReadAll(res.Body)
	json.Unmarshal(resBody, &w)
	wMap["İl"] = "İzmir"
	wMap["En düşük sıcaklık"] = fmt.Sprintf("%v", w.Forecast.Forecastday[0].Day.MintempC)
	wMap["En yüksek sıcaklık"] = fmt.Sprintf("%v", w.Forecast.Forecastday[0].Day.MaxtempC)
	wMap["Ortalama sıcaklık"] = fmt.Sprintf("%v", w.Forecast.Forecastday[0].Day.AvgtempC)
	defer res.Body.Close()
}
