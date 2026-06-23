package main

import (
	"log"
	"net/http"
	"stage/api"
	"stage/handlers"
	"stage/service"
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

	go func() {
		for {
			log.Println("Refreshing artist cache in background...")

			err := cache.Refresh()
			if err != nil {
				log.Println("Background cache refresh failed:", err)
			} else {
				log.Println("Background cache updated successfully")
			}

			time.Sleep(24 * time.Hour)
		}
	}()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", h.ArtistsHandler)
	http.HandleFunc("/singleArtists/", s.SingleArtistHandler)
	log.Println("Server currently running on port:http://localhost:8001")
	err = http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
