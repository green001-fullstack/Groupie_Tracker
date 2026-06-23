package api

import (
	"stage/models"
	"sync"
	"os"
	"encoding/json"
)

type GeoCache struct{
	geoCache map[string]models.Geolocation
	geoMutex sync.RWMutex
}

func NewGeoCache () *GeoCache{
	return &GeoCache{
		geoCache : make(map[string]models.Geolocation),
	}
}

func (g *GeoCache) Get(key string) (models.Geolocation, bool){
	g.geoMutex.RLock()
	defer g.geoMutex.RUnlock()

	result, found := g.geoCache[key]

	return result, found
}

func (g *GeoCache) Set(key string, value models.Geolocation){
	g.geoMutex.Lock()
	defer g.geoMutex.Unlock()
	g.geoCache[key] = value
}

func (g *GeoCache) LoadCacheFromFile() error {
	data, err := os.ReadFile("geocache.json")
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &g.geoCache)
}

func (g *GeoCache)SaveCacheToFile()error{
	g.geoMutex.RLock()
	data, err := json.MarshalIndent(g.geoCache, "", " ")
	g.geoMutex.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile("geocache.json", data, 0644)
}