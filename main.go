package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/heroku/x/hmetrics/onload"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Session struct {
SessionId         string   `json:"session_id"`
Date              string   `json:"date"`
AvailableCapacity int      `json:"available_capacity"`
MinAgeLimit       int      `json:"min_age_limit"`
Vaccine           string   `json:"vaccine"`
Slots             []string `json:"slots"`
}

type VaccineFee struct {
Vaccine string `json:"vaccine"`
Fee     string `json:"fee"`
}

type Centers struct {
	CenterId     int    `json:"center_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	StateName    string `json:"state_name"`
	DistrictName string `json:"district_name"`
	BlockName    string `json:"block_name"`
	Pincode      int    `json:"pincode"`
	Lat          int    `json:"lat"`
	Long         int    `json:"long"`
	From         string `json:"from"`
	To           string `json:"to"`
	FeeType      string `json:"fee_type"`
	Sessions     []Session `json:"sessions"`
	VaccineFees []VaccineFee `json:"vaccine_fees,omitempty"`
}

type Response struct {
	Centers []Centers `json:"centers"`
}

func messageTelegram(bot_message string)  {
	var botToken = "1865363469:AAFyk4LsnskUMSKIYylJ1cx7cSXXikA8c0o"
	var botChatid = "578229642"
	var sendTextURL = "https://api.telegram.org/bot"+ botToken +"/sendMessage?chat_id=-"+ botChatid +"&text="+ url.QueryEscape(bot_message)
	log.Print(sendTextURL)
	_, err := http.Get(sendTextURL)
	if err != nil {
		log.Fatalln(err)
	}
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Print(body)
}

func task() {
	URL := "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?"

	payload := url.Values{}
	payload.Add("district_id", "392")
	payload.Add("date", time.Now().Format("02-01-2006"))

	//log.Println(bytes.NewBufferString(payload.Encode()))
	req, err := http.NewRequest("GET", URL + payload.Encode(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")

	q := req.URL.Query()
	q.Add("pincode", "400001")
	q.Add("date", time.Now().Format("30-12-2006"))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	//sb := string(body)
	//log.Printf(sb)

	//Convert json response to struct
	var data Response
	json.Unmarshal(body, &data)
	//log.Printf("Results: %v\n", data)
	//log.Printf("%d", len(data.Centers))

	//var center Message
	//var count int
	for i := 0; i < len(data.Centers); i++ {
		var for18 = false
		for j := 0; j < len(data.Centers[i].Sessions); j++ {
			if data.Centers[i].Sessions[j].MinAgeLimit ==18{
				for18 = true
			}
		}
		if for18 {
			//count ++
			messageTelegram(fmt.Sprintf("%#v\n", data.Centers[i]))
			//center.message += fmt.Sprintf("%#v\n", data.Centers[i])
		}
	}
	//log.Printf("%v\n", center)
}

func main() {
	port := "8080"

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	s := gocron.NewScheduler(time.UTC)
	s.Every(24).Seconds().Do(task)
	s.StartAsync()

	router.Run(":" + port)
}
