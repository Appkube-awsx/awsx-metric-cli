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
	"net/http"
	"os"
	"time"
)

var AwsxCloudWatchMetricsCmd = &cobra.Command{
	Use:   "getAwsCloudWatchMetrics",
	Short: "getAwsCloudWatchMetrics command gets cloudwatch metrics data",
	Long:  `getAwsCloudWatchMetrics command gets cloudwatch metrics data`,

	Run: func(cmd *cobra.Command, args []string) {
		cloudElementId, err := cmd.Flags().GetString("cloudElementId")
		if err != nil {
			log.Fatalf("Error while getting cloudElementId command line parameter. Error %v", err)
			cmd.Help()
			return
		}
		if cloudElementId != "" {
			cmdbCloudElementApiUrl, err := cmd.Flags().GetString("cmdbCloudElementApiUrl")
			if err != nil {
				log.Fatalf("Error while getting cmdbCloudElementApiUrl command line parameter. Error %v", err)
				cmd.Help()
				return
			}
			if cmdbCloudElementApiUrl == "" {
				log.Fatalf("cmdb cloud-element url not provided")
				cmd.Help()
				return
			}

			cmdbCloudCredsApiUrl, err := cmd.Flags().GetString("cmdbCloudCredsApiUrl")
			if err != nil {
				log.Fatalf("Error while getting cmdbCloudCredsApiUrl command line parameter. Error %v", err)
				cmd.Help()
				return
			}
			if cmdbCloudCredsApiUrl == "" {
				log.Fatalf("cmdb cloud credential url not provided")
				cmd.Help()
				return
			}
			cmdbResp, cmdbStatusCode, _, err := getCmdbData(cloudElementId, cmdbCloudElementApiUrl)
			if err != nil {
				fmt.Println("error in cmdb cloud-element api call", "error", err.Error())
				return
			}
			if cmdbStatusCode >= http.StatusBadRequest {
				fmt.Println("CMDB cloud-element api call failed", "error", err.Error())
				return
			}
			vaultResp, vaultStatusCode, _, err := getAwsCredentials(cmdbResp.LandingzoneId, cmdbCloudCredsApiUrl)
			if err != nil {
				fmt.Println("error in cmdb cloud-creds api call", "error", err.Error())
				return
			}
			if vaultStatusCode >= http.StatusBadRequest {
				fmt.Println("CMDB cloud-creds api call failed", "error", err.Error())
				return
			}
			// create new cmd and set vaultResp acess key secret key
			//dCmd := dynamicCmd("zone", vaultResp.Region)
			dCmd := &cobra.Command{}

			dCmd.PersistentFlags().String("zone", vaultResp.Region, "aws region")
			dCmd.PersistentFlags().String("accessKey", vaultResp.AccessKey, "aws access key")
			dCmd.PersistentFlags().String("secretKey", vaultResp.SecretKey, "aws secret key")
			dCmd.PersistentFlags().String("crossAccountRoleArn", vaultResp.CrossAccountRoleArn, "aws cross account role arn")
			dCmd.PersistentFlags().String("externalId", vaultResp.ExternalId, "aws external id")

			authFlag, clientAuth, err := authenticate.CommandAuth(dCmd)
			if err != nil {
				log.Fatalf("Error during authentication: %v", err)
				cmd.Help()
				return
			}
			if authFlag {
				// Retrieve JSON input from command-line flag
				jsonInput, err := cmd.Flags().GetString("jsonInput")
				if err != nil {
					log.Fatalf("Error retrieving JSON input: %v", err)
					cmd.Help()
					return
				}
				// Call GetMetricData with clientAuth, JSON input, and dimensions
				if err := controller.GetMetricData(*clientAuth, jsonInput); err != nil {
					log.Fatalf("Error getting metric data: %v", err)
				}
			} else {
				cmd.Help()
				return
			}

		} else {
			authFlag, clientAuth, err := authenticate.CommandAuth(cmd)
			if err != nil {
				log.Fatalf("Error during authentication: %v", err)
				cmd.Help()
				return
			}
			if authFlag {
				// Retrieve JSON input from command-line flag
				jsonInput, err := cmd.Flags().GetString("jsonInput")
				if err != nil {
					log.Fatalf("Error retrieving JSON input: %v", err)
					cmd.Help()
					return
				}
				// Call GetMetricData with clientAuth, JSON input, and dimensions
				if err := controller.GetMetricData(*clientAuth, jsonInput); err != nil {
					log.Fatalf("Error getting metric data: %v", err)
				}
			} else {
				cmd.Help()
				return
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
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("vaultUrl", "", "vault end point")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("vaultToken", "", "vault token")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("accountId", "", "aws account number")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("zone", "", "aws region")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("accessKey", "", "aws access key")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("secretKey", "", "aws secret key")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("crossAccountRoleArn", "", "aws cross account role arn")
	AwsxCloudWatchMetricsCmd.PersistentFlags().String("externalId", "", "aws external id")
	AwsxCloudWatchMetricsCmd.Flags().String("jsonInput", "", "JSON input for dynamic metric queries")
	AwsxCloudWatchMetricsCmd.Flags().String("cloudElementId", "", "cloud element id")
	AwsxCloudWatchMetricsCmd.Flags().String("cmdbCloudElementApiUrl", "", "cmdb cloud element api")
	AwsxCloudWatchMetricsCmd.Flags().String("cmdbCloudCredsApiUrl", "", "cloud cloud credential api")

	// Use ParseFlags instead of PersistentFlags().Parse
	if err := AwsxCloudWatchMetricsCmd.ParseFlags(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
}
