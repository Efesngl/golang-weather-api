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
	r.GET("/getweather/:method", func(c *gin.Context) {
		method := c.Param("method")
		if getWeather(&w, weatherMap, method){
			c.JSON(http.StatusOK, weatherMap)
		}else{
			c.Status(http.StatusNotFound)
		}
		
	})
	r.Run()
}
func getWeather(w *Weather, wMap map[string]string, method string) bool {
	methods:=[]string{"today","yesterday"}
	isMethodValid:=false
	for i := 0; i < len(methods); i++ {
		if method==methods[i] {
			isMethodValid=true
		}
	}
	if !isMethodValid{
		return false
	}
	link := "http://api.weatherapi.com/v1/forecast.json?key=fa802e2cc1714f7281c100339231912&q=Izmir&days=1&aqi=no&alerts=no"
	if method == "yesterday" {
		yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		link = fmt.Sprintf("http://api.weatherapi.com/v1/history.json?key=fa802e2cc1714f7281c100339231912&q=Izmir&dt=%v", yesterday)
	}
	req, _ := http.NewRequest("GET", link, nil)
	client := http.Client{}
	res, _ := client.Do(req)
	resBody, _ := io.ReadAll(res.Body)
	json.Unmarshal(resBody, &w)
	wMap["İl"] = "İzmir"
	wMap["En düşük sıcaklık"] = fmt.Sprintf("%v", w.Forecast.Forecastday[0].Day.MintempC)
	wMap["En yüksek sıcaklık"] = fmt.Sprintf("%v", w.Forecast.Forecastday[0].Day.MaxtempC)
	wMap["Ortalama sıcaklık"] = fmt.Sprintf("%v", w.Forecast.Forecastday[0].Day.AvgtempC)
	defer res.Body.Close()
	return true
}
