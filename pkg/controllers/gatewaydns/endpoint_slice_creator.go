// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package gatewaydns

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	discoveryv1beta1 "k8s.io/api/discovery/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	connectivityv1alpha1 "github.com/vmware-tanzu/cross-cluster-connectivity/apis/connectivity/v1alpha1"
)

func ConvertGatewaysToEndpointSlices(clusterGateways []ClusterGateway, gatewayDNSNamespace, controllerNamespace string) []discoveryv1beta1.EndpointSlice {
	var endpointSlices []discoveryv1beta1.EndpointSlice
	for _, clusterGateway := range clusterGateways {
		endpointSlices = append(endpointSlices, convertServiceToEndpointSlice(clusterGateway.Gateway, clusterGateway.ClusterName, gatewayDNSNamespace, controllerNamespace))
	}
	return endpointSlices
}

func convertServiceToEndpointSlice(service corev1.Service, clusterName string, gatewayDNSNamespace string, controllerNamespace string) discoveryv1beta1.EndpointSlice {
	// TODO: xcc.test TLD should be a configuration option
	hostname := fmt.Sprintf("*.gateway.%s.%s.clusters.xcc.test", clusterName, gatewayDNSNamespace)
	name := fmt.Sprintf("%s-%s-gateway", gatewayDNSNamespace, clusterName)
	addresses := []string{}

	for _, ingress := range service.Status.LoadBalancer.Ingress {
		addresses = append(addresses, ingress.IP)
	}

	return discoveryv1beta1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: controllerNamespace,
			Annotations: map[string]string{
				connectivityv1alpha1.DNSHostnameAnnotation: hostname,
			},
		},
		AddressType: discoveryv1beta1.AddressTypeIPv4,
		Endpoints: []discoveryv1beta1.Endpoint{
			{
				Addresses: addresses,
			},
		},
	}
}
