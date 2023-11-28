package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
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

type QueryInput struct {
	RefID      string      `json:"RefID"`
	Namespace  string      `json:"Namespace"`
	MetricName string      `json:"MetricName"`
	Period     int64       `json:"Period"`
	Stat       string      `json:"Stat"`
	Dimensions []Dimension `json:"Dimensions"`
}

type MetricQueryInputV2 struct {
	RefID        string       `json:"RefID"`
	MaxDataPoint int          `json:"MaxDataPoint"`
	Interval     int          `json:"Interval"`
	TimeRange    TimeRange    `json:"TimeRange"`
	Query        []QueryInput `json:"Query"`
}

func GetMetricData(clientAuth client.Auth, jsonInput string) *cloudwatch.GetMetricDataOutput {
	var metricQueriesV2 []MetricQueryInputV2
	err := json.Unmarshal([]byte(jsonInput), &metricQueriesV2)
	if err != nil {
		fmt.Println("Error parsing JSON input:", err)
		return nil
	}

	// Create the metric queries dynamically
	queries := make([]*cloudwatch.MetricDataQuery, len(metricQueriesV2))
	for i, queryInputV2 := range metricQueriesV2 {
		for _, queryInput := range queryInputV2.Query {
			if queryInput.Dimensions == nil {
				queryInput.Dimensions = make([]Dimension, 0)
			}

			query := &cloudwatch.MetricDataQuery{
				Id:         aws.String(fmt.Sprintf(queryInput.RefID)),
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

	cloudWatchClient := client.GetClient(clientAuth, client.CLOUDWATCH).(*cloudwatch.CloudWatch)

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
		return nil
	}

	// Process the result
	for _, metricDataResult := range result.MetricDataResults {
		for i, timestamp := range metricDataResult.Timestamps {
			fmt.Printf("Data for Metric at Timestamp %v: %f\n", *timestamp, *metricDataResult.Values[i])
		}
	}

	return result
}

// Existing buildDimensions function...

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
