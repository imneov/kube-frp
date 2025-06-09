// Copyright 2018 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	v1 "github.com/imneov/kube-frp/pkg/config/v1"
	"github.com/imneov/kube-frp/pkg/k8s"
	"github.com/spf13/pflag"
)

var (
	internetAddr  string
	frpServerName string
	kubeconfig    string
	k8sFlags      = pflag.NewFlagSet("Kubernetes", pflag.ExitOnError)
)

func init() {
	k8sFlags.StringVar(&internetAddr, "frps-internet-host", "", "[Kubernetes] Public IP or DNS name used by clients to connect to frps")
	k8sFlags.StringVar(&frpServerName, "frp-server-cr", "", "[Kubernetes] Name of the FRPServer CustomResource (CR) to load configuration")
	k8sFlags.StringVar(&kubeconfig, "kubeconfig", "", "[Kubernetes] Absolute path to the kubeconfig file (if running out-of-cluster)")

	rootCmd.PersistentFlags().AddFlagSet(k8sFlags)
}

// InitManagerAndLoadConfig loads configuration from FRPServer CR
func InitManagerAndLoadConfig(frpServerName string, kubeconfig string) (*v1.ServerConfig, error) {

	if internetAddr == "" {
		internetAddr = getOutboundIP().String()
	}

	client, err := k8s.NewClient(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	err = k8s.Init(client, frpServerName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize frp server client: %v", err)
	}

	svrCfg, err := k8s.GetConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration from FRPServer CR: %v", err)
	}

	if svrCfg != nil {
		svrCfg.Complete()
	}

	svrCfg.InternetAddr = internetAddr

	return svrCfg, nil
}

func getOutboundIP() net.IP {
	// 连接到一个“理论上的”公网地址（Google DNS），不会真正建立连接
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("Failed to determine outbound IP: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
