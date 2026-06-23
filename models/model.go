package models

import "html/template"


type Artist struct {
    Id           int	`json:"id"`
    Image        string	`json:"image"`
    Name         string	`json:"name"`
    Members      []string	`json:"members"`
    CreationDate int	`json:"creationDate"`
    FirstAlbum   string	`json:"firstAlbum"`
}

type RelationItem struct{
	Id int	`json:"id"`
	DatesLocations map[string][]string	`json:"datesLocations"`
}

type Relations struct{
	Index []RelationItem `json:"index"`
}

type LocationItem struct{
	Id int `json:"id"`
	Locations []string `json:"locations"`
	DatesUrl string `json:"dates"`
}

type Locations struct{
	Index []LocationItem `json:"index"`
}

type FullArtist struct{
	Artist
	DatesLocations map[string][]string
	Locations []LocationInfo
}

type SearchResult struct{
	Search string
	Artists []FullArtist
}

type ErrorPage struct{
	Code int
	Message string
}

type Geolocation struct{
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type LocationInfo struct {
	Name string
	Lat  string
	Lon  string
	Dates []string
}

type ArtistPageData struct {
    Artist
    DatesLocations map[string][]string
    Locations       template.JS
}

type ArtistsPageData struct{
	Artists []FullArtist
	PageNumbers []int
	PageNo int
	Query string
	Sort string
	TotalArtists int
}