package api

import(
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stage/models"
	"stage/utils"
	"strings"
)

func (g *GeoCache)GeocodeLocation(location string) (models.Geolocation, error) {
	formattedLocation := utils.FormatForGeocoding(location)

	cacheKey := strings.ToLower(formattedLocation)
	result, found := g.Get(cacheKey)
	if found {
		return result, nil
	}

	url := "https://nominatim.openstreetmap.org/search?q=" +
		strings.ReplaceAll(formattedLocation, " ", "+") +
		"&format=json&limit=1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return models.Geolocation{}, err
	}
	req.Header.Set("User-Agent", "groupie-tracker-app")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Geolocation{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Geolocation{}, err
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println(string(body[:200]))

	var data []geoResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.Geolocation{}, err
	}

	if len(data) == 0 {
		fmt.Println("No location found")
		return models.Geolocation{}, fmt.Errorf("no location found")
	}

	result = models.Geolocation{
		Lat: data[0].Lat,
		Lon: data[0].Lon,
	}

	g.Set(cacheKey, result)

	return result, nil
}

func GetArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error fetching api", err)
		return []models.Artist{}, err
	}
	defer resp.Body.Close()

	fmt.Println("Status code is:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad Status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Artist{}, err
	}

	var artists []models.Artist

	err = json.Unmarshal(body, &artists)
	if err != nil {
		return []models.Artist{}, err
	}
	return artists, nil
}

func GetRelations() (map[int]map[string][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad Status %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var relations models.Relations

	err = json.Unmarshal(data, &relations)
	if err != nil {
		return nil, err
	}

	relationMap := make(map[int]map[string][]string)

	for _, item := range relations.Index {
		relationMap[item.Id] = item.DatesLocations
	}
	return relationMap, nil
}

