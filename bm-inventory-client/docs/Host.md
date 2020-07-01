# Host

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**kind** | **str** | Indicates the type of this object. Will be &#39;Host&#39; if this is a complete object or &#39;HostLink&#39; if it is just a link. | 
**id** | **str** | Unique identifier of the object. | 
**href** | **str** | Self link. | 
**cluster_id** | **str** | The cluster that this host is associated with. | [optional] 
**status** | **str** |  | 
**status_info** | **str** |  | 
**status_updated_at** | **datetime** | The last time that the host status has been updated | [optional] 
**connectivity** | **str** |  | [optional] 
**hardware_info** | **str** |  | [optional] 
**inventory** | **str** |  | [optional] 
**free_addresses** | **str** |  | [optional] 
**role** | **str** |  | [optional] 
**bootstrap** | **bool** |  | [optional] 
**installer_version** | **str** | Installer version | [optional] 
**updated_at** | **datetime** |  | [optional] 
**created_at** | **datetime** |  | [optional] 
**checked_in_at** | **datetime** | The last time the host&#39;s agent communicated with the service. | [optional] 
**discovery_agent_version** | **str** |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


