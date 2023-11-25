package ec2_test

import (
	"github.com/ahsan-n/aws-cost/internal/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetResponse(url string) (*ec2.SpotPricing, error) {

	args := m.Called(url)
	return args.Get(0).(*ec2.SpotPricing), args.Error(1)
}

func TestGetSpotPricing(t *testing.T) {
	url := "https://website.spot.ec2.aws.a2z.com/spot.json"
	mockService := new(MockService)
	expectedResponse := &ec2.SpotPricing{
		Vers: 0,
		Config: struct {
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
		}{Rate: "perhr"},
	}

	t.Run("SuccessfulFetch", func(t *testing.T) {
		mockService.On("GetResponse", url).Return(expectedResponse, nil)

		pricing, err := ec2.GetSpotPricing(mockService)
		if err != nil {
			return
		}

		assert.Equal(t, 1.00, pricing)

		mockService.AssertExpectations(t)
	})
}
