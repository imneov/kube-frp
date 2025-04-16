/*
Copyright 2024 The frp Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TransportConfig defines the transport configuration
type TransportConfig struct {
	// BandwidthLimit specifies the bandwidth limit
	// +optional
	BandwidthLimit string `json:"bandwidthLimit,omitempty"`

	// BandwidthLimitMode specifies the bandwidth limit mode
	// +optional
	BandwidthLimitMode string `json:"bandwidthLimitMode,omitempty"`
}

// LoadBalancerConfig defines the load balancer configuration
type LoadBalancerConfig struct {
	// Group specifies the load balancer group
	// +optional
	Group string `json:"group,omitempty"`
}

// HealthCheckConfig defines the health check configuration
type HealthCheckConfig struct {
	// Type specifies the health check type
	// +optional
	Type string `json:"type,omitempty"`

	// IntervalSeconds specifies the health check interval in seconds
	// +optional
	IntervalSeconds int32 `json:"intervalSeconds,omitempty"`
}

// ProxyConfig defines the proxy configuration
type ProxyConfig struct {
	// Name specifies the proxy name
	Name string `json:"name"`

	// Type specifies the proxy type
	Type string `json:"type"`

	// Transport specifies the transport configuration
	// +optional
	Transport *TransportConfig `json:"transport,omitempty"`

	// LoadBalancer specifies the load balancer configuration
	// +optional
	LoadBalancer *LoadBalancerConfig `json:"loadBalancer,omitempty"`

	// HealthCheck specifies the health check configuration
	// +optional
	HealthCheck *HealthCheckConfig `json:"healthCheck,omitempty"`

	// LocalIP specifies the local IP address
	// +optional
	LocalIP string `json:"localIP,omitempty"`

	// RemotePort specifies the remote port
	// +optional
	RemotePort int32 `json:"remotePort,omitempty"`
}

// TrafficStats defines the traffic statistics
type TrafficStats struct {
	// Name specifies the name of the entity
	Name string `json:"name"`

	// TrafficIn specifies the incoming traffic history
	TrafficIn []int64 `json:"trafficIn"`

	// TrafficOut specifies the outgoing traffic history
	TrafficOut []int64 `json:"trafficOut"`
}

// ConnectionStatus defines the status of a single connection
type ConnectionStatus struct {
	// ClientName is the name of the client that established the connection
	ClientName string `json:"clientName"`

	// ClientVersion is the version of the client
	ClientVersion string `json:"clientVersion"`

	// ProxyName is the name of the proxy that this connection belongs to
	ProxyName string `json:"proxyName"`

	// ProxyType is the type of the proxy (tcp, udp, http, https, etc.)
	ProxyType string `json:"proxyType"`

	// ProxyConfig contains the proxy configuration
	ProxyConfig ProxyConfig `json:"proxyConfig"`

	// LocalAddr is the local address of the connection
	LocalAddr string `json:"localAddr"`

	// RemoteAddr is the remote address of the connection
	RemoteAddr string `json:"remoteAddr"`

	// StartTime is when the connection was established
	// +optional
	StartTime metav1.Time `json:"startTime"`

	// LastCloseTime is when the connection was last closed
	// +optional
	LastCloseTime metav1.Time `json:"lastCloseTime,omitempty"`

	// BytesIn is the number of bytes received
	BytesIn int64 `json:"bytesIn"`

	// BytesOut is the number of bytes sent
	BytesOut int64 `json:"bytesOut"`

	// TodayTrafficIn is the number of bytes received today
	TodayTrafficIn int64 `json:"todayTrafficIn"`

	// TodayTrafficOut is the number of bytes sent today
	TodayTrafficOut int64 `json:"todayTrafficOut"`

	// CurrentConnections is the current number of connections
	CurrentConnections int32 `json:"currentConnections"`

	// Status is the current status of the connection
	Status string `json:"status"`

	// TrafficStats contains the traffic statistics
	// +optional
	TrafficStats *TrafficStats `json:"trafficStats,omitempty"`
}

// ProxyTypeCount defines the count of proxies by type
type ProxyTypeCount struct {
	// Type is the proxy type
	Type string `json:"type"`
	// Count is the number of proxies of this type
	Count int32 `json:"count"`
}

// ServiceStatus defines the overall status of the FRP service
type ServiceStatus struct {
	// TotalConnections is the total number of active connections
	TotalConnections int32 `json:"totalConnections"`

	// TotalClients is the total number of connected clients
	TotalClients int32 `json:"totalClients"`

	// TotalProxies is the total number of configured proxies
	TotalProxies int32 `json:"totalProxies"`

	// TotalBytesIn is the total number of bytes received
	TotalBytesIn int64 `json:"totalBytesIn"`

	// TotalBytesOut is the total number of bytes sent
	TotalBytesOut int64 `json:"totalBytesOut"`

	// Uptime is the duration since the service started
	// +optional
	Uptime metav1.Duration `json:"uptime"`

	// LastUpdateTime is when the status was last updated
	// +optional
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// FRPServerStatus defines the observed state of FRPServer
type FRPServerStatus struct {
	// Conditions represents the latest available observations of an object's current state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Ready indicates whether the FRPServer is ready to serve requests
	// +optional
	Ready bool `json:"ready,omitempty"`

	// ObservedGeneration is the most recent generation observed by the controller
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// ServiceStatus contains the overall status of the FRP service
	// +optional
	ServiceStatus ServiceStatus `json:"serviceStatus,omitempty"`

	// ActiveConnections contains the status of all active connections
	// +optional
	ActiveConnections []*ConnectionStatus `json:"activeConnections,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Connections",type="integer",JSONPath=".status.serviceStatus.totalConnections"
// +kubebuilder:printcolumn:name="Clients",type="integer",JSONPath=".status.serviceStatus.totalClients"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// FRPServer is the Schema for the frpservers API
type FRPServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FRPServerSpec   `json:"spec,omitempty"`
	Status FRPServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FRPServerList contains a list of FRPServer
type FRPServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FRPServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FRPServer{}, &FRPServerList{})
}
