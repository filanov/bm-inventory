# bm_inventory_client.EventsApi

All URIs are relative to *http://api.openshift.com/api/assisted-install/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**list_events**](EventsApi.md#list_events) | **GET** /events/{entity_id} | Lists events for an entity_id


# **list_events**
> EventList list_events(entity_id)

Lists events for an entity_id

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.EventsApi()
entity_id = 'entity_id_example' # str | 

try:
    # Lists events for an entity_id
    api_response = api_instance.list_events(entity_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling EventsApi->list_events: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **entity_id** | [**str**](.md)|  | 

### Return type

[**EventList**](EventList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

