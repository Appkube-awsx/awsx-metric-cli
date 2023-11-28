package command

import (
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-metric/controller"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var AwsxCloudWatchMetricsCmd = &cobra.Command{
	Use:   "getAwsCloudWatchMetrics",
	Short: "getAwsCloudWatchMetrics command gets cloudwatch metrics data",
	Long:  `getAwsCloudWatchMetrics command gets cloudwatch metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

		authFlag, clientAuth, err := authenticate.CommandAuth(cmd)
		if err != nil {
			log.Println("Error during authentication: %v", err)
			cmd.Help()
			return
		}
		if authFlag {
			// Retrieve JSON input from command-line flag
			query, err := cmd.Flags().GetString("query")
			if err != nil {
				log.Println("Error retrieving JSON input: %v", err)
				cmd.Help()
				return
			}
			// Call GetMetricData with clientAuth, JSON input, and dimensions
			if err := controller.GetMetricData(*clientAuth, query); err != nil {
				log.Println("Error getting metric data: %v", err)
			}
		}
	},
}

func Execute() {
	if err := AwsxCloudWatchMetricsCmd.Execute(); err != nil {
		log.Println("error executing command: %v", err)
	}
}

func init() {
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("cloudElementId", "", "cloud element id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("cloudElementApiUrl", "", "cloud element api")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("query", "", "aws cloudwatch metric queries")

	// Use ParseFlags instead of PersistentFlags().Parse
	if err := AwsxCloudWatchMetricsCmd.ParseFlags(os.Args[1:]); err != nil {
		log.Println("error parsing flags: %v", err)
	}
}
