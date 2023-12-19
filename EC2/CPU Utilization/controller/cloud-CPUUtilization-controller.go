package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"strings"
	"time"
)

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

const (
	defaultStartTimeOffset = -1
	defaultTimeRange       = 1 // in hours
)

func getMetricDataInternal(clientAuth *client.Auth, cloudWatchQueries string, numQueries int) (*cloudwatch.GetMetricDataOutput, error) {
	var outerQueries []OuterQuery
	err := json.Unmarshal([]byte(cloudWatchQueries), &outerQueries)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON input: %v", err)
	}

	queries := make([]*cloudwatch.MetricDataQuery, numQueries)
	for i, outerQueryInput := range outerQueries {
		for _, queryInput := range outerQueryInput.Query {
			if queryInput.Dimensions == nil {
				queryInput.Dimensions = make([]Dimension, 0)
			}

			query := &cloudwatch.MetricDataQuery{
				Id:         aws.String(strings.ToLower(outerQueryInput.RefID)),
				ReturnData: aws.Bool(true),
				MetricStat: &cloudwatch.MetricStat{
					Metric: &cloudwatch.Metric{
						Namespace:  aws.String(queryInput.Namespace),
						MetricName: aws.String(queryInput.MetricName),
						Dimensions: buildDimensions(queryInput.Dimensions),
					},
					Period: aws.Int64(queryInput.Period),
					Stat:   aws.String(queryInput.Stat),
				},
			}
			queries[i] = query
		}
	}

	cloudWatchClient := client.GetClient(*clientAuth, client.CLOUDWATCH).(*cloudwatch.CloudWatch)

	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: queries,
		StartTime:         aws.Time(time.Now().Add(time.Duration(defaultStartTimeOffset) * time.Hour)),
		EndTime:           aws.Time(time.Now()),
	}

	result, err := cloudWatchClient.GetMetricData(input)
	if err != nil {
		return nil, fmt.Errorf("error making request to CloudWatch: %v", err)
	}
	// Extract the values
	values := result.MetricDataResults
	// Prepare output in JSON format
	output := map[string]interface{}{}
	if len(values) > 0 {
		output["CurrentCPUUtilization"] = fmt.Sprintf("%.2f%%", *values[0].Values[0]) // First query is for current CPU utilization

		if len(values) > 1 {
			output["MaxCPUUtilization24h"] = fmt.Sprintf("%.2f%%", *values[1].Values[0]) // Second query is for maximum CPU utilization (past 24 hours)
		}

		if len(values) > 2 {
			output["AvgCPUUtilization24h"] = fmt.Sprintf("%.2f%%", *values[2].Values[0]) // Third query is for average CPU utilization (past 24 hours)
		}
	}

	// Convert output to JSON with new lines
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return nil, err
	}

	// Print JSON output
	fmt.Println(string(jsonOutput))

	return result, nil
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

func GetMetricData(clientAuth *client.Auth, cloudWatchQueries string) (*cloudwatch.GetMetricDataOutput, error) {
	var outerQueries []OuterQuery
	err := json.Unmarshal([]byte(cloudWatchQueries), &outerQueries)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON input: %v", err)
	}

	return getMetricDataInternal(clientAuth, cloudWatchQueries, len(outerQueries))
}
