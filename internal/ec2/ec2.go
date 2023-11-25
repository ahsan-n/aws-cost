package ec2

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Service interface {
	GetResponse(url string) (*SpotPricing, error)
}
type SpotPricing struct {
	Vers   float64 `json:"vers"`
	Config struct {
		Rate         string   `json:"rate"`
		ValueColumns []string `json:"valueColumns"`
		Currencies   []string `json:"currencies"`
		Regions      []struct {
			Region    string `json:"region"`
			Footnotes struct {
				NAMING_FAILED string `json:"*"`
			} `json:"footnotes"`
			InstanceTypes []struct {
				Type  string `json:"type"`
				Sizes []struct {
					Size         string `json:"size"`
					ValueColumns []struct {
						Name   string `json:"name"`
						Prices struct {
							USD string `json:"USD"`
						} `json:"prices"`
					} `json:"valueColumns"`
				} `json:"sizes"`
			} `json:"instanceTypes"`
		} `json:"regions"`
	} `json:"config"`
}

type Spot struct{}

func (s *Spot) GetResponse(url string) (*SpotPricing, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)

	}
	defer resp.Body.Close() // Don't forget to close the response body

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data SpotPricing
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
	return &data, nil

}

func GetSpotPricing(service Service) (float64, error) {
	body, err := service.GetResponse("https://website.spot.ec2.aws.a2z.com/spot.json")
	if err != nil {
		return 0, err
	}

	return body.Vers + 1, nil

}
