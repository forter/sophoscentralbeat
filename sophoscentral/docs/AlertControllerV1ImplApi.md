# \AlertControllerV1ImplApi

All URIs are relative to *http://api1.central.sophos.com/gateway*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAlertsUsingGET1**](AlertControllerV1ImplApi.md#GetAlertsUsingGET1) | **Get** /siem/v1/alerts | Get alerts for customer based on the parameters provided



## GetAlertsUsingGET1

> AlertAggregate GetAlertsUsingGET1(ctx, xApiKey, authorization, optional)
Get alerts for customer based on the parameters provided

Note: Alerts are retrieved for timestamps within last 24 hours

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**xApiKey** | **string**|  | 
**authorization** | **string**|  | 
 **optional** | ***GetAlertsUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a GetAlertsUsingGET1Opts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **limit** | **optional.Int32**| The maximum number of items to return, default is 200, max is 1000 | [default to 200]
 **fromDate** | **optional.Int64**| The starting date from which alerts will be retrieved defined as Unix timestamp in UTC. Ignored if cursor is set. Must be within last 24 hours. | 
 **cursor** | **optional.String**| Identifier for next item in the list, this value is available in response as next_cursor. Response will default to last 24 hours if cursor is not within last 24 hours. | 
 **xTimestamp** | **optional.String**|  | 

### Return type

[**AlertAggregate**](AlertAggregate.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

