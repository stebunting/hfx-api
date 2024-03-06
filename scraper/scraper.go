package scraper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stebunting/hfx-backend/currency"
	"github.com/stebunting/hfx-backend/model"
)

type Scraper struct {
	codeValidator regexp.Regexp
}

func NewScraper() Scraper {
	re, err := regexp.Compile("^[A-Z]{3}$")
	if err != nil {
		panic(err)
	}

	return Scraper{
		codeValidator: *re,
	}
}

type XeResponse struct {
	PageProps XePageProps `json:"pageProps"`
}

type XePageProps struct {
	HistoricRates []XeHistoricRate `json:"historicRates"`
}

type XeHistoricRate struct {
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
}

func (s *Scraper) GetCurrency(from currency.Code, to currency.Code, date time.Time) ([]model.Exchange, error) {
	url, err := url.Parse("https://xe.com/_next/data/Y8_6DNxLmlBwzHvjNAjvA/en/currencytables.json")
	if err != nil {
		panic(err)
	}
	query := url.Query()
	query.Set("from", from.String())
	query.Set("date", s.FormatDate(date))
	url.RawQuery = query.Encode()

	response, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	xe := XeResponse{}
	err = json.Unmarshal(body, &xe)
	if err != nil {
		panic(err)
	}

	var data []model.Exchange
	for _, r := range xe.PageProps.HistoricRates {
		code, err := currency.NewCode(r.Currency)
		if err == nil {
			if code == to {
				data = append(data, model.Exchange{
					Date:     date,
					FromCode: from.String(),
					ToCode:   code.String(),
					Rate:     r.Rate,
				})
			}
		}
	}

	return data, nil
}

func (s *Scraper) ScrapeCurrency(from currency.Code, to currency.Code, date time.Time) ([]model.Exchange, error) {
	url, err := url.Parse("https://www.xe.com/currencytables")
	if err != nil {
		panic(err)
	}
	query := url.Query()
	query.Set("from", from.String())
	query.Set("date", s.FormatDate(date))
	url.RawQuery = query.Encode()

	document, err := s.getDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	var data []model.Exchange
	tableClass := "#table-section > section > div > div > table > tbody > tr"
	document.Find(tableClass).Each(func(i int, sel *goquery.Selection) {
		code, err := currency.NewCode(sel.Find("th a").Text())
		if err == nil {
			rate := s.floatFromStr(sel.Find("td").Eq(1).Text())

			if code == to {
				data = append(data, model.Exchange{
					Date:     date,
					FromCode: from.String(),
					ToCode:   code.String(),
					Rate:     rate,
				})
			}
		}
	})

	return data, nil
}

func (s *Scraper) ScrapeCurrencies() (*[]model.Currency, error) {
	url, err := url.Parse("https://www.xe.com/iso4217.php")
	if err != nil {
		return nil, err
	}

	document, err := s.getDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	var currencies []model.Currency

	selector := "table#currencyTable > tbody > tr"
	document.Find(selector).Each(func(i int, sel *goquery.Selection) {
		c := sel.Find("td").Eq(0).Find("a").Text()
		code, err := currency.NewCode(strings.Replace(c, "*", "", -1))
		if err == nil {
			name := sel.Find("td").Eq(1).Text()

			if name != "" {
				newCurrency := model.Currency{
					Code: code.String(),
					Name: name,
				}
				currencies = append(currencies, newCurrency)
			}
		}
	})

	return &currencies, nil
}

func (s *Scraper) ScrapeSymbols() (*[]model.Currency, error) {
	url, err := url.Parse("https://www.xe.com/symbols.php")
	if err != nil {
		return nil, err
	}

	document, err := s.getDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	var symbols []model.Currency

	sectionClass := "section > section > ul > li"
	document.Find(sectionClass).Each(func(i int, sel *goquery.Selection) {
		code, err := currency.NewCode(sel.Find("div").Eq(3).Text())
		if err == nil {
			symbol := sel.Find("div").Eq(7).Text()

			if symbol != "" {
				newSymbol := model.Currency{
					Code:   code.String(),
					Symbol: symbol,
				}
				symbols = append(symbols, newSymbol)
			}
		}
	})

	return &symbols, nil
}

func (s *Scraper) getDocument(url *url.URL) (*goquery.Document, error) {
	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (s *Scraper) floatFromStr(n string) float64 {
	number := strings.Replace(n, ",", "", -1)
	f, err := strconv.ParseFloat(number, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (s *Scraper) ParseDate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, date+"T00:00:00.000Z")
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func (s *Scraper) FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}
