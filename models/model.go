package models

type CmdbCloudElementResponse struct {
	Id                       int64                  `json:"id"`
	ElementType              string                 `json:"elementType,omitempty"`
	HostedServices           map[string]interface{} `json:"hostedServices,omitempty"`
	Arn                      string                 `json:"arn,omitempty"`
	InstanceId               string                 `json:"instanceId,omitempty"`
	InstanceName             string                 `json:"instanceName,omitempty"`
	Category                 string                 `json:"category,omitempty"`
	SlaJson                  map[string]interface{} `json:"slaJson,omitempty"`
	CostJson                 map[string]interface{} `json:"costJson,omitempty"`
	ViewJson                 map[string]interface{} `json:"viewJson,omitempty"`
	ConfigJson               map[string]interface{} `json:"configJson,omitempty"`
	ComplianceJson           map[string]interface{} `json:"complianceJson,omitempty"`
	Status                   string                 `json:"status,omitempty"`
	CreatedBy                string                 `json:"createdBy,omitempty"`
	UpdatedBy                string                 `json:"updatedBy,omitempty"`
	CreatedOn                string                 `json:"createdOn,omitempty"`
	UpdatedOn                string                 `json:"updatedOn,omitempty"`
	LandingzoneId            int64                  `json:"landingzoneId"`
	LandingZone              string                 `json:"landingZone,omitempty"`
	DbCategoryId             int64                  `json:"dbCategoryId"`
	DbCategoryName           string                 `json:"dbCategoryName,omitempty"`
	ProductEnclaveId         int64                  `json:"productEnclaveId"`
	ProductEnclaveInstanceId string                 `json:"productEnclaveInstanceId,omitempty"`
}

type AwsCredential struct {
	Region              string `json:"region,omitempty"`
	AccessKey           string `json:"accessKey,omitempty"`
	SecretKey           string `json:"secretKey,omitempty"`
	CrossAccountRoleArn string `json:"crossAccountRoleArn,omitempty"`
	ExternalId          string `json:"externalId,omitempty"`
}
