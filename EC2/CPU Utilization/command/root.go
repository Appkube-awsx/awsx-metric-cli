package command

import (
	"log"
	"os"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-cpu-utilization/controller"
	"github.com/spf13/cobra"
)

var AwsxCpuUtilizationCmd = &cobra.Command{
	Use:   "CPUUtilization",
	Short: "CPUUtilization command gets cloudwatch metrics data",
	Long:  `CPUUtilization command gets cloudwatch metrics data`,

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
	if err := AwsxCpuUtilizationCmd.Execute(); err != nil {
		log.Println("error executing command: %v", err)
	}
}

func init() {
	AwsxCpuUtilizationCmd.PersistentFlags().String("cloudElementId", "", "cloud element id")
	AwsxCpuUtilizationCmd.PersistentFlags().String("cloudElementApiUrl", "", "cloud element api")
	AwsxCpuUtilizationCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxCpuUtilizationCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxCpuUtilizationCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxCpuUtilizationCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxCpuUtilizationCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxCpuUtilizationCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxCpuUtilizationCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxCpuUtilizationCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxCpuUtilizationCmd.PersistentFlags().String("cloudWatchQueries", "", "aws cloudwatch metric queries")

	// Use ParseFlags instead of PersistentFlags().Parse
	if err := AwsxCpuUtilizationCmd.ParseFlags(os.Args[1:]); err != nil {
		log.Println("error parsing flags: %v", err)
	}
}
