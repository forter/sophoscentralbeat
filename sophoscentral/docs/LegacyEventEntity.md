# LegacyEventEntity

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppCerts** | [**[]EndpointCoreEventCertificate**](EndpointCoreEventCertificate.md) | Certificate info of the application associated with the threat, if available | [optional] 
**AppSha256** | **string** | SHA 256 hash of the application associated with the threat, if available | [optional] 
**CoreRemedyItems** | [**CoreRemedyItems**](CoreRemedyItems.md) |  | [optional] 
**CreatedAt** | **string** | The date at which the event was created | [optional] 
**CustomerId** | **string** | The identifier of the customer for which record is created | [optional] 
**Details** | [**[]EventDetailProperty**](EventDetailProperty.md) |  | [optional] 
**EndpointId** | **string** | The corresponding endpoint id associated with the record | [optional] 
**EndpointType** | **string** | The corresponding endpoint type associated with the record | [optional] 
**Group** | **string** | The group associated with the group | [optional] 
**Id** | **string** | The Identifier for the event | [optional] 
**Location** | **string** | The location captured for this record | [optional] 
**Name** | **string** | The name of the record created | [optional] 
**Origin** | **string** | originating component of a detection | [optional] 
**Severity** | **string** | The severity for this alert | [optional] 
**Source** | **string** | The source for this record | [optional] 
**SourceInfo** | **map[string]string** | Detailed source information for this record | [optional] 
**Threat** | **string** | The threat associated with the record | [optional] 
**Type** | **string** | The type of this record | [optional] 
**UserId** | **string** | The identifier of the user for which record is created | [optional] 
**When** | **string** | The date at which the event was created | [optional] 
**WhitelistProperties** | [**[]EndpointWhitelistProperties**](EndpointWhitelistProperties.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


