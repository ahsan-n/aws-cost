package ec2_test

import (
	"fmt"
	"github.com/ahsan-n/aws-cost/internal/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetResponse(url string, os string) (*ec2.SpotPricing, error) {

	args := m.Called(url)
	return args.Get(0).(*ec2.SpotPricing), args.Error(1)
}

func TestGetSpotPricing(t *testing.T) {
	url := "https://website.spot.ec2.aws.a2z.com/spot.json"
	mockService := new(MockService)
	expectedResponse := response("")

	t.Run("ValidateRegion", func(t *testing.T) {
		mockService.On("GetResponse", url).Return(expectedResponse, nil)

		actual, err := ec2.GetSpotPricing(mockService, "ap-southeast-1", "")
		if err != nil {
			return
		}

		assert.Equal(t, 1, len(*actual))
		assert.Equal(t, "ap-southeast-1", (*actual)[0].Region)

		mockService.AssertExpectations(t)
	})

	t.Run("AllRegions", func(t *testing.T) {
		mockService.On("GetResponse", url).Return(expectedResponse, nil)

		actual, err := ec2.GetSpotPricing(mockService, "", "")

		if err != nil {
			return
		}

		assert.Equal(t, 2, len(*actual))

		mockService.AssertExpectations(t)
	})

	t.Run("OSFiltering", func(t *testing.T) {
		mockService.On("GetResponse", url).Return(expectedResponse, nil)

		regions, err := ec2.GetSpotPricing(mockService, "ap-southeast-1", "linux")
		if err != nil {
			assert.Fail(t, "error should be nil", err)
		}
		for _, r := range *regions {
			for _, v := range r.InstanceTypes[0].Sizes[0].ValueColumns {
				if v.Name == "mswin" {
					assert.Fail(t, "unable to filter windows machines")
				}
				fmt.Println(r.InstanceTypes[0].Sizes[0].ValueColumns)
			}
		}

		mockService.AssertExpectations(t)
	})
	t.Run("RegionNotFound", func(t *testing.T) {
		mockService.On("GetResponse", url).Return(expectedResponse, nil)

		_, err := ec2.GetSpotPricing(mockService, "eu-west-1", "")
		if err != nil {
			assert.Error(t, err)

			assert.Equal(t, fmt.Errorf("region eu-west-1 not found"), err)
		}

		mockService.AssertExpectations(t)
	})

}

func response(os string) *ec2.SpotPricing {
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
	if os == "linux" {
		expectedResponse = &ec2.SpotPricing{
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
											}, {
												Name: "mswin",
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
	}
	return expectedResponse
}
