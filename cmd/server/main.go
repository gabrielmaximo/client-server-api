package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type economiaApiResponse struct {
	bid string
}

type CotacaoEntity struct {
	ID  string `gorm:"primary key"`
	bid string
}

const EconomiaApiUrl = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {
	mux := http.NewServeMux()
	dialect := sqlite.Open(":memory:")
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&CotacaoEntity{})
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		res, err := http.NewRequestWithContext(ctx, "GET", EconomiaApiUrl, nil)
		if err != nil {
			panic(err)
		}

		var decodedResponse = new(economiaApiResponse)
		err = json.NewDecoder(res.Body).Decode(decodedResponse)
		if err != nil {
			panic(err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		db.WithContext(ctx).Create(CotacaoEntity{ID: uuid.New().String(), bid: decodedResponse.bid})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(decodedResponse)
		if err != nil {
			panic(err)
		}
	})

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
