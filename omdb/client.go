package omdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey  string
	baseURL string
}

type Response struct {
	Title        string   `json:"Title"`
	Year         string   `json:"Year"`
	Rated        string   `json:"Rated"`
	Released     string   `json:"Released"`
	Runtime      string   `json:"Runtime"`
	Genre        string   `json:"Genre"`
	Director     string   `json:"Director"`
	Writer       string   `json:"Writer"`
	Actors       string   `json:"Actors"`
	Plot         string   `json:"Plot"`
	Language     string   `json:"Language"`
	Country      string   `json:"Country"`
	Awards       string   `json:"Awards"`
	Poster       string   `json:"Poster"`
	Ratings      []Rating `json:"Ratings"`
	Metascore    string   `json:"Metascore"`
	ImdbRating   string   `json:"imdbRating"`
	ImdbVotes    string   `json:"imdbVotes"`
	ImdbID       string   `json:"imdbID"`
	Type         string   `json:"Type"`
	DVD          string   `json:"DVD"`
	BoxOffice    string   `json:"BoxOffice"`
	Production   string   `json:"Production"`
	Website      string   `json:"Website"`
	Response     string   `json:"Response"`
	Error        string   `json:"Error"`
	TotalSeasons string   `json:"totalSeasons,omitempty"`
}

type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "http://www.omdbapi.com/",
	}
}

func (client *Client) GetById(id string) (*Response, error) {
	url := client.baseURL + "?i=" + id + "&apiKey=" + client.apiKey
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if res.Response != "True" {
		return nil, fmt.Errorf("error: %s", res.Error)
	}

	return &res, nil
}
