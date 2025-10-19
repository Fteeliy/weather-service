package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Fteeliy/weather-service/internal/client/http/geocoding"
	openmeteo "github.com/Fteeliy/weather-service/internal/client/http/open_meteo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

const port = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	geocodingClient := geocoding.NewClient(&http.Client{Timeout: 10 * time.Second})
	openMeteoClient := openmeteo.NewClient(&http.Client{Timeout: 10 * time.Second})

	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")

		geoRes, err := geocodingClient.GetCoords(city)
		if err != nil {
			http.Error(w, "Failed to get coordinates", http.StatusInternalServerError)
			return
		}

		openMeteoRes, err := openMeteoClient.GetTemperature(geoRes.Latitude, geoRes.Longitude)
		if err != nil {
			http.Error(w, "Failed to get temperature", http.StatusInternalServerError)
			return
		}

		raw, err := json.Marshal(openMeteoRes)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(raw)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Error creating scheduler: %v", err)
	}

	jobs, err := initJons(s)
	if err != nil {
		log.Fatalf("Error initializing jobs: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Printf("Starting server on port %s\n", port)
		err := http.ListenAndServe(port, r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()

		fmt.Printf("Starting scheduler with job: %v\n", jobs[0].ID())
		s.Start()
	}()

	wg.Wait()
}

func initJons(scheduler gocron.Scheduler) ([]gocron.Job, error) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("Hello")
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}
