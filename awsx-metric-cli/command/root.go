package command

import (
	"github.com/Appkube-awsx/awsx-metric-cli/auth"
	"github.com/Appkube-awsx/awsx-metric-cli/command/encryptdecrypt"
	"github.com/Appkube-awsx/awsx-metric-cli/controller"
	"github.com/spf13/cobra"
	"log"
)

var AwsxCloudWatchMetricsCmd = &cobra.Command{
	Use:   "getAwsCloudWatchMetrics",
	Short: "getAwsCloudWatchMetrics command gets cloudwatch metrics data",
	Long:  `getAwsCloudWatchMetrics command gets cloudwatch metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

		var authFlag, clientAuth, err = authenticate.CommandAuth(cmd)
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
	if err := AwsxCloudWatchMetricsCmd.Execute(); err != nil {
		log.Println("error executing command: %v", err)
	}
}

func init() {
	AwsxCloudWatchMetricsCmd.AddCommand(encryptdecrypt.EncryptDecrypt)
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
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")

}
