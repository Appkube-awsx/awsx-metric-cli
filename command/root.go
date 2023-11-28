/*
Copyright © 2023 Manoj Sharma manoj.sharma@synectiks.com
*/
package command

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-metric/apiclient"
	"github.com/Appkube-awsx/awsx-metric/controller"
	"github.com/Appkube-awsx/awsx-metric/models"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var AwsxCloudWatchMetricsCmd = &cobra.Command{
	Use:   "getAwsCloudWatchMetrics",
	Short: "getAwsCloudWatchMetrics command gets cloudwatch metrics data",
	Long:  `getAwsCloudWatchMetrics command gets cloudwatch metrics data`,

	Run: func(cmd *cobra.Command, args []string) {

		authFlag, clientAuth, err := authenticate.CommandAuth(cmd)
		if err != nil {
			log.Fatalf("Error during authentication: %v", err)
			cmd.Help()
			return
		}
		if authFlag {
			// Retrieve JSON input from command-line flag
			query, err := cmd.Flags().GetString("query")
			if err != nil {
				log.Fatalf("Error retrieving JSON input: %v", err)
				cmd.Help()
				return
			}
			// Call GetMetricData with clientAuth, JSON input, and dimensions
			if err := controller.GetMetricData(*clientAuth, query); err != nil {
				log.Fatalf("Error getting metric data: %v", err)
			}
		}
	},
}

func dynamicCmd(name string, desc string) *cobra.Command {
	return &cobra.Command{}
}
func Execute() {
	if err := AwsxCloudWatchMetricsCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
func getCmdbData(cloudElementId string, cmdbCloudElementApi string) (*models.CmdbCloudElementResponse, int, time.Duration, error) {
	url := cmdbCloudElementApi + cloudElementId
	fmt.Println("Request to get landing-zone, CMDB landing-zone url: " + url)
	cmdbResp, cmdbStatusCode, duration, err := apiclient.ProcessApiCall(url, nil)
	if err != nil {
		fmt.Println("CMDB landing-zone api call failed. Error: ", "error", err.Error())
		return nil, cmdbStatusCode, duration, err
	}

	var cmdbSlice []interface{}

	switch v := cmdbResp.(type) {
	case []interface{}:
		cmdbSlice = v
	default:
		// Handle other types or report an error
		log.Fatalf("Unexpected type for cmdbResp: %T", v)
		return nil, cmdbStatusCode, duration, fmt.Errorf("unexpected type for cmdbResp: %T", v)
	}

	if len(cmdbSlice) > 0 {
		cmdbByte, err := json.Marshal(cmdbSlice[0]) // assuming the first element contains the desired data
		if err != nil {
			fmt.Println("Error in marshalling cmdb cloud-element response", "error", err.Error())
			return nil, cmdbStatusCode, duration, err
		}

		var out models.CmdbCloudElementResponse
		err = json.Unmarshal(cmdbByte, &out)
		if err != nil {
			fmt.Println("Error in parsing cmdb cloud-element response", "error", err.Error())
			return nil, cmdbStatusCode, duration, err
		}

		return &out, cmdbStatusCode, duration, nil
	}

	return nil, cmdbStatusCode, duration, fmt.Errorf("No data found in CMDB response")
}

func getAwsCredentials(landingZoneId int64, cmdbLandingCloudCredsApiUrl string) (*models.AwsCredential, int, time.Duration, error) {
	fmt.Println("Request to get AWS credentials")
	url := cmdbLandingCloudCredsApiUrl + strconv.FormatInt(landingZoneId, 10)
	fmt.Println("Request to get cloud credentials, CMDB cloud-creds URL: " + url)

	cloudCredResp, cmdbStatusCode, duration, err := apiclient.ProcessApiCall(url, nil)
	if err != nil {
		fmt.Println("CMDB cloud-creds API call failed. Error:", err.Error())
		return nil, cmdbStatusCode, duration, err
	}

	// Check the type of cloudCredResp
	switch v := cloudCredResp.(type) {
	case map[string]interface{}:
		// Convert the map to JSON
		cmdbByte, err := json.Marshal(v)
		if err != nil {
			fmt.Println("Error in marshalling cloud-creds response", "error", err.Error())
			return nil, cmdbStatusCode, duration, err
		}

		// Unmarshal the JSON into the AwsCredential struct
		var out models.AwsCredential
		err = json.Unmarshal(cmdbByte, &out)
		if err != nil {
			fmt.Println("Error in parsing cloud-creds response", "error", err.Error())
			return nil, cmdbStatusCode, duration, err
		}

		return &out, cmdbStatusCode, duration, nil
	default:
		// Handle other types or report an error
		return nil, cmdbStatusCode, duration, fmt.Errorf("unexpected type for cloud-creds response: %T", v)
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
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("query", "", "dynamic metric queries")

	// Use ParseFlags instead of PersistentFlags().Parse
	if err := AwsxCloudWatchMetricsCmd.ParseFlags(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
}
