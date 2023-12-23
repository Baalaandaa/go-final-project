package location_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LatLngLiteral struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius float64 `json:"radius"`
}

type Driver struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	DriverId string  `json:"id"`
}

type LocationServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func New(baseUrl string) LocationServiceClient {
	return LocationServiceClient{
		BaseURL: baseUrl,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c LocationServiceClient) GetDriverLocations(lat, lng, radius float64) ([]Driver, error) {
	params := LatLngLiteral{
		Lat:    lat,
		Lng:    lng,
		Radius: radius,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request params: %w", err)
	}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/location/v1/drivers", c.BaseURL),
		bytes.NewBuffer(paramsBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request to location service: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making get request to location service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body from location service: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("location service returned non-OK status: %s", resp.Status)
	}

	var result []Driver
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling driver locations: %w", err)
	}

	return result, nil
}
