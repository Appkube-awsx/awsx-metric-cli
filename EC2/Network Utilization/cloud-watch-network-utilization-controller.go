package Network_Utilization

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-metric-cli/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"strings"
	"time"
)

// Existing Dimension and MetricQueryInput structures...
type Dimension struct {
	Name  string
	Value string
}
type TimeRange struct {
	From     string `json:"From"`
	To       string `json:"To"`
	TimeZone string `json:"TimeZone"`
}
type InnerQuery struct {
	Namespace  string      `json:"Namespace"`
	MetricName string      `json:"MetricName"`
	Period     int64       `json:"Period"`
	Stat       string      `json:"Stat"`
	Dimensions []Dimension `json:"Dimensions"`
}
type OuterQuery struct {
	RefID        string       `json:"RefID"`
	MaxDataPoint int          `json:"MaxDataPoint"`
	Interval     int          `json:"Interval"`
	TimeRange    TimeRange    `json:"TimeRange"`
	Query        []InnerQuery `json:"Query"`
}

func GetNerworkUtilizationMetricData(clientAuth *client.Auth, cloudWatchQueries string) (*cloudwatch.GetMetricDataOutput, error) {
	var outerQuery []OuterQuery
	err := json.Unmarshal([]byte(cloudWatchQueries), &outerQuery)
	if err != nil {
		fmt.Println("Error parsing JSON input:", err)
		return nil, err
	}
	// Create the metric queries dynamically
	queries := make([]*cloudwatch.MetricDataQuery, 0)
	for i, outerQueryInput := range outerQuery {
		for _, queryInput := range outerQueryInput.Query {
			if queryInput.Dimensions == nil {
				queryInput.Dimensions = make([]Dimension, 0)
			}
			query := &cloudwatch.MetricDataQuery{
				Id:         aws.String(strings.ToLower(fmt.Sprintf("%s_%d", outerQueryInput.RefID, i))),
				ReturnData: aws.Bool(true),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String(queryInput.Namespace),
						MetricName: aws.String(queryInput.MetricName),
						Dimensions: buildDimensions(queryInput.Dimensions),
					},
					Period: aws.Int64(queryInput.Period),
					Stat:   aws.String(queryInput.Stat),
					Unit:   aws.String("Bytes"),
				},
			}
			queries = append(queries, query)
		}
	}
	cloudWatchClient := client.GetClient(*clientAuth, client.CLOUDWATCH).(*cloudwatch.CloudWatch)
	// Specify the request input with multiple queries
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: queries,
		StartTime:         aws.Time(time.Now().Add(time.Duration(-1) * time.Hour)),
		EndTime:           aws.Time(time.Now()),
	}
	// Make the request to CloudWatch
	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	var inbound, outbound float64
	for _, metricDataResult := range result.MetricDataResults {
		switch *metricDataResult.Label {
		case "NetworkIn":
			processMetricData("NetworkIn", metricDataResult, &inbound)
		case "NetworkOut":
			processMetricData("NetworkOut", metricDataResult, &outbound)
		default:
			fmt.Printf("Unknown metric label: %v\n", *metricDataResult.Label)
		}
	}

	inboundBytes := inbound
	outboundBytes := outbound
	// Output the result
	//fmt.Printf("{ Inbound: %v, Outbound: %v, DataTransfer: %v }\n", inbound, outbound, inbound+outbound)

	fmt.Printf("{ Inbound: %.2f Bytes, Outbound: %.2f Bytes, DataTransfer: %.2f Bytes }\n", inboundBytes, outboundBytes, inboundBytes+outboundBytes)
	//for _, metricDataResult := range result.MetricDataResults {
	// // Assuming you are interested in NetworkIn metric
	// if len(metricDataResult.Values) > 0 {
	//    inbound := *metricDataResult.Values[0]
	//    outbound := 0.0 // Assume outbound is not available in the response
	//    dataTransfer := inbound + outbound
	//
	//    fmt.Printf("{ Inbound: %v, Outbound: %v, DataTransfer: %v }\n", inbound, outbound, dataTransfer)
	// } else {
	//    fmt.Println("No data available for NetworkIn")
	// }
	// if *metricDataResult.Label == "NetworkIn" {
	//    for i, timestamp := range metricDataResult.Timestamps {
	//       fmt.Printf("Data for NetworkIn at Timestamp %v: %f\n", *timestamp, *metricDataResult.Values[i])
	//    }
	//
	// }
	//}
	return result, nil
}
func processMetricData(metricLabel string, metricDataResult *cloudwatch.MetricDataResult, value *float64) {
	fmt.Printf("Data for %s at Timestamps:\n", metricLabel)
	for i, timestamp := range metricDataResult.Timestamps {
		fmt.Printf("  Timestamp %v: %f\n", *timestamp, *metricDataResult.Values[i])
	}
	if len(metricDataResult.Values) > 0 {
		// Corrected code: Dereference the pointer
		*value = *metricDataResult.Values[0]
		fmt.Printf("{ %s: %v }\n", metricLabel, *value)
	} else {
		fmt.Printf("No data available for %s\n", metricLabel)
	}
}
func buildDimensions(dimensions []Dimension) []*cloudwatch.Dimension {
	var cloudWatchDimensions []*cloudwatch.Dimension
	for _, d := range dimensions {
		dimension := &cloudwatch.Dimension{
			Name:  aws.String(d.Name),
			Value: aws.String(d.Value),
		}
		cloudWatchDimensions = append(cloudWatchDimensions, dimension)
	}
	return cloudWatchDimensions
}
