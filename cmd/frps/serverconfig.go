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
	"os"
	"path/filepath"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/k8s"
	"github.com/spf13/pflag"
)

var (
	frpServerName string
	kubeconfig    string
	k8sFlags      = pflag.NewFlagSet("Kubernetes", pflag.ExitOnError)
)

func init() {
	k8sFlags.StringVar(&frpServerName, "kube-frps-config", "", "[Kubernetes] name of the FRPServer CR to load config from")
	k8sFlags.StringVar(&kubeconfig, "kubeconfig", "", "[Kubernetes] path to the kubeconfig file (defaults to standard locations if not specified)")

	rootCmd.PersistentFlags().AddFlagSet(k8sFlags)
}

// InitManagerAndLoadConfig loads configuration from FRPServer CR
func InitManagerAndLoadConfig(frpServerName string, kubeconfig string) (*v1.ServerConfig, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %v", err)
			}
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
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

	return svrCfg, nil
}
