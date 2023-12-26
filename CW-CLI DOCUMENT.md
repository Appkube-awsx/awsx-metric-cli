
## AWS cloudwatch-metric-cli Documentation

### Overview

cli to monitors AWS resources using cloudwatch metric queries. It is written in Go and customizable with various parameters.

### Prerequisites
- Go installed.

### Command Details
```
go run .\main.go \
  --crossAccountRoleArn=arn:aws:iam::<account-id>:role/CrossAccount \
  --cloudWatchQueries="[{\"RefID\": \"A\",\"MaxDataPoint\": 100,\"Interval\": 60,\"TimeRange\": {\"From\": \"\",\"To\": \"\",\"TimeZone\": \"UTC\"},\"Query\": [{\"Namespace\": \"AWS/EC2\",\"MetricName\": \"CPUUtilization\",\"Period\": 300,\"Stat\": \"Average\"}]}]"

```
### Command Parameter:
- --crossAccountRoleArn: AWS IAM role ARN for cross-account access.
- -cloudWatchQueries: JSON array of CloudWatch queries.
       mandatory paramters of cloudWatchQueries
            1. RefID
            2. Namespace
            3. MetricName
            4. Period
            5. Stat
    
### Logic to get GLOBAL_AWS_SECRETS (access/secret key) in cli: 
        Since we are only passing crossAccountRoleArn, we need GLOBAL_AWS_SECRETS (access/secret key) from vault. It can be retrieved by two ways explaind below: 
            1. make vault call with static key (GLOBAL_AWS_SECRETS)
            2. If vault is not available, get the GLOBAL_AWS_SECRETS from environment variable
            3. If GLOBAL_AWS_SECRETS not found in environment variable, program should exit with error - clien connection could not be established. access/secret key not found

# Proposed changes
        1. awsx-common 
        GLOBAL_AWS_SECRETS logic described in above para should be implemented in awsx-common layer. awsx-common is responsible to do authentication and provide aws connection based on aws element types (e.g. cloudwatch-metric etc..)
        2. Appkube-cloud-datasource
        Once cli changes are done and validated the above command, we need to make following changes in Appkube-cloud-datasource
            2.1 crossRoleArn for the elementId 
	        2.2 Full and final query 
                    NOTE: since appkube-cloud-datasource is able to make a json with all the required query params, we don't need any tranformation in this json in api layer. So pass this query json to cli as it is. cli will parse this json to make cloudwatch-query-input

# Integration with awsx-metric api
    http://<server>:port/awsx-metrics


