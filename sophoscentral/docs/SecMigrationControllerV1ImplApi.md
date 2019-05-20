# \SecMigrationControllerV1ImplApi

All URIs are relative to *http://api1.central.sophos.com/gateway*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EndpointsUsingGET1**](SecMigrationControllerV1ImplApi.md#EndpointsUsingGET1) | **Get** /migration-tool/v1/endpoints | Get endpoints for customer based on the parameters provided
[**GetCurrentLicensesUsingGET1**](SecMigrationControllerV1ImplApi.md#GetCurrentLicensesUsingGET1) | **Get** /migration-tool/v1/licenses/current | Get current licenses for customer
[**GetCustomerFeaturesUsingGET1**](SecMigrationControllerV1ImplApi.md#GetCustomerFeaturesUsingGET1) | **Get** /migration-tool/v1/features/current | Get current features for a customer
[**GetInstallerInfoUsingGET1**](SecMigrationControllerV1ImplApi.md#GetInstallerInfoUsingGET1) | **Get** /migration-tool/v1/deployment/agent/locations | Get the installer information.
[**HashesUsingGET1**](SecMigrationControllerV1ImplApi.md#HashesUsingGET1) | **Get** /migration-tool/v1/download/hashes | Get SHA1 hashes for all available installer templates.



## EndpointsUsingGET1

> EndpointsResponse EndpointsUsingGET1(ctx, xApiKey, authorization, optional)
Get endpoints for customer based on the parameters provided

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**xApiKey** | **string**|  | 
**authorization** | **string**|  | 
 **optional** | ***EndpointsUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a EndpointsUsingGET1Opts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **since** | **optional.String**| The timestamp to start searching from | 
 **offset** | **optional.Int32**| The paging offset | 
 **xTimestamp** | **optional.String**|  | 

### Return type

[**EndpointsResponse**](EndpointsResponse.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCurrentLicensesUsingGET1

> CurrentLicensesResponse GetCurrentLicensesUsingGET1(ctx, xApiKey, authorization, optional)
Get current licenses for customer

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**xApiKey** | **string**|  | 
**authorization** | **string**|  | 
 **optional** | ***GetCurrentLicensesUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a GetCurrentLicensesUsingGET1Opts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xTimestamp** | **optional.String**|  | 

### Return type

[**CurrentLicensesResponse**](CurrentLicensesResponse.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCustomerFeaturesUsingGET1

> CustomerFeaturesResponse GetCustomerFeaturesUsingGET1(ctx, xApiKey, authorization, optional)
Get current features for a customer

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**xApiKey** | **string**|  | 
**authorization** | **string**|  | 
 **optional** | ***GetCustomerFeaturesUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a GetCustomerFeaturesUsingGET1Opts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xTimestamp** | **optional.String**|  | 

### Return type

[**CustomerFeaturesResponse**](CustomerFeaturesResponse.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetInstallerInfoUsingGET1

> InstallerInfoResponse GetInstallerInfoUsingGET1(ctx, xApiKey, authorization, optional)
Get the installer information.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**xApiKey** | **string**|  | 
**authorization** | **string**|  | 
 **optional** | ***GetInstallerInfoUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a GetInstallerInfoUsingGET1Opts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xTimestamp** | **optional.String**|  | 

### Return type

[**InstallerInfoResponse**](InstallerInfoResponse.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## HashesUsingGET1

> HashesResponse HashesUsingGET1(ctx, xApiKey, authorization, optional)
Get SHA1 hashes for all available installer templates.

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**xApiKey** | **string**|  | 
**authorization** | **string**|  | 
 **optional** | ***HashesUsingGET1Opts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a HashesUsingGET1Opts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xTimestamp** | **optional.String**|  | 

### Return type

[**HashesResponse**](HashesResponse.md)

### Authorization

[api_key](../README.md#api_key)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

