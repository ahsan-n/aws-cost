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
	expectedResponse := response()

	t.Run("SuccessfulFetch", func(t *testing.T) {
		mockService.On("GetResponse", url).Return(expectedResponse, nil)

		pricing, err := ec2.GetSpotPricing(mockService)
		if err != nil {
			return
		}

		assert.Equal(t, 2.00, pricing)

		mockService.AssertExpectations(t)
	})

}

func response() *ec2.SpotPricing {
	expectedResponse := &ec2.SpotPricing{
		Vers: 1.0,
		Config: ec2.PricingConfig{
			Rate:         "perhour",
			ValueColumns: []string{"value1", "value2"},
			Currencies:   []string{"USD", "EUR"},
			Regions: []ec2.Region{
				{
					Region: "ap-southeast-1",
					Footnotes: ec2.Footnotes{
						NAMING_FAILED: "example footnote",
					},
					InstanceTypes: []ec2.InstanceType{
						{
							Type: "generalCurrentGen",
							Sizes: []ec2.Size{
								{
									Size: "t3a.micro",
									ValueColumns: []ec2.ValueColumn{
										{
											Name: "linux",
											Prices: ec2.Prices{
												USD: "0.025",
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Region: "ap-southeast-2",
					Footnotes: ec2.Footnotes{
						NAMING_FAILED: "example footnote",
					},
					InstanceTypes: []ec2.InstanceType{
						{
							Type: "generalCurrentGen",
							Sizes: []ec2.Size{
								{
									Size: "t2.micro",
									ValueColumns: []ec2.ValueColumn{
										{
											Name: "linux",
											Prices: ec2.Prices{
												USD: "0.02125",
											},
										},
									},
								},
							},
						},
					},
				},
				// Add more regions as needed
			},
		},
	}
	return expectedResponse
}
