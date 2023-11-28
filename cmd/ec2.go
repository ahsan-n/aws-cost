/*
Copyright © 2023 Ahsan Naseem

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/ahsan-n/aws-cost/internal/ec2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		limit, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		region, err := cmd.Flags().GetString("region")

		if err != nil {
			return err
		}

		os, err := cmd.Flags().GetString("operating-system")
		if err != nil {
			return err
		}
		return pricing(limit, region, os)
	},
}

func beautify(region *[]ec2.Region) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAutoIndex(true)
	t.AppendHeader(table.Row{"size", "linux/windows", "price"})
	for _, r := range *region {
		for _, i := range r.InstanceTypes {
			for _, s := range i.Sizes {
				for _, v := range s.ValueColumns {
					t.AppendRow(table.Row{s.Size, v.Name, v.Prices.USD})

				}
			}
		}

	}

	t.AppendSeparator()
	t.Render()
}
func pricing(limit bool, region string, os string) error {

	s := &ec2.Spot{}
	sp, err := ec2.GetSpotPricing(s, region, os)
	beautify(sp)
	//fmt.Println(spotPricing)
	if err != nil {
		return err
	}

	return nil

}

func init() {
	ec2Cmd.PersistentFlags().BoolP("all", "a", false, "get all ec2 from all regions (expects a boolean value)")
	ec2Cmd.PersistentFlags().StringP("region", "r", "ap-southeast-1", "aws region name")
	ec2Cmd.PersistentFlags().StringP("operating-system", "o", "linux", "include os")
	getCmd.AddCommand(ec2Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
