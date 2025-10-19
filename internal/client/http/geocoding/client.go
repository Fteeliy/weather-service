package geocoding

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *client {
	return &client{
		httpClient: httpClient,
	}
}

func (c *client) GetCoords(city string) (Response, error) {
	res, err := c.httpClient.Get(
		fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json", city),
	)

	if err != nil {
		log.Printf("Error making request to geocoding API: %v", err)
		return Response{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("geocoding api returned non-200 status: %d", res.StatusCode)
	}

	var geoResponse struct {
		Results []Response `json:"results"`
	}

	err = json.NewDecoder(res.Body).Decode(&geoResponse)
	if err != nil {
		return Response{}, err
	}

	if len(geoResponse.Results) == 0 {
		return Response{}, fmt.Errorf("no results found for city: %s", city)
	}

	return geoResponse.Results[0], nil
}
