package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stebunting/hfx-backend/currency"
	"github.com/stebunting/hfx-backend/model"
	"github.com/stebunting/hfx-backend/scraper"
)

type Route struct {
	db *pg.DB
}

type Response struct {
	Status  string
	Details interface{}
}

func ConfigRoutes(db *pg.DB) Route {
	return Route{db}
}

func (s *Route) Wake(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		response, err := json.Marshal(Response{
			Status:  "OK",
			Details: "Awake",
		})
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s *Route) GetRate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		scrape := scraper.NewScraper()

		f := r.URL.Query()["from_code"]
		if len(f) < 1 {
			s.returnError(w, "From currency required")
			return
		}
		from, err := currency.NewCode(f[0])
		if err != nil {
			s.returnError(w, "Invalid from currency")
			return
		}

		t := r.URL.Query()["to_code"]
		if len(t) < 1 {
			s.returnError(w, "To currency required")
			return
		}
		to, err := currency.NewCode(t[0])
		if err != nil {
			s.returnError(w, "Invalid to currency")
			return
		}

		d := r.URL.Query()["date"]
		if len(d) < 1 {
			s.returnError(w, "Date required")
			return
		}
		date, err := scrape.ParseDate(d[0])
		if err != nil {
			s.returnError(w, "Invalid date")
			return
		}

		var result []model.Exchange
		err = s.db.Model(&model.Exchange{}).
			Column("date", "from_code", "to_code", "rate").
			Where("from_code = ?", from.String()).
			Where("to_code = ?", to.String()).
			Where("date = ?", date).
			Select(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}

		if len(result) == 0 {
			result, err = scrape.ScrapeCurrency(from, to, date)
			if err != nil {
				log.Fatal(err)
			}

			_, err = s.db.Model(&result).Insert()
			if err != nil {
				log.Fatal(err)
			}
		}

		var res Response
		if len(result) > 0 {
			res = Response{
				Status:  "OK",
				Details: result[0],
			}
		} else {
			res = Response{
				Status:  "Error",
				Details: "no rate found",
			}
		}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s *Route) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		var currencies []model.Currency
		s.db.Model(&model.Currency{}).Order("code ASC").Select(&currencies)

		response, err := json.Marshal(Response{
			Status:  "OK",
			Details: currencies,
		})
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s *Route) UpdateCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		s.ScrapeCurrencies()

		res, _ := json.Marshal(Response{
			Status:  "OK",
			Details: "Currencies Updated",
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}

func (s *Route) DbInit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")

		models := []interface{}{
			(*model.Currency)(nil),
			(*model.Exchange)(nil),
		}

		for _, model := range models {
			err := s.db.Model(model).CreateTable(&orm.CreateTableOptions{})
			if err != nil {
				log.Fatal(err)
			}
		}

		s.ScrapeCurrencies()

		response, err := json.Marshal(Response{
			Status: "OK",
		})
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s *Route) ScrapeCurrencies() {
	scrape := scraper.NewScraper()

	result, err := scrape.ScrapeCurrencies()
	if err != nil {
		log.Fatal(err)
	}

	s.db.Model(&model.Currency{}).Where("TRUE").Delete()
	_, err = s.db.Model(result).Insert()
	if err != nil {
		log.Fatal(err)
	}

	result, err = scrape.ScrapeSymbols()
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.db.Model(result).ExcludeColumn("name").WherePK().Update()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Route) returnError(w http.ResponseWriter, errorDescription string) {
	w.WriteHeader(http.StatusBadRequest)

	response, _ := json.Marshal(Response{
		Status:  "Error",
		Details: errorDescription,
	})
	w.Write([]byte(response))
}
