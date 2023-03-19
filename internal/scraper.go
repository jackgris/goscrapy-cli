/*
	Copyright © 2023 Gabriel Pozo

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package internal

import (
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"

	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// This structs will be use for match JSON data in the web page
type Product struct {
	MainEntity  MainEntity `json:"mainEntityOfPage"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Offers      Offers
}

type MainEntity struct {
	Id string `json:"@id"`
}

type Offers struct {
	Price          string `json:"price"`
	Availability   string `json:"availability"`
	InventoryLevel InventoryLevel
}

type InventoryLevel struct {
	Stock string `json:"value"`
}

// ///////////////////////////////////////////////////////////////////////////////////////
// Setup data
// ///////////////////////////////////////////////////////////////////////////////////////
type Wholesalers struct {
	Login        string `json:"login"`
	User         string `json:"user"`
	Pass         string `json:"pass"`
	Searchpage   string `json:"searchpage"`
	Name         string `json:"name"`
	EndPhrase    string `json:"endphrase"`
	EndPhraseDiv string `json:"endphrasediv"`
}

type DBProduct struct {
	Id          string
	Name        string
	Image       string
	Description string
	Price       string
	Stock       string
	Wholesaler  string
}

//////////////////////////////////////////////////////////////////////////////////////////////

// GetData will get data from the web of wholesalers, and save that information on the database
func GetData(file *os.File, w Wholesalers, log *logrus.Logger) error {

	// Starting data collector
	c := colly.NewCollector()
	err := c.Limit(&colly.LimitRule{
		Delay:        5 * time.Second,
		DomainRegexp: w.Searchpage + "*",
	})
	if err != nil {
		log.Error("Getdata: " + err.Error())
	}
	// With this we know when is the last page of catalog
	end := false

	// Authenticate
	err = c.Post(w.Login, map[string]string{"username": w.User, "password": w.Pass})
	if err != nil {
		log.Fatal("Get Data authenticate: ", err)
	}

	// Attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("Response received: ", r.StatusCode, " URL: ", r.Request.URL)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {

		// Goquery selection of the HTMLElement is in e.DOM
		goquerySelection := e.DOM

		// Check here, there are or there are no products
		goquerySelection.Find(w.EndPhraseDiv).Each(func(i int, el *goquery.Selection) {

			if end {
				return
			}
			not := w.EndPhrase
			end = strings.Contains(el.Text(), not)
		})

		if end {
			return
		}

		// Finding JSON data from scripts
		goquerySelection.Find("script").Each(func(i int, el *goquery.Selection) {

			// Create json struct
			p := Product{}
			_ = json.Unmarshal([]byte(el.Text()), &p)

			// Is data is ok, saving product on database
			if p.MainEntity.Id != "" {
				// saving data here
				//saveData(collection, ctx, p)
				product := DBProduct{
					Id:          p.MainEntity.Id,
					Name:        p.Name,
					Image:       p.Image,
					Description: p.Description,
					Price:       p.Offers.Price,
					Stock:       p.Offers.InventoryLevel.Stock,
					Wholesaler:  w.Name,
				}

				jsonData, err := json.Marshal(product)
				if err != nil {
					log.Warn("Error marshaling data:", err)
					return
				}

				if _, err = file.Write(jsonData); err != nil {
					log.Warn("Can´t save product: ", err)
				}

				if n, err := file.WriteString(",\n"); err != nil {
					log.Fatalln(n, err)
				}
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("OnRequest")
	})

	// The approach can be changed, at least for now, this will look for pages until the maxNumber
	// of pages or when the scraper detects there are no more pages to visit
	maxNumber := 1000
	for i := 1; i < maxNumber; i++ {
		// Check when there are no products
		if end {
			log.Debug("Searching end")
			break
		}

		num := strconv.Itoa(i)
		URL := w.Searchpage + num

		err := c.Visit(URL)
		if err != nil {
			log.Debug("Error visiting site: ", err)
		}
	}

	log.Debug(c.String())

	return nil
}
