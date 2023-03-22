package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type EconomiaApiResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type ServerResponse struct {
	Bid string `json:"bid"`
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
		ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", EconomiaApiUrl, nil)
		if err != nil {
			panic(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer func(b io.ReadCloser) {
			err := b.Close()
			if err != nil {
				panic(err)
			}
		}(res.Body)

		var decodedResponse = new(EconomiaApiResponse)
		err = json.NewDecoder(res.Body).Decode(decodedResponse)
		if err != nil {
			panic(err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		db.WithContext(ctx).Create(CotacaoEntity{ID: uuid.New().String(), bid: decodedResponse.USDBRL.Bid})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(ServerResponse{Bid: decodedResponse.USDBRL.Bid})
		if err != nil {
			panic(err)
		}
	})

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
