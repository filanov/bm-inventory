# bm_inventory_client.InstallerApi

All URIs are relative to *http://api.openshift.com/api/assisted-install/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**deregister_cluster**](InstallerApi.md#deregister_cluster) | **DELETE** /clusters/{cluster_id} | Deletes an OpenShift bare metal cluster definition.
[**deregister_host**](InstallerApi.md#deregister_host) | **DELETE** /clusters/{cluster_id}/hosts/{host_id} | Deregisters an OpenShift bare metal host.
[**disable_host**](InstallerApi.md#disable_host) | **DELETE** /clusters/{cluster_id}/hosts/{host_id}/actions/enable | Disables a host for inclusion in the cluster.
[**download_cluster_files**](InstallerApi.md#download_cluster_files) | **GET** /clusters/{cluster_id}/downloads/files | Downloads files relating to the installed/installing cluster.
[**download_cluster_iso**](InstallerApi.md#download_cluster_iso) | **GET** /clusters/{cluster_id}/downloads/image | Downloads the OpenShift per-cluster discovery ISO.
[**download_cluster_kubeconfig**](InstallerApi.md#download_cluster_kubeconfig) | **GET** /clusters/{cluster_id}/downloads/kubeconfig | Downloads the kubeconfig file for this cluster.
[**enable_host**](InstallerApi.md#enable_host) | **POST** /clusters/{cluster_id}/hosts/{host_id}/actions/enable | Enables a host for inclusion in the cluster.
[**generate_cluster_iso**](InstallerApi.md#generate_cluster_iso) | **POST** /clusters/{cluster_id}/downloads/image | Creates a new OpenShift per-cluster discovery ISO.
[**get_cluster**](InstallerApi.md#get_cluster) | **GET** /clusters/{cluster_id} | Retrieves the details of the OpenShift bare metal cluster.
[**get_credentials**](InstallerApi.md#get_credentials) | **GET** /clusters/{cluster_id}/credentials | Get the the cluster admin credentials.
[**get_host**](InstallerApi.md#get_host) | **GET** /clusters/{cluster_id}/hosts/{host_id} | Retrieves the details of the OpenShift bare metal host.
[**get_next_steps**](InstallerApi.md#get_next_steps) | **GET** /clusters/{cluster_id}/hosts/{host_id}/instructions | Retrieves the next operations that the host agent needs to perform.
[**install_cluster**](InstallerApi.md#install_cluster) | **POST** /clusters/{cluster_id}/actions/install | Installs the OpenShift bare metal cluster.
[**list_clusters**](InstallerApi.md#list_clusters) | **GET** /clusters | Retrieves the list of OpenShift bare metal clusters.
[**list_hosts**](InstallerApi.md#list_hosts) | **GET** /clusters/{cluster_id}/hosts | Retrieves the list of OpenShift bare metal hosts.
[**post_step_reply**](InstallerApi.md#post_step_reply) | **POST** /clusters/{cluster_id}/hosts/{host_id}/instructions | Posts the result of the operations from the host agent.
[**register_cluster**](InstallerApi.md#register_cluster) | **POST** /clusters | Creates a new OpenShift bare metal cluster definition.
[**register_host**](InstallerApi.md#register_host) | **POST** /clusters/{cluster_id}/hosts | Registers a new OpenShift bare metal host.
[**set_debug_step**](InstallerApi.md#set_debug_step) | **POST** /clusters/{cluster_id}/hosts/{host_id}/actions/debug | Sets a single shot debug step that will be sent next time the host agent will ask for a command.
[**update_cluster**](InstallerApi.md#update_cluster) | **PATCH** /clusters/{cluster_id} | Updates an OpenShift bare metal cluster definition.
[**update_host_install_progress**](InstallerApi.md#update_host_install_progress) | **PUT** /clusters/{clusterId}/hosts/{hostId}/progress | Update installation progress
[**upload_cluster_ingress_cert**](InstallerApi.md#upload_cluster_ingress_cert) | **POST** /clusters/{cluster_id}/uploads/ingress-cert | Transfer the ingress certificate for the cluster.


# **deregister_cluster**
> deregister_cluster(cluster_id)

Deletes an OpenShift bare metal cluster definition.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 

try:
    # Deletes an OpenShift bare metal cluster definition.
    api_instance.deregister_cluster(cluster_id)
except ApiException as e:
    print("Exception when calling InstallerApi->deregister_cluster: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deregister_host**
> deregister_host(cluster_id, host_id)

Deregisters an OpenShift bare metal host.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
host_id = 'host_id_example' # str | 

try:
    # Deregisters an OpenShift bare metal host.
    api_instance.deregister_host(cluster_id, host_id)
except ApiException as e:
    print("Exception when calling InstallerApi->deregister_host: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **host_id** | [**str**](.md)|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **disable_host**
> disable_host(cluster_id, host_id)

Disables a host for inclusion in the cluster.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
host_id = 'host_id_example' # str | 

try:
    # Disables a host for inclusion in the cluster.
    api_instance.disable_host(cluster_id, host_id)
except ApiException as e:
    print("Exception when calling InstallerApi->disable_host: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **host_id** | [**str**](.md)|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **download_cluster_files**
> file download_cluster_files(cluster_id, file_name)

Downloads files relating to the installed/installing cluster.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
file_name = 'file_name_example' # str | 

try:
    # Downloads files relating to the installed/installing cluster.
    api_response = api_instance.download_cluster_files(cluster_id, file_name)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->download_cluster_files: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **file_name** | **str**|  | 

### Return type

[**file**](file.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/octet-stream

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **download_cluster_iso**
> str download_cluster_iso(cluster_id)

Downloads the OpenShift per-cluster discovery ISO.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 

try:
    # Downloads the OpenShift per-cluster discovery ISO.
    api_response = api_instance.download_cluster_iso(cluster_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->download_cluster_iso: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 

### Return type

**str**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/octet-stream

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **download_cluster_kubeconfig**
> str download_cluster_kubeconfig(cluster_id)

Downloads the kubeconfig file for this cluster.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 

try:
    # Downloads the kubeconfig file for this cluster.
    api_response = api_instance.download_cluster_kubeconfig(cluster_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->download_cluster_kubeconfig: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 

### Return type

**str**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/octet-stream

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **enable_host**
> enable_host(cluster_id, host_id)

Enables a host for inclusion in the cluster.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
host_id = 'host_id_example' # str | 

try:
    # Enables a host for inclusion in the cluster.
    api_instance.enable_host(cluster_id, host_id)
except ApiException as e:
    print("Exception when calling InstallerApi->enable_host: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **host_id** | [**str**](.md)|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **generate_cluster_iso**
> Cluster generate_cluster_iso(cluster_id, image_create_params)

Creates a new OpenShift per-cluster discovery ISO.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
image_create_params = bm_inventory_client.ImageCreateParams() # ImageCreateParams | 

try:
    # Creates a new OpenShift per-cluster discovery ISO.
    api_response = api_instance.generate_cluster_iso(cluster_id, image_create_params)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->generate_cluster_iso: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **image_create_params** | [**ImageCreateParams**](ImageCreateParams.md)|  | 

### Return type

[**Cluster**](Cluster.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_cluster**
> Cluster get_cluster(cluster_id)

Retrieves the details of the OpenShift bare metal cluster.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 

try:
    # Retrieves the details of the OpenShift bare metal cluster.
    api_response = api_instance.get_cluster(cluster_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->get_cluster: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 

### Return type

[**Cluster**](Cluster.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_credentials**
> Credentials get_credentials(cluster_id)

Get the the cluster admin credentials.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 

try:
    # Get the the cluster admin credentials.
    api_response = api_instance.get_credentials(cluster_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->get_credentials: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 

### Return type

[**Credentials**](Credentials.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_host**
> Host get_host(cluster_id, host_id)

Retrieves the details of the OpenShift bare metal host.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
host_id = 'host_id_example' # str | 

try:
    # Retrieves the details of the OpenShift bare metal host.
    api_response = api_instance.get_host(cluster_id, host_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->get_host: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **host_id** | [**str**](.md)|  | 

### Return type

[**Host**](Host.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **get_next_steps**
> Steps get_next_steps(cluster_id, host_id)

Retrieves the next operations that the host agent needs to perform.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
host_id = 'host_id_example' # str | 

try:
    # Retrieves the next operations that the host agent needs to perform.
    api_response = api_instance.get_next_steps(cluster_id, host_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->get_next_steps: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **host_id** | [**str**](.md)|  | 

### Return type

[**Steps**](Steps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **install_cluster**
> Cluster install_cluster(cluster_id)

Installs the OpenShift bare metal cluster.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 

try:
    # Installs the OpenShift bare metal cluster.
    api_response = api_instance.install_cluster(cluster_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->install_cluster: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 

### Return type

[**Cluster**](Cluster.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **list_clusters**
> ClusterList list_clusters()

Retrieves the list of OpenShift bare metal clusters.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()

try:
    # Retrieves the list of OpenShift bare metal clusters.
    api_response = api_instance.list_clusters()
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->list_clusters: %s\n" % e)
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**ClusterList**](ClusterList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **list_hosts**
> HostList list_hosts(cluster_id)

Retrieves the list of OpenShift bare metal hosts.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 

try:
    # Retrieves the list of OpenShift bare metal hosts.
    api_response = api_instance.list_hosts(cluster_id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->list_hosts: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 

### Return type

[**HostList**](HostList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **post_step_reply**
> post_step_reply(cluster_id, host_id, reply=reply)

Posts the result of the operations from the host agent.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
host_id = 'host_id_example' # str | 
reply = bm_inventory_client.StepReply() # StepReply |  (optional)

try:
    # Posts the result of the operations from the host agent.
    api_instance.post_step_reply(cluster_id, host_id, reply=reply)
except ApiException as e:
    print("Exception when calling InstallerApi->post_step_reply: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **host_id** | [**str**](.md)|  | 
 **reply** | [**StepReply**](StepReply.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **register_cluster**
> Cluster register_cluster(new_cluster_params)

Creates a new OpenShift bare metal cluster definition.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
new_cluster_params = bm_inventory_client.ClusterCreateParams() # ClusterCreateParams | 

try:
    # Creates a new OpenShift bare metal cluster definition.
    api_response = api_instance.register_cluster(new_cluster_params)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->register_cluster: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **new_cluster_params** | [**ClusterCreateParams**](ClusterCreateParams.md)|  | 

### Return type

[**Cluster**](Cluster.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **register_host**
> Host register_host(cluster_id, new_host_params)

Registers a new OpenShift bare metal host.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
new_host_params = bm_inventory_client.HostCreateParams() # HostCreateParams | 

try:
    # Registers a new OpenShift bare metal host.
    api_response = api_instance.register_host(cluster_id, new_host_params)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->register_host: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **new_host_params** | [**HostCreateParams**](HostCreateParams.md)|  | 

### Return type

[**Host**](Host.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **set_debug_step**
> set_debug_step(cluster_id, host_id, step)

Sets a single shot debug step that will be sent next time the host agent will ask for a command.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
host_id = 'host_id_example' # str | 
step = bm_inventory_client.DebugStep() # DebugStep | 

try:
    # Sets a single shot debug step that will be sent next time the host agent will ask for a command.
    api_instance.set_debug_step(cluster_id, host_id, step)
except ApiException as e:
    print("Exception when calling InstallerApi->set_debug_step: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **host_id** | [**str**](.md)|  | 
 **step** | [**DebugStep**](DebugStep.md)|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **update_cluster**
> Cluster update_cluster(cluster_id, cluster_update_params)

Updates an OpenShift bare metal cluster definition.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
cluster_update_params = bm_inventory_client.ClusterUpdateParams() # ClusterUpdateParams | 

try:
    # Updates an OpenShift bare metal cluster definition.
    api_response = api_instance.update_cluster(cluster_id, cluster_update_params)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling InstallerApi->update_cluster: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **cluster_update_params** | [**ClusterUpdateParams**](ClusterUpdateParams.md)|  | 

### Return type

[**Cluster**](Cluster.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **update_host_install_progress**
> update_host_install_progress(cluster_id, host_id, host_install_progress_params)

Update installation progress

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | The ID of the cluster to retrieve
host_id = 'host_id_example' # str | The ID of the host to retrieve
host_install_progress_params = bm_inventory_client.HostInstallProgressParams() # HostInstallProgressParams | New progress value

try:
    # Update installation progress
    api_instance.update_host_install_progress(cluster_id, host_id, host_install_progress_params)
except ApiException as e:
    print("Exception when calling InstallerApi->update_host_install_progress: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)| The ID of the cluster to retrieve | 
 **host_id** | [**str**](.md)| The ID of the host to retrieve | 
 **host_install_progress_params** | [**HostInstallProgressParams**](HostInstallProgressParams.md)| New progress value | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **upload_cluster_ingress_cert**
> upload_cluster_ingress_cert(cluster_id, ingress_cert_params)

Transfer the ingress certificate for the cluster.

### Example
```python
from __future__ import print_function
import time
import bm_inventory_client
from bm_inventory_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = bm_inventory_client.InstallerApi()
cluster_id = 'cluster_id_example' # str | 
ingress_cert_params = bm_inventory_client.IngressCertParams() # IngressCertParams | 

try:
    # Transfer the ingress certificate for the cluster.
    api_instance.upload_cluster_ingress_cert(cluster_id, ingress_cert_params)
except ApiException as e:
    print("Exception when calling InstallerApi->upload_cluster_ingress_cert: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cluster_id** | [**str**](.md)|  | 
 **ingress_cert_params** | [**IngressCertParams**](IngressCertParams.md)|  | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

