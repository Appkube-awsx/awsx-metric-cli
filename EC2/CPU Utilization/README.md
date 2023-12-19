- [What is awsx-CPUUtilization](#awsx-cpuutilization)
- [How to write plugin subcommand](#how-to-write-plugin-subcommand)
- [How to build / Test](#how-to-build--test)
- [what it does ](#what-it-does)
- [command input](#command-input)
- [command output](#command-output)
- [How to run ](#how-to-run)

# awsx-cpuutilization
This is a plugin subcommand for awsx cli ( https://github.com/Appkube-awsx/awsx#awsx ) cli.

For details about awsx commands and how its used in Appkube platform , please refer to the diagram below:

![alt text](https://raw.githubusercontent.com/AppkubeCloud/appkube-architectures/main/LayeredArchitecture-phase2.svg)

This plugin subcommand will implement the Apis' related to EC2 services , primarily the following API's:

- getCurrentUsage
- getAverageUsage
- getMaxUsage

This cli collect data from metric / logs / traces of the EC2 services and produce the data in a form that Appkube Platform expects.

This CLI , interacts with other Appkube services like Appkube vault , Appkube cloud CMDB so that it can talk with cloud services as 
well as filter and sort the information in terms of product/env/ services, so that Appkube platform gets the data that it expects from the cli.

# How to write plugin subcommand 
Please refer to the instruction -
https://github.com/Appkube-awsx/awsx#how-to-write-a-plugin-subcommand

It has detailed instruction on how to write a subcommand plugin , build / test / debug  / publish and integrate into the main commmand.

# How to build / Test
            go run main.go
                - Program will print Calling awsx-cpuutilization on console 

            Another way of testing is by running go install command
            go install
            - go install command creates an exe with the name of the module (e.g. awsx-cpuutilization) and save it in the GOPATH
            - Now we can execute this command on command prompt as below
            awsx-cloudelements --vaultURL=vault.dummy.net --accountId=xxxxxxxxxx --zone=us-west-2

# what it does 
This subcommand implement the following functionalities -
   getElementDetails - It  will get the resource count summary for a given AWS account id and region.

# command input
  --valutURL = URL location of vault - that stores credentials to call API
  --acountId = The AWS account id.
  --zone = AWS region
#  command output
{
  "AvgCPUUtilization24h": "19.70%",
  "CurrentCPUUtilization": "19.70%",
  "MaxCPUUtilization24h": "19.70%"
}
2023/12/19 11:55:19 {
  MetricDataResults: [{
      Id: "afreen1",
      Label: "CPUUtilization",
      StatusCode: "Complete",
      Timestamps: [
        2023-12-19 06:20:00 +0000 UTC,
        2023-12-19 06:15:00 +0000 UTC,
        2023-12-19 06:10:00 +0000 UTC,
        2023-12-19 06:05:00 +0000 UTC,
        2023-12-19 06:00:00 +0000 UTC,
        2023-12-19 05:55:00 +0000 UTC,
        2023-12-19 05:50:00 +0000 UTC,
        2023-12-19 05:45:00 +0000 UTC,
        2023-12-19 05:40:00 +0000 UTC,
        2023-12-19 05:35:00 +0000 UTC,
        2023-12-19 05:30:00 +0000 UTC,
        2023-12-19 05:25:00 +0000 UTC
      ],
      Values: [
        19.700000000000003,
        20.041666666666664,
        18.57,
        19.66166666666667,
        20.37630558564564,
        19.99666328229722,
        19.48334205762963,
        19.746666666666666,
        20.186666666666667,
        19.726666666666667,
        19.69833333333333,
        19.445
      ]


# How to run 
  From main awsx command , it is called as follows:
  awsx getElementDetails  --vaultURL=vault.dummy.net --accountId=xxxxxxxxxx --zone=us-west-2
  If you build it locally , you can simply run it as standalone command as 
  awsx-cloudelements --vaultURL=vault.dummy.net --accountId=xxxxxxxxxx --zone=us-west-2






# awsx-cpuutilization
CPUUtilization extension

# AWSX Commands for AWSX-CPUUtilization Cli's :

1. CMD used to get list of CPUUtilization instance's :

./awsx-cpuutilization --zone=us-east-1 --accessKey=<> --secretKey=<> --crossAccountRoleArn=<>  --externalId=<> --cloudWatchQueries  "[{\"RefID\": \"afreen1\",\"Query\": [{\"Namespace\": \"AWS/EC2\",\"MetricName\": \"CPUUtilization\",\"Period\": 300,\"Stat\": \"Average\",\"Dimensions\": [{\"Name\": \"InstanceId\",\"Value\": \"i-05e4e6757f13da657\"}]}]},{\"RefID\": \"afreen2\",\"Query\": [{\"Namespace\": \"AWS/EC2\",\"MetricName\": \"CPUUtilization\",\"Period\": 300,\"Stat\": \"Average\",\"Dimensions\": [{\"Name\": \"InstanceId\",\"Value\": \"i-05e4e6757f13da657\"}]}]},{\"RefID\": \"afreen3\",\"Query\": [{\"Namespace\": \"AWS/EC2\",\"MetricName\": \"CPUUtilization\",\"Period\": 300,\"Stat\": \"Average\",\"Dimensions\": [{\"Name\": \"InstanceId\",\"Value\": \"i-05e4e6757f13da657\"}]}]}]" 




