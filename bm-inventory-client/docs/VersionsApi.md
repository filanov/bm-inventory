# bm_inventory_client.VersionsApi

All URIs are relative to *http://api.openshift.com/api/assisted-install/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**list_component_versions**](VersionsApi.md#list_component_versions) | **GET** /component_versions | List of componenets versions


# **list_component_versions**
> ListVersions list_component_versions()

List of componenets versions

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.VersionsApi()

try:
    # List of componenets versions
    api_response = api_instance.list_component_versions()
    pprint(api_response)
except ApiException as e:
    print("Exception when calling VersionsApi->list_component_versions: %s\n" % e)
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**ListVersions**](ListVersions.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

