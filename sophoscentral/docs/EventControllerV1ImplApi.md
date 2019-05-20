# \EventControllerV1ImplApi

All URIs are relative to *http://api1.central.sophos.com/gateway*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetEventsUsingGET1**](EventControllerV1ImplApi.md#GetEventsUsingGET1) | **Get** /siem/v1/events | Get events for customer based on the parameters provided



## GetEventsUsingGET1

> EventAggregate GetEventsUsingGET1(ctx, xApiKey, authorization, optional)
Get events for customer based on the parameters provided

Note: Events are retrieved for timestamps within last 24 hours

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**xApiKey** | **string**|  | 
**authorization** | **string**|  | 
 **optional** | ***GetEventsUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a GetEventsUsingGET1Opts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **limit** | **optional.Int32**| The maximum number of items to return, default is 200, max is 1000 | [default to 200]
 **cursor** | **optional.String**| Identifier for next item in the list, this value is available in response as next_cursor. Response will default to last 24 hours if cursor is not within last 24 hours. | 
 **fromDate** | **optional.Int64**| The starting date from which alerts will be retrieved defined as Unix timestamp in UTC.Ignored if cursor is set. Must be within last 24 hours | 
 **excludeTypes** | **optional.String**| The String of list of types of events to be excluded | 
 **xTimestamp** | **optional.String**|  | 

### Return type

[**EventAggregate**](EventAggregate.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

