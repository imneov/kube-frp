package k8s

import (
	"context"
	"fmt"

	frpv1alpha1 "github.com/fatedier/frp/api/v1alpha1"
	configtypes "github.com/fatedier/frp/pkg/config/types"
	config "github.com/fatedier/frp/pkg/config/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// GroupName is the group name used in this package
	GroupName = "frp.io"
	// Version is the version of the API
	Version = "v1alpha1"
)

// NewClient creates a new Kubernetes client
func NewClient(kubeconfig string) (client.Client, error) {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %v", err)
	}

	scheme := runtime.NewScheme()
	if err := frpv1alpha1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add scheme: %v", err)
	}

	cl, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return cl, nil
}

// LoadServerConfig loads the server configuration from a FRPServer CR
func LoadServerConfig(ctx context.Context, c client.Client, name string) (*config.ServerConfig, error) {
	frpServer := &frpv1alpha1.FRPServer{}
	err := c.Get(ctx, client.ObjectKey{Name: name}, frpServer)
	if err != nil {
		return nil, fmt.Errorf("failed to load FRPServer CR: %v", err)
	}

	// Convert FRPServerSpec to ServerConfig
	conf := &config.ServerConfig{
		Auth: config.AuthServerConfig{
			Method: config.AuthMethod(frpServer.Spec.Auth.Method),
		},
		BindAddr:              frpServer.Spec.BindAddr,
		BindPort:              frpServer.Spec.BindPort,
		KCPBindPort:           frpServer.Spec.KCPBindPort,
		QUICBindPort:          frpServer.Spec.QUICBindPort,
		ProxyBindAddr:         frpServer.Spec.ProxyBindAddr,
		VhostHTTPPort:         frpServer.Spec.VhostHTTPPort,
		VhostHTTPSPort:        frpServer.Spec.VhostHTTPSPort,
		TCPMuxHTTPConnectPort: frpServer.Spec.TCPMuxHTTPConnectPort,
		SubDomainHost:         frpServer.Spec.SubDomainHost,
		MaxPortsPerClient:     int64(frpServer.Spec.MaxPortsPerClient),
		AllowPorts:            parsePortsRange(frpServer.Spec.AllowPorts),
		EnablePrometheus:      frpServer.Spec.EnablePrometheus,
		UserConnTimeout:       int64(frpServer.Spec.UserConnTimeout),
		WebServer: config.WebServerConfig{
			Addr:     frpServer.Spec.WebServer.Addr,
			Port:     frpServer.Spec.WebServer.Port,
			User:     frpServer.Spec.WebServer.User,
			Password: frpServer.Spec.WebServer.Password,
		},
	}

	if frpServer.Spec.Transport.TLS.Force {
		conf.Transport.TLS = config.TLSServerConfig{
			Force: true,
			TLSConfig: config.TLSConfig{
				CertFile:      frpServer.Spec.Transport.TLS.CertFile,
				KeyFile:       frpServer.Spec.Transport.TLS.KeyFile,
				TrustedCaFile: frpServer.Spec.Transport.TLS.TrustedCaFile,
			},
		}
	}

	return conf, nil
}

func parsePortsRange(ports []frpv1alpha1.PortsRange) []configtypes.PortsRange {
	if ports == nil {
		return nil
	}
	confPorts := make([]configtypes.PortsRange, len(ports))
	for i, port := range ports {
		confPorts[i] = configtypes.PortsRange{
			Start: int(port.Start),
			End:   int(port.End),
		}
	}
	return confPorts
}
