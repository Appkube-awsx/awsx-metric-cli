- [What is awsx-NetworkUtilization](#awsx-networkutilization)
- [How to write plugin subcommand](#how-to-write-plugin-subcommand)
- [How to build / Test](#how-to-build--test)
- [what it does ](#what-it-does)
- [command input](#command-input)
- [command output](#command-output)
- [How to run ](#how-to-run)

# awsx-networkutilization
This is a plugin subcommand for awsx cli ( https://github.com/Appkube-awsx/awsx#awsx ) cli.

For details about awsx commands and how its used in Appkube platform , please refer to the diagram below:

![alt text](https://raw.githubusercontent.com/AppkubeCloud/appkube-architectures/main/LayeredArchitecture-phase2.svg)

This plugin subcommand will implement the Apis' related to EC2 services , primarily the following API's:

- Network In
- Network Out
- Data transferred

This cli collect data from metric / logs / traces of the EC2 services and produce the data in a form that Appkube Platform expects.

This CLI , interacts with other Appkube services like Appkube vault , Appkube cloud CMDB so that it can talk with cloud services as 
well as filter and sort the information in terms of product/env/ services, so that Appkube platform gets the data that it expects from the cli.

# How to write plugin subcommand 
Please refer to the instruction -
https://github.com/Appkube-awsx/awsx#how-to-write-a-plugin-subcommand

It has detailed instruction on how to write a subcommand plugin , build / test / debug  / publish and integrate into the main commmand.

# How to build / Test
            go run main.go
                - Program will print Calling awsx-networkutilization on console 

            Another way of testing is by running go install command
            go install
            - go install command creates an exe with the name of the module (e.g. awsx-networkutilization) and save it in the GOPATH
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
{ Inbound: 3012053.60 Bytes, Outbound: 1998826.60 Bytes, DataTransfer: 5010880.20 Bytes }
2023/12/22 11:42:07 {
  MetricDataResults: [{
      Id: "a_0",
      Label: "NetworkIn",
      StatusCode: "Complete",
      Timestamps: [
        2023-12-22 06:07:00 +0000 UTC,
        2023-12-22 06:02:00 +0000 UTC,
        2023-12-22 05:57:00 +0000 UTC,
        2023-12-22 05:52:00 +0000 UTC,
        2023-12-22 05:47:00 +0000 UTC,
        2023-12-22 05:42:00 +0000 UTC,
        2023-12-22 05:37:00 +0000 UTC,
        2023-12-22 05:32:00 +0000 UTC,
        2023-12-22 05:27:00 +0000 UTC,
        2023-12-22 05:22:00 +0000 UTC,
        2023-12-22 05:17:00 +0000 UTC,
        2023-12-22 05:12:00 +0000 UTC
      ],
      Values: [
        3.0120536e+06,
        2.2863358e+06,
        2.4806378e+06,
        2.686089e+06,
        2.3005934e+06,
        2.8208448e+06,
        2.5373932e+06,
        2.4092096e+06,
        2.5031996e+06,
        2.330849e+06,
        2.6198352e+06,
        2.5706224e+06
      ]
    },{
      Id: "a_1",
      Label: "NetworkOut",
      StatusCode: "Complete",
      Timestamps: [
        2023-12-22 06:07:00 +0000 UTC,
        2023-12-22 06:02:00 +0000 UTC,
        2023-12-22 05:57:00 +0000 UTC,
        2023-12-22 05:52:00 +0000 UTC,
        2023-12-22 05:47:00 +0000 UTC,
        2023-12-22 05:42:00 +0000 UTC,
        2023-12-22 05:37:00 +0000 UTC,
        2023-12-22 05:32:00 +0000 UTC,
        2023-12-22 05:27:00 +0000 UTC,
        2023-12-22 05:22:00 +0000 UTC,
        2023-12-22 05:17:00 +0000 UTC,
        2023-12-22 05:12:00 +0000 UTC
      ],
      Values: [
        1.9988266e+06,
        1.850179e+06,
        1.8612084e+06,
        1.8479928e+06,
        1.8441908e+06,
        1.9280652e+06,
        1.9135358e+06,
        1.8273886e+06,
        2.5422918e+06,
        1.890389e+06,
        1.9081236e+06,
        1.9417618e+06
      ]
    }]
}

    
# How to run 
  From main awsx command , it is called as follows:
  awsx getElementDetails  --vaultURL=vault.dummy.net --accountId=xxxxxxxxxx --zone=us-west-2
  If you build it locally , you can simply run it as standalone command as 
  awsx-cloudelements --vaultURL=vault.dummy.net --accountId=xxxxxxxxxx --zone=us-west-2


# awsx-networkutilization
NetworkUtilization extension

# AWSX Commands for AWSX-NetworkUtilization Cli's :

1. CMD used to get list of CPUUtilization instance's :

./awsx-networkutilization --zone=us-east-1 --accessKey=<> --secretKey=<> --crossAccountRoleArn=<>  --externalId=<> --cloudWatchQueries  "[{\"RefID\": \"A\",\"Query\": [{\"Namespace\": \"AWS/EC2\",\"MetricName\": \"NetworkIn\",\"Period\": 300,\"Stat\": \"Average\",\"Dimensions\": [{\"Name\": \"InstanceId\",\"Value\": \"i-05e4e6757f13da657\"}]}]},{\"RefID\": \"A\",\"Query\": [{\"Namespace\": \"AWS/EC2\",\"MetricName\": \"NetworkOut\",\"Period\": 300,\"Stat\": \"Average\",\"Dimensions\": [{\"Name\": \"InstanceId\",\"Value\": \"i-05e4e6757f13da657\"}]}]}]"





