package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

var boulderbars map[string]int = make(map[string]int)
var lastRefresh time.Time = time.Now()

func collect() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("shop.boulderbar.net:8080", "shop.boulderbar-sbg.at:8081"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("html body div.progress-radial2", func(e *colly.HTMLElement) {
		gym := e.Text
		gymName := strings.TrimSpace(gym[:len(gym)-len("plätze frei")-2])
		freePlaces, err := strconv.Atoi(gym[len(gym)-len("plätze frei")-2 : len(gym)-len("plätze frei")])
		if err != nil {
			freePlaces = -1
		}
		if freePlaces > 49 {
			gymName = gymName[:len(gymName)-len("über")]
		}
		boulderbars[gymName] = freePlaces
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
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

type conf struct {
	Token string `yaml:"token"`
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	var c conf
	c.getConf()

	b, err := tb.NewBot(tb.Settings{
		Token:  c.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
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
