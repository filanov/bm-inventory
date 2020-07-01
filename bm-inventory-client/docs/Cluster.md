# Cluster

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**kind** | **str** | Indicates the type of this object. Will be &#39;Cluster&#39; if this is a complete object or &#39;ClusterLink&#39; if it is just a link. | 
**id** | **str** | Unique identifier of the object. | 
**href** | **str** | Self link. | 
**name** | **str** | Name of the OpenShift cluster. | [optional] 
**openshift_version** | **str** | Version of the OpenShift cluster. | [optional] 
**image_info** | [**ImageInfo**](ImageInfo.md) |  | 
**base_dns_domain** | **str** | Base domain of the cluster. All DNS records must be sub-domains of this base and include the cluster name. | [optional] 
**cluster_network_cidr** | **str** | IP address block from which Pod IPs are allocated This block must not overlap with existing physical networks. These IP addresses are used for the Pod network, and if you need to access the Pods from an external network, configure load balancers and routers to manage the traffic. | [optional] 
**cluster_network_host_prefix** | **int** | The subnet prefix length to assign to each individual node. For example, if clusterNetworkHostPrefix is set to 23, then each node is assigned a /23 subnet out of the given cidr (clusterNetworkCIDR), which allows for 510 (2^(32 - 23) - 2) pod IPs addresses. If you are required to provide access to nodes from an external network, configure load balancers and routers to manage the traffic. | [optional] 
**service_network_cidr** | **str** | The IP address pool to use for service IP addresses. You can enter only one IP address pool. If you need to access the services from an external network, configure load balancers and routers to manage the traffic. | [optional] 
**api_vip** | **str** | Virtual IP used to reach the OpenShift cluster API. | [optional] 
**machine_network_cidr** | **str** | A CIDR that all hosts belonging to the cluster should have an interfaces with IP address that belongs to this CIDR. The api_vip belongs to this CIDR. | [optional] 
**ingress_vip** | **str** | Virtual IP used for cluster ingress traffic. | [optional] 
**ssh_public_key** | **str** | SSH public key for debugging OpenShift nodes. | [optional] 
**status** | **str** | Status of the OpenShift cluster. | 
**status_info** | **str** | Additional information pertaining to the status of the OpenShift cluster. | 
**status_updated_at** | **datetime** | The last time that the cluster status has been updated | [optional] 
**hosts** | [**list[Host]**](Host.md) | Hosts that are associated with this cluster. | [optional] 
**updated_at** | **datetime** | The last time that this cluster was updated. | [optional] 
**created_at** | **datetime** | The time that this cluster was created. | [optional] 
**install_started_at** | **datetime** | The time that this cluster began installation. | [optional] 
**install_completed_at** | **datetime** | The time that this cluster completed installation. | [optional] 
**host_networks** | [**list[HostNetwork]**](HostNetwork.md) | List of host networks to be filled during query. | [optional] 
**pull_secret_set** | **bool** | True if the pull-secret has been added to the cluster | [optional] 
**ignition_generator_version** | **str** |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


