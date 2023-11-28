package ec2

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Service interface {
	GetResponse(url string, os string) (*SpotPricing, error)
}
type SpotPricing struct {
	Vers   float64       `json:"vers"`
	Config PricingConfig `json:"config"`
}

type PricingConfig struct {
	Rate         string   `json:"rate"`
	ValueColumns []string `json:"valueColumns"`
	Currencies   []string `json:"currencies"`
	Regions      []Region `json:"regions"`
}

type Region struct {
	Region        string         `json:"region"`
	Footnotes     Footnotes      `json:"footnotes"`
	InstanceTypes []InstanceType `json:"instanceTypes"`
}

type Footnotes struct {
	NAMING_FAILED string `json:"*"`
}

type InstanceType struct {
	Type  string `json:"type"`
	Sizes []Size `json:"sizes"`
}

type Size struct {
	Size         string        `json:"size"`
	ValueColumns []ValueColumn `json:"valueColumns"`
}

type ValueColumn struct {
	Name   string `json:"name"`
	Prices Prices `json:"prices"`
}

type Prices struct {
	USD string `json:"USD"`
}

type Spot struct{}

func (s *Spot) GetResponse(url string, os string) (*SpotPricing, error) {

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

	if os == "" {
		return &data, nil
	}
	fmt.Println(os)
	excludeMswinEntries(&data, os)

	//fmt.Println(data)
	return &data, nil

}
func excludeMswinEntries(spotPricing *SpotPricing, os string) {
	for i, region := range spotPricing.Config.Regions {
		for j, instanceType := range region.InstanceTypes {
			for k, size := range instanceType.Sizes {
				if os == "linux" {
					filteredValueColumns := filterValueColumns(size.ValueColumns, "mswin")
					spotPricing.Config.Regions[i].InstanceTypes[j].Sizes[k].ValueColumns = filteredValueColumns

				} else if os == "mswin" {
					filteredValueColumns := filterValueColumns(size.ValueColumns, "linux")
					spotPricing.Config.Regions[i].InstanceTypes[j].Sizes[k].ValueColumns = filteredValueColumns

				}
			}
		}
	}
}

func filterValueColumns(columns []ValueColumn, excludeName string) []ValueColumn {
	filtered := []ValueColumn{}
	for _, col := range columns {
		if col.Name != excludeName {
			filtered = append(filtered, col)
		}
	}
	return filtered
}
func GetSpotPricing(service Service, region string, os string) (*[]Region, error) {
	body, err := service.GetResponse("https://website.spot.ec2.aws.a2z.com/spot.json", os)
	if err != nil {
		return nil, err
	}

	if region == "" {
		return &body.Config.Regions, nil
	}

	for _, r := range body.Config.Regions {

		if r.Region == region {
			return &[]Region{r}, nil
		}

	}
	return nil, fmt.Errorf("region %s not found", region)

}
