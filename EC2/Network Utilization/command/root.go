package command

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-network-utilization/controller"
	"github.com/spf13/cobra"
)

var AwsxNetworkUtilizationCmd = &cobra.Command{
	Use:   "NetworkUtilization",
	Short: "NetworkUtilization command gets cloudwatch metrics data",
	Long:  `NetworkUtilization command gets cloudwatch metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

		authFlag, clientAuth, err := authenticate.CommandAuth(cmd)
		if err != nil {
			log.Println("Error during authentication: %v", err)
			cmd.Help()
			return
		}
		if authFlag {
			// Retrieve JSON input from command-line flag

			cloudWatchQueries, err := cmd.PersistentFlags().GetString("cloudWatchQueries")
			if err != nil {
				log.Println("Error retrieving JSON input: %v", err)
				cmd.Help()
				return
			}
			if cloudWatchQueries == "" {
				log.Println("cloud-watch query not provided. program exit")
				cmd.Help()
				return
			}
			// Call GetMetricData with clientAuth, JSON input, and dimensions
			res, err := controller.GetMetricData(clientAuth, cloudWatchQueries)

			if err != nil {
				log.Println("Error getting metric data: %v", err)
				return
			}
			log.Println(res)
		}
	},
}

func Execute() {
	if err := AwsxNetworkUtilizationCmd.Execute(); err != nil {
		log.Println("error executing command: %v", err)
	}
}

func init() {
	AwsxNetworkUtilizationCmd.PersistentFlags().String("cloudElementId", "", "cloud element id")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("cloudElementApiUrl", "", "cloud element api")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxNetworkUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")

	// Use ParseFlags instead of PersistentFlags().Parse
	if err := AwsxNetworkUtilizationCmd.ParseFlags(os.Args[1:]); err != nil {
		log.Println("error parsing flags: %v", err)
	}
}
