NOTE - Pass an instance flag for all instances or specific instance, along with other params

1. CPU utilization panel, - Metric Namespace - CPUUtilization, Stat Average(5m), Average(24Hours), Maximum(24h), Units - %, Dimmensions - All instances

execution - awsx-metric-cli cpuutilization --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " " 


2. Memory Utilization panel, - Metric Namespace - MemoryUtilization, Stat Average(5m), Avergae (24h), Maximum(24h),Units - Bytes, Dimmensions - All instances.

execution - awsx-metric-cli MemoryUtilization --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

Note - This panel is fetched through cloudwatch agent / metrics.


3. Storage Utilization - AWS service -EBS, Metric Namespace - Storage Utilization, Stat Average(5m), Avergae (24h), Maximum(24h),Units - Bytes. Dimmensions - for all volumes.
Please write a program for Root Volume usage by defining MountedPath (/), and the rest shall be retrieved by Volume id for avg (24h),max (24h)

execution - awsx-metric-cli StorageUtilization --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

4. Network Utilization - AWS service - EC2, stat Average, Units - Bytes. Dimmensions - for all instances. Metric used - net_bytes_recv (InBound), net_bytes_sent (OutBound) --

execution - awsx-metric-cli NetworkUtilization --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

5. cpu_user - AWS service EC2, stat - average, Metric Namespace - AWS EC2, stat average, Units - None, Dimmensions - all instances, Metric name - cpu_time_user

Note - This metric is defined by CW agent

execution - awsx-metric-cli cpu_time_user --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "


6. cpu_idle - AWS service EC2, stat - average, Metric Namespace - AWS EC2, Units - None, Dimmensions - all instances, metric name - cpu_usage_idle.
Note - This metric is defined by CW agent

execution - awsx-metric-cli cpu_idle --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

7. cpu_sys - AWS service EC2, stat - average, Metric Namespace - AWS EC2, Units - None, Dimmensions - all instances, metric name - cpu_time_system.
Note - This metric is defined by CW agent

execution - awsx-metric-cli cpu_sys --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

8.  cpu_nice - AWS service EC2, stat - average, Metric Namespace - AWS EC2, Units - None, Dimmensions - all instances, metric name - cpu_time_nice.
Note - This metric is defined by CW agent

execution - awsx-metric-cli cpu_nice --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

9. mem_total - AWS service EC2, stat - average, Metric Namespace - AWS EC2, Units none, Dimmensions - all instances, metric name - mem_total.
Note - this metric is a part of CW agent.

execution - awsx-metric-cli mem_total --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

9. mem_free - AWS service EC2, stat - average, Metric Namespace - AWS EC2, Units none, Dimmensions - all instances, metric name - mem_free.
Note - this metric is a part of CW agent.

execution - awsx-metric-cli mem_free --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

10. mem_used - AWS service EC2, stat - average , Metric Namespace -AWS EC2, Units None, Dimmensions - all instances, metric name - mem_free.
Note - this metric is a part of Cloudwatch Agent.

execution - awsx-metric-cli mem_used --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

11. mem_physicalRam - 



12. disk_reads - AWS service EC2, stat - average, Metric namespace - AWS EC2, Units Count, Dimmensions - all instances, metric name - diskio_reads.
Note - this metric is defined by CW agent, The only statistic that should be used for this metric is Sum. Do not use Average.

execution - awsx-metric-cli diskio_reads --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "


13. disk_writes - AWS service EC2, stat - average, Metric namespace - AWS EC2, Units Count, Dimmensions - all instances, metric name - diskio_writes.
Note this metric is defined by CW agent, The only statistic that should be used for this metric is Sum. Do not use Average.

execution - awsx-metric-cli diskio_writes --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "


14. disk_used - AWS service EC2, stat - average, Metric Namespace - AWS EC2, units - bytes, Dimmensions - all instances, metric name - disk_used,
Note - This metric is part of CW agent. 

execution - awsx-metric-cli disk_used --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

15. Net_InPackets - AWS service EC2, stat - sum, metric namespace - AWS EC2, units -  count , Dimmensions - all instances, metric name - NetworkPacketsIn.

execution - awsx-metric-cli Net_InPackets --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

16. Net_OutPackets - AWS service EC2, stat - sum, metric namespace - AWS EC2, units - count, Dimmensions - all instances, metric name - NetworkPacketsOut.

execution - awsx-metric-cli Net_OutPackets --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

17. Net_InBytes - AWS service EC2, stat - sum, metric namespace - AWS ec2, units - Bytes, Dimmensions - all instances, metric name - net_bytes_recv.
Note - This is a CW agent metric, The only statistic that should be used for this metric is Sum. Do not use Average.

execution - awsx-metric-cli Net_InBytes --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

18. Net_outBytes AWS service EC2, stat - sum, metric namespace - AWS ec2, units - Bytes, Dimmensions - all instances, metric name - net_bytes_sent.
Note - This a cw agent metric, The only statistic that should be used for this metric is Sum. Do not use Average.

execution - awsx-metric-cli Net_outBytes --acKey= " " --secKey= " " crossaccount = " " zone = " " externalID = " "

--> Availability dashboard -- 

19. Instance status panel - Custom API's 
(i)   DescribeInstanceStatus - This will retrieve the instance status and Availability zone
(ii)  DescribeInstance -  This will get the metadata for an instance but we have to only parse the instance type
(iii) DescribeAlarms - This is a cloudwatch API where we can describe alarms for an instance if attached. 

execution - awsx-metric-cli Instancestatus --acKey = " " --secKey = " " --externalId = " " --crossaccountrolearn = " " --zone = " " --instanceID = " " 

20. Alerts and Notifications - Custom API
(i) DescribeAlarms - This will retrieve the information for alerts and only parse Date, Alarm Name & description.

execution - awsx-metric-cli AlertsandNotifications --acKey = " " --secKey = " " --externalId = " " --crossaccountrolearn = " " --zone = " " --instanceID = " "

21. Instances Health Check - Custom API
(i)  Describe instances - This API will provide metadata of an ec2 instance we have to parse Availability zone, instance ID, instance type and instance state.
(ii) DescribeInstanceStatus - This will retrieve the instance status which is system check in ec2 design.

execution - awsx-metric-cli HealthCheck --acKey = " " --secKey = " " --externalId = " " --crossaccountrolearn = " " --zone = " " --instanceID = " "

