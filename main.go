package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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

var regexWeight = regexp.MustCompile(`(?im)weight\s*:\s*([0-9]{2,4}(?:\s*[-–]\s*[0-9]{2,4})?\s*(?:g|grs?|grams?))`)
var regexShape = regexp.MustCompile(`(?im)(?:shape|form)\s*:\s*(diamond|round|drop\s*/\s*tear)`)
var regexFace = regexp.MustCompile(`(?im)\b(?:face|faces|heads|cara)\b\s*(?:[:·\-–])\s*` + `(.+?)` + `\s*(?:\bwith\b|\b[A-Z][A-Z0-9 /-]{2,20}\b\s*(?:[:·\-–])|$)`)
var regexFiberFallback = regexp.MustCompile(`(?im)\b(fiberglass(?:\s+silver)?|fiber\s+glass(?:\s+\d+k)?(?:\s+silver)?)\b`)
var brand = "NOX"

func ScrapeRacketSpecificInfo(url string) (Weight string, Shape string, Material string) {

	weight := ""
	shape := ""
	material := ""

	c := colly.NewCollector(colly.AllowedDomains("noxsport.com"))
	c.OnHTML("div.cc-accordion-item__content", func(e *colly.HTMLElement) {
		Description := strings.Join(strings.Fields(e.Text), " ")

		if weight != "" && shape != "" && material != "" {
			return
		}

		if !strings.Contains(strings.ToUpper(Description), "FACE:") && !strings.Contains(strings.ToUpper(Description), "FACES:") && !strings.Contains(strings.ToUpper(Description), "HEADS:") && !strings.Contains(strings.ToUpper(Description), "CARA:") {
			return
		}

		tmpWeight := regexWeight.FindStringSubmatch(Description)
		if len(tmpWeight) >= 2 {
			weight = strings.ReplaceAll(strings.TrimSpace(tmpWeight[1]), " ", "")
		}

		tmpShape := regexShape.FindStringSubmatch(Description)
		if len(tmpShape) >= 2 {
			shape = strings.TrimSpace(tmpShape[1])
			shape = strings.ReplaceAll(shape, " ", "")
			shape = strings.ReplaceAll(shape, "drop/tear", "Drop/Tear")
			shape = strings.ReplaceAll(shape, "diamond", "Diamond")
			shape = strings.ReplaceAll(shape, "round", "Round")
		}

		tmpFace := regexFace.FindStringSubmatch(Description)
		if len(tmpFace) >= 2 {
			material = strings.TrimSpace(tmpFace[1])
		} else {
			tmpFiberFallback := regexFiberFallback.FindStringSubmatch(Description)
			if len(tmpFiberFallback) >= 2 {
				material = strings.TrimSpace(tmpFiberFallback[1])
			}
		}

	})

	if err := c.Visit(url); err != nil {
		fmt.Println("Visiting error: ", err)
	}
	return weight, shape, material
}

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

		racketWeight, racketShape, racketMaterial := ScrapeRacketSpecificInfo(item.RacketPage)
		item.Weight = racketWeight
		item.Shape = racketShape
		item.Material = racketMaterial

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
	SliceOfSeries := []string{"signature-series-luxury", "serie-nfa-padel", "serie-classic", "series-exclusive-edition-padel", "serie-advance", "serie-essential", "serie-ultralight"}
	for i := 0; i < len(SliceOfSeries); i++ {
		ScrapeRacketPage("https://noxsport.com/en/collections/" + SliceOfSeries[i])
	}

}
