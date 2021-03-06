/*
 * Sophos Public API
 *
 * Swagger Specifications for public APIs
 *
 * API version: beta
 * Contact: support@sophos.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package sophoscentral

type EndpointEntity struct {
	AdSyncInfo                      EndpointAdSyncInfo     `json:"adSyncInfo,omitempty"`
	AlertStatus                     int32                  `json:"alert_status,omitempty"`
	AssignedProducts                []string               `json:"assignedProducts,omitempty"`
	AwsInfo                         map[string]string      `json:"awsInfo,omitempty"`
	AzureInfo                       map[string]string      `json:"azureInfo,omitempty"`
	Beta                            bool                   `json:"beta,omitempty"`
	Cloned                          bool                   `json:"cloned,omitempty"`
	DeclonedFrom                    string                 `json:"decloned_from,omitempty"`
	DeletedAt                       string                 `json:"deleted_at,omitempty"`
	DeviceEncryptionStatusUnmanaged bool                   `json:"device_encryption_status_unmanaged,omitempty"`
	EarlyAccessProgram              string                 `json:"early_access_program,omitempty"`
	EndpointType                    string                 `json:"endpoint_type,omitempty"`
	FeatureCodes                    []string               `json:"feature_codes,omitempty"`
	GroupFullName                   string                 `json:"group_full_name,omitempty"`
	GroupId                         string                 `json:"group_id,omitempty"`
	GroupName                       string                 `json:"group_name,omitempty"`
	HealthStatus                    int32                  `json:"health_status,omitempty"`
	HeartbeatUtmName                string                 `json:"heartbeat_utm_name,omitempty"`
	Id                              string                 `json:"id,omitempty"`
	Info                            map[string]interface{} `json:"info,omitempty"`
	IsAdsyncGroup                   bool                   `json:"is_adsync_group,omitempty"`
	IsCachingProxy                  bool                   `json:"is_caching_proxy,omitempty"`
	JavaId                          string                 `json:"java_id,omitempty"`
	LastActivity                    string                 `json:"last_activity,omitempty"`
	LastUser                        string                 `json:"last_user,omitempty"`
	LastUserId                      string                 `json:"last_user_id,omitempty"`
	LicenseCodes                    []string               `json:"license_codes,omitempty"`
	MachineId                       string                 `json:"machine_id,omitempty"`
	Name                            string                 `json:"name,omitempty"`
	RegisteredAt                    string                 `json:"registered_at,omitempty"`
	Status                          map[string]interface{} `json:"status,omitempty"`
	TamperProtection                TamperProtectionEntity `json:"tamper_protection,omitempty"`
	Transport                       string                 `json:"transport,omitempty"`
}
