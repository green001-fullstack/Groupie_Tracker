package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stage/api"
	"stage/handlers"
	"stage/service"
	"sync"
	"time"
)

func main() {
	var err error

	geoCache := api.NewGeoCache()
	err = geoCache.LoadCacheFromFile()
	if err != nil {
		log.Println("No geo cache found")
	}

	cache := api.NewArtistCache(geoCache)
	err = cache.LoadArtistsFromFile()
	if err != nil {
		log.Println("No cache file found, starting fresh")
	}

	if len(cache.GetAllArtists()) == 0 {
		err := cache.Refresh()
		if err != nil {
			log.Fatal("Failed to refresh cache:", err)
		}
	}

	artistService := service.NewArtistService(cache)
	s := handlers.NewSingleArtist(artistService)

	h := &handlers.Handler{
		Service: artistService,
	}

	var wg sync.WaitGroup

	workerCtx, stopWorker := context.WithCancel(context.Background())

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select{
			case <- ctx.Done():
				return
			case	<- time.After(24 * time.Hour):
			}
			
			log.Println("Refreshing artist cache in background...")

			err := cache.Refresh()
			if err != nil {
				log.Println("Background cache refresh failed:", err)
			} else {
				log.Println("Background cache updated successfully")
			}

		}
	}(workerCtx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	srv := &http.Server{
		Addr: ":8001",
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", h.ArtistsHandler)
	http.HandleFunc("/singleArtists/", s.SingleArtistHandler)

	go func() {
		log.Println("Server currently running on port:http://localhost:8001")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-sigChan
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stopWorker()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Println(ctx.Err())
	}

	wg.Wait()

	log.Println("Server shutdown complete")
}
