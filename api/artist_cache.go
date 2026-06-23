package api

import (
	"stage/models"
	"sync"
	"log"
	"os"
	"encoding/json"
)


func (a *ArtistCache) Refresh() error{
	artists, err := GetArtists()
	if err != nil {
		return  err
	}

	relation, err := GetRelations()
	if err != nil {
		return err
	}

	var newCache []models.FullArtist
	var artistMapByID = make(map[int]models.FullArtist)

	semaphore := make(chan struct{}, 5)

	for _, artist := range artists {

		locationMap := relation[artist.Id]

		// create slice with correct size
		locations := make([]models.LocationInfo, len(locationMap))

		var wg sync.WaitGroup

		// channel to safely collect results
		type result struct {
			location models.LocationInfo
		}

		results := make(chan result, len(locationMap))

		// launch goroutines
		for location, dates := range locationMap {

			wg.Add(1)

			go func(location string, dates []string) {

				defer wg.Done()

				semaphore <- struct{}{}

				defer func() {
					<-semaphore
				}()

				coordinates, err := a.geoCache.GeocodeLocation(location)
				if err != nil {
					log.Println("geocode failed:", err)
					coordinates = models.Geolocation{}
				}

				results <- result{
					location: models.LocationInfo{
						Name:  location,
						Lat:   coordinates.Lat,
						Lon:   coordinates.Lon,
						Dates: dates,
					},
				}

			}(location, dates)
		}

		// wait for all goroutines
		wg.Wait()

		// close results channel
		close(results)

		// safely fill locations slice
		i := 0
		for res := range results {
			locations[i] = res.location
			i++
		}

		info := models.FullArtist{
			Artist:         artist,
			DatesLocations: relation[artist.Id],
			Locations:      locations,
		}

		newCache = append(newCache, info)
		artistMapByID[artist.Id] = info
	}
	a.mu.Lock()
	a.artistsCache = newCache
	a.artistById = artistMapByID
	a.mu.Unlock()

	// SaveArtistsToCache()
	err = a.SaveArtistsToCache()
	if err != nil{
		return err
	}
	
	a.geoCache.SaveCacheToFile()
	return nil
}

func (a *ArtistCache) GetAllArtists()[]models.FullArtist{
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := make([]models.FullArtist, len(a.artistsCache))
	copy(out, a.artistsCache)
	return out
}

func (a *ArtistCache) LoadArtistsFromFile()error{
	data, err := os.ReadFile("artistsCache.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &a.artistsCache)
	if err != nil {
		return err
	}

	a.artistById = make(map[int]models.FullArtist)

	for _, artist := range a.artistsCache {
		a.artistById[artist.Artist.Id] = artist
	}

	return nil
}

func (a *ArtistCache) SaveArtistsToCache() error{
	a.mu.RLock()
	data, err := json.MarshalIndent(a.artistsCache, "", " ")
	a.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile("artistsCache.json", data, 0644)
}

func (a *ArtistCache) GetArtistByID(Id int) (models.FullArtist, bool){
	a.mu.RLock()
	defer a.mu.RUnlock()

	artist, found := a.artistById[Id]
	return artist, found
}