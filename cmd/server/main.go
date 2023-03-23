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
	Bid string
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
		ctx2, cancel2 := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel2()

		req, err := http.NewRequestWithContext(ctx2, "GET", EconomiaApiUrl, nil)
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

		ctx2, cancel2 = context.WithTimeout(context.Background(), 10*time.Millisecond)
		if ctx2.Err() != nil {
			panic(ctx2.Err())
		}
		defer cancel2()

		db.WithContext(ctx2).Create(CotacaoEntity{ID: uuid.New().String(), Bid: decodedResponse.USDBRL.Bid})
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
