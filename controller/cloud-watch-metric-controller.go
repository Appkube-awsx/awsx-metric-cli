package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/client"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"net/http"
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

type Query struct {
	ElementType string     `json:"elementType,omitempty"`
	ElementId   int64      `json:"elementId,omitempty"`
	URL         string     `json:"url"`
	URLOptions  URLOptions `json:"url_options"`
}
type URLOptions struct {
	Method           string                  `json:"method"` // 'GET' | 'POST'
	Params           []URLOptionKeyValuePair `json:"params"`
	Headers          []URLOptionKeyValuePair `json:"headers"`
	Body             string                  `json:"data"`
	BodyType         string                  `json:"body_type"`
	BodyContentType  string                  `json:"body_content_type"`
	BodyForm         []URLOptionKeyValuePair `json:"body_form"`
	BodyGraphQLQuery string                  `json:"body_graphql_query"`
	// BodyGraphQLVariables string           `json:"body_graphql_variables"`
}
type URLOptionKeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Client struct {
	HttpClient *http.Client
	IsMock     bool
}

// GetMetricData retrieves metric data from AWS CloudWatch based on the provided JSON input.
func GetMetricData(clientAuth client.Auth, jsonInput string) *cloudwatch.GetMetricDataOutput {
	//cloudElementId, err := getCloudElementIdFromDB(cloudElementId)
	//if err != nil {
	//	fmt.Println("Error fetching cloudElementId from the database:", err)
	//	return nil
	//}

	var metricQueriesV2 []MetricQueryInputV2
	// Use a different variable name for err in the second declaration
	json.Unmarshal([]byte(jsonInput), &metricQueriesV2)
	//if err != nil {
	//	fmt.Println("Error parsing JSON input:", err)
	//	return nil
	//}

	// Create the metric queries dynamically
	queries := make([]*cloudwatch.MetricDataQuery, 0)
	for i, queryInputV2 := range metricQueriesV2 {
		for _, queryInput := range queryInputV2.Query {
			// Ensure Dimensions is not nil
			if queryInput.Dimensions == nil {
				queryInput.Dimensions = make([]Dimension, 0)
			}

			// Construct a MetricDataQuery
			query := &cloudwatch.MetricDataQuery{
				Id:         aws.String(fmt.Sprintf("q%d", i+1)),
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
			queries = append(queries, query)
		}
	}

	// Get the CloudWatch client
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
		fmt.Println("Error getting metric data:", err)
		return nil
	}

	// Check for errors in the metric data results
	for _, metricDataResult := range result.MetricDataResults {
		if metricDataResult.StatusCode != nil && *metricDataResult.StatusCode != "Complete" {
			fmt.Printf("Error in metric data result (ID: %s): %s\n", *metricDataResult.Id, *metricDataResult.StatusCode)
		} else {
			for i, timestamp := range metricDataResult.Timestamps {
				fmt.Printf("Data for Metric at Timestamp %v: %f\n", *timestamp, *metricDataResult.Values[i])
			}
		}
	}

	return result
}

// buildDimensions converts an array of Dimensions to an array of CloudWatch Dimensions.
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

//func getCloudElementIdFromDB(cloudElementId string, query Query, infClient Client) (string, error) {
//	query.URL = "http://34.199.12.114:6057/api/cloud-element/search?id=" + strconv.Itoa(int(query.ElementId))
//	fmt.Println("query::::::::::::::::::::::::::::", query)
//	return cloudElementId, nil
//}
