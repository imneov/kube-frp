package k8s

import (
	"context"
	"fmt"

	v1alpha1 "github.com/fatedier/frp/api/v1alpha1"
	config "github.com/fatedier/frp/pkg/config/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	mgr *Manager
)

// Init initializes the frp server client
func Init(client client.Client, serverName string) error {
	manager, err := NewManager(client, serverName)
	if err != nil {
		klog.Errorf("failed to create frp server client: %v", err)
		return err
	}
	mgr = manager
	klog.Info("Initialized frp server client")
	return nil
}

func GetFrpServer(ctx context.Context) (*v1alpha1.FRPServer, error) {
	if mgr == nil {
		return nil, fmt.Errorf("frp server client not initialized")
	}
	return mgr.GetFrpServer(ctx)
}

func GetConfig(ctx context.Context) (*config.ServerConfig, error) {
	if mgr == nil {
		return nil, fmt.Errorf("frp server client not initialized")
	}
	return mgr.GetConfig(ctx)
}

// func HandleNewProxy(ctx context.Context, inMsg *msg.NewProxy, userInfo plugin.UserInfo, remoteAddr string) error {
// 	if mgr == nil {
// 		return fmt.Errorf("frp server client not initialized")
// 	}
// 	//更新frpserver的status
// 	err := mgr.UpdateClientConnection(ctx, v1alpha1.ConnectionStatus{
// 		ProxyName:  inMsg.ProxyName,
// 		ProxyType:  inMsg.ProxyType,
// 		RemoteAddr: remoteAddr,
// 	})

// 	if err != nil {
// 		klog.Errorf("failed to update frp server status: %v", err)
// 		return err
// 	}
// 	return nil
// }

// func HandleCloseProxy(ctx context.Context, proxyName, proxyType string) error {
// 	if mgr == nil {
// 		return fmt.Errorf("frp server client not initialized")
// 	}
// 	err := mgr.RemoveClientConnection(ctx, proxyName)
// 	if err != nil {
// 		klog.Errorf("failed to remove frp server status: %v", err)
// 		return err
// 	}
// 	return nil
// }

func UpdateProxyStats(ctx context.Context, connectionStatus []*v1alpha1.ConnectionStatus) error {
	if mgr == nil {
		return fmt.Errorf("frp server client not initialized")
	}
	return mgr.UpdateProxyStats(ctx, connectionStatus)
}
