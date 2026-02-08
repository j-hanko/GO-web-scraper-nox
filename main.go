package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Racket struct {
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Price      string `json:"price"`
	ImageUrl   string `json:"imageUrl"`
	RacketPage string `json:"racketPage"`
	Weight     string `json:"weight"`
	Shape      string `json:"shape"`
	Material   string `json:"material"`
	Series     string `json:"series"`
}

var brand = "NOX"

func ScrapeRacketPage(url string) {
	var items []Racket
	c := colly.NewCollector(colly.AllowedDomains("noxsport.com"))

	series := strings.Split(url, "/")[5]
	series = strings.Replace(series, "-", " ", -1)
	series = strings.Title(series)
	series = strings.ReplaceAll(series, " ", "")

	c.OnHTML("div.filters-adjacent div.block-inner-inner", func(e *colly.HTMLElement) {
		item := Racket{
			Brand:      brand,
			Model:      e.ChildText("div.product-block__title"),
			Price:      e.ChildText("span.money"),
			RacketPage: "https://noxsport.com" + e.ChildAttr("a", "href"),
		}
		item.Series = series

		imgSrcSet := e.ChildAttr("img.inline-image__image", "srcset")
		if imgSrcSet == "" {
			imgSrcSet = e.ChildAttr("img.inline-image__image", "data-srcset")
		}
		imgSrcSet = strings.SplitN(imgSrcSet, "?", 2)[0]
		imgSrcSet = strings.TrimPrefix(imgSrcSet, "//")
		item.ImageUrl = "https://" + imgSrcSet

		items = append(items, item)
	})

	if err := c.Visit(url); err != nil {
		fmt.Println("Visiting error: ", err)
	}

	content, err := json.Marshal(items)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err := os.WriteFile(brand+"Racket"+series+".json", content, 0644); err != nil {
		fmt.Println("File saving error: ", err)
	}
}

func main() {
	ScrapeRacketPage("https://noxsport.com/en/collections/signature-series-luxury")
}
