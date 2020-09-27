package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	tb "gopkg.in/tucnak/telebot.v2"
)

var boulderbars map[string]int = make(map[string]int)
var lastRefresh time.Time = time.Now()

func collect() {
	c := colly.NewCollector(
		colly.AllowedDomains("shop.boulderbar.net:8080", "shop.boulderbar-sbg.at:8081"),
	)

	c.OnHTML("html body div.progress-radial2", func(e *colly.HTMLElement) {
		gym := strings.ToLower(e.Text)
		gymName := "unknown"
		if strings.Contains(gym, "hauptbahnhof") {
			gymName = "Hauptbahnhof"
		}
		if strings.Contains(gym, "hannovergasse") {
			gymName = "Hannovergasse"
		}
		if strings.Contains(gym, "wienerberg") {
			gymName = "Wienerberg"
		}
		if strings.Contains(gym, "salzburg") {
			gymName = "Salzburg"
		}
		freePlaces, err := strconv.Atoi(gym[len(gymName) : len(gym)-len("plätze frei")])
		if err != nil {
			freePlaces, err = strconv.Atoi(gym[len(gymName)+len("über")+1 : len(gym)-len("plätze frei")])
			if err != nil {
				freePlaces = -1
			}

		}
		boulderbars[gymName] = freePlaces
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://shop.boulderbar.net:8080/modules/bbext/CurrentCustomer.php")
	c.Visit("https://shop.boulderbar-sbg.at:8081/modules/bbext/CurrentCustomer.php")
}

func createResponse() string {
	responseSlice := make([]string, 5)
	var builder strings.Builder
	if len(boulderbars) > 0 {
		responseSlice[0] = "Current available places\n\n"
	}
	for key, value := range boulderbars {
		if value > 0 {
			builder.WriteString(key)
			builder.WriteString(": ")
			if value > 49 {
				builder.WriteString("more than ")
			}
			builder.WriteString(strconv.Itoa(value))
			builder.WriteString("\n")
			responseSlice = append(responseSlice, builder.String())
			builder.Reset()
		}
	}
	builder.Reset()
	sort.Strings(responseSlice)
	for _, value := range responseSlice {
		builder.WriteString(value)
	}
	return builder.String()
}

func test() {
	collect()
	fmt.Println(createResponse())
}

func main() {

	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: webhook,
	})

	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "This is a bot to fetch the available places at the boulderbars in Vienna from boulderbar.net.")
		b.Send(m.Sender, "Due to the COVID-19 pandemic only a limited amount of people are allowed at a gym at the same time.")
		b.Send(m.Sender, "Type in /status to get the current available places.")
	})

	b.Handle("/status", func(m *tb.Message) {
		if time.Now().After(lastRefresh.Add(time.Second * 60)) {
			collect()
			lastRefresh = time.Now()
		}
		b.Send(m.Sender, createResponse())
	})

	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "Type in /status to get the current utilization.")
	})

	b.Start()
}
