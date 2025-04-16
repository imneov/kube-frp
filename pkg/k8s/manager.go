package k8s

import (
	"context"
	"fmt"
	"sync"

	v1alpha1 "github.com/fatedier/frp/api/v1alpha1"
	config "github.com/fatedier/frp/pkg/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Manager handles the Kubernetes integration for the frp server
type Manager struct {
	client     client.Client
	serverName string
	namespace  string
	mu         sync.RWMutex
	server     *config.ServerConfig
}

// NewManager creates a new Kubernetes manager
func NewManager(client client.Client, serverName string) (*Manager, error) {
	frpServer, err := LoadServerConfig(context.Background(), client, serverName)
	if err != nil {
		return nil, fmt.Errorf("failed to load server config: %v", err)
	}
	manager := &Manager{
		client:     client,
		serverName: serverName,
		server:     frpServer,
	}
	return manager, nil
}

func (m *Manager) GetFrpServer(ctx context.Context) (*v1alpha1.FRPServer, error) {
	frpServer := &v1alpha1.FRPServer{}
	err := m.client.Get(ctx, client.ObjectKey{Namespace: m.namespace, Name: m.serverName}, frpServer)
	if err != nil {
		return nil, fmt.Errorf("failed to get FrpServer: %v", err)
	}
	return frpServer, nil
}

// GetConfig returns the current FrpServer configuration
func (m *Manager) GetConfig(ctx context.Context) (*config.ServerConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.server, nil
}

// // UpdateClientConnection updates the status with a new client connection
// func (m *Manager) UpdateClientConnection(ctx context.Context, newConn v1alpha1.ConnectionStatus) error {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	server, err := m.GetFrpServer(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to get FrpServer: %v", err)
// 	}

// 	// Add the new connection to the status
// 	isExist := false
// 	for i, conn := range server.Status.ActiveConnections {
// 		if conn.ProxyName == newConn.ProxyName {
// 			isExist = true
// 			server.Status.ActiveConnections[i].ProxyType = newConn.ProxyType
// 			server.Status.ActiveConnections[i].RemoteAddr = newConn.RemoteAddr
// 			server.Status.ActiveConnections[i].ClientName = newConn.ClientName
// 			server.Status.ActiveConnections[i].LastHeartbeatTime = metav1.NewTime(time.Now())
// 		}
// 	}
// 	if !isExist {
// 		newConn.LastHeartbeatTime = metav1.NewTime(time.Now())
// 		newConn.StartTime = metav1.NewTime(time.Now())
// 		server.Status.ActiveConnections = append(server.Status.ActiveConnections, newConn)
// 	}

// 	// Update the status in Kubernetes
// 	return m.client.Status().Update(ctx, server)
// }

// // RemoveClientConnection removes a client connection from the status
// func (m *Manager) RemoveClientConnection(ctx context.Context, clientID string) error {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	server, err := m.GetFrpServer(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to get FrpServer: %v", err)
// 	}

// 	// Find and remove the connection
// 	for i, conn := range server.Status.ActiveConnections {
// 		if conn.ClientName == clientID {
// 			server.Status.ActiveConnections = append(
// 				server.Status.ActiveConnections[:i],
// 				server.Status.ActiveConnections[i+1:]...,
// 			)
// 			break
// 		}
// 	}

// 	// Update the status in Kubernetes
// 	return m.client.Status().Update(ctx, server)
// }

func (m *Manager) UpdateProxyStats(ctx context.Context, connectionStatus []*v1alpha1.ConnectionStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	server, err := m.GetFrpServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FrpServer: %v", err)
	}

	// Find and remove the connection
	server.Status.ActiveConnections = connectionStatus
	server.Status.Conditions = []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "Running",
			Message:            "frp server is running",
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: server.Generation,
		},
	}

	// Update the status in Kubernetes
	return m.client.Status().Update(ctx, server)
}
