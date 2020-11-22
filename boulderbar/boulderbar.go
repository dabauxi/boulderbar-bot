package boulderbar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var lastRefresh time.Time = time.Now()

const BOULDER_BAR_CUSTOMER_PAGE_WIEN = "https://shop.boulderbar.net:8080/modules/bbext/CurrentCustomer.php"
const BOULDER_BAR_CUSTOMER_PAGE_SBG = "https://shop.boulderbar-sbg.at:8081/modules/bbext/CurrentCustomer.php"

func CollectStatus() map[string]int {

	var boulderbars map[string]int = make(map[string]int)

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

	c.Visit(BOULDER_BAR_CUSTOMER_PAGE_WIEN)
	c.Visit(BOULDER_BAR_CUSTOMER_PAGE_SBG)

	return boulderbars
}

type GeneralInfo struct {
	Info [4]string
}

func GetGeneralInfo() GeneralInfo {
	var a [4]string
	a[0] = "This is a bot to fetch the available places at the boulderbars in Vienna from boulderbar.net."
	a[1] = "Due to the COVID-19 pandemic only a limited amount of people are allowed at a gym at the same time."
	a[2] = "Type in /status to get the current available places."
	a[3] = "Type in /locations to get the boulderbar locations."
	return GeneralInfo{Info: a}
}
