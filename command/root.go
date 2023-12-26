package command

import (
	"encoding/json"
	"fmt"
	ec2 "github.com/Appkube-awsx/awsx-metric-cli/EC2/CPUUtilization"
	"github.com/Appkube-awsx/awsx-metric-cli/auth"
	"github.com/Appkube-awsx/awsx-metric-cli/command/encryptdecrypt"
	"github.com/spf13/cobra"
	"log"
)

type CloudWatchQuery struct {
	MetricName string `json:"MetricName"`
}

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
			cloudWatchQueries, err := cmd.PersistentFlags().GetString("cloudWatchQueries")
			if err != nil {
				log.Println("Error retrieving JSON input: %v", err)
				cmd.Help()
				return
			}
			if cloudWatchQueries != "" {
				var queries []struct {
					RefID string `json:"RefID"`
					Query []struct {
						Namespace  string `json:"Namespace"`
						MetricName string `json:"MetricName"`
						Period     int    `json:"Period"`
						Stat       string `json:"Stat"`
						Dimensions []struct {
							Name  string `json:"Name"`
							Value string `json:"Value"`
						} `json:"Dimensions"`
					} `json:"Query"`
				}
				if err := json.Unmarshal([]byte(cloudWatchQueries), &queries); err != nil {
					fmt.Println("Error parsing JSON:", err)
					return
				}
				var metricNames []CloudWatchQuery
				var metricNamesMap = make(map[string]bool)
				for _, query := range queries {
					for _, q := range query.Query {
						metricName := q.MetricName
						if _, exists := metricNamesMap[metricName]; !exists {
							metricNamesMap[metricName] = true
							metricNames = append(metricNames, CloudWatchQuery{MetricName: metricName})
						}
					}
				}
				metricNamesJSON, err := json.Marshal(metricNames)
				if err != nil {
					log.Printf("Error marshaling MetricName: %v", err)
					return
				}

				fmt.Println(string(metricNamesJSON))
				for _, metric := range metricNames {
					fmt.Println("MetricName:", metric.MetricName)
					switch metric.MetricName {
					case "CPUUtilization":
						fmt.Println("Processing CPUUtilization")
						res, err := ec2.GetMetricDataCPU(clientAuth, cloudWatchQueries)
						if err != nil {
							log.Printf("Error getting CPUUtilization metric data: %v", err)
							return
						}
						log.Println(res)
						fmt.Println("Processing complete for CPUUtilization")
					default:
						log.Printf("Unsupported metric name: %s", metric.MetricName)
					}
				}

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
