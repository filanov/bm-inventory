# ClusterUpdateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** | OpenShift cluster name | [optional] 
**base_dns_domain** | **str** | Base domain of the cluster. All DNS records must be sub-domains of this base and include the cluster name. | [optional] 
**cluster_network_cidr** | **str** | IP address block from which Pod IPs are allocated This block must not overlap with existing physical networks. These IP addresses are used for the Pod network, and if you need to access the Pods from an external network, configure load balancers and routers to manage the traffic. | [optional] 
**cluster_network_host_prefix** | **int** | The subnet prefix length to assign to each individual node. For example, if clusterNetworkHostPrefix is set to 23, then each node is assigned a /23 subnet out of the given cidr (clusterNetworkCIDR), which allows for 510 (2^(32 - 23) - 2) pod IPs addresses. If you are required to provide access to nodes from an external network, configure load balancers and routers to manage the traffic. | [optional] 
**service_network_cidr** | **str** | The IP address pool to use for service IP addresses. You can enter only one IP address pool. If you need to access the services from an external network, configure load balancers and routers to manage the traffic. | [optional] 
**api_vip** | **str** | Virtual IP used to reach the OpenShift cluster API. | [optional] 
**ingress_vip** | **str** | Virtual IP used for cluster ingress traffic. | [optional] 
**pull_secret** | **str** | The pull secret that obtained from the Pull Secret page on the Red Hat OpenShift Cluster Manager site. | [optional] 
**ssh_public_key** | **str** | SSH public key for debugging OpenShift nodes. | [optional] 
**hosts_roles** | [**list[ClusterupdateparamsHostsRoles]**](ClusterupdateparamsHostsRoles.md) | The desired role for hosts associated with the cluster. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


