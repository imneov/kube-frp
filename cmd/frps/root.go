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

	"github.com/spf13/cobra"

	"github.com/imneov/kube-frp/pkg/config"
	v1 "github.com/imneov/kube-frp/pkg/config/v1"
	"github.com/imneov/kube-frp/pkg/config/v1/validation"
	"github.com/imneov/kube-frp/pkg/util/log"
	"github.com/imneov/kube-frp/pkg/util/version"
	"github.com/imneov/kube-frp/server"

	ctrl "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	cfgFile          string
	showVersion      bool
	strictConfigMode bool

	serverCfg v1.ServerConfig
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file of frps")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of frps")
	rootCmd.PersistentFlags().BoolVarP(&strictConfigMode, "strict_config", "", true, "strict config parsing mode, unknown fields will cause errors")
	config.RegisterServerConfigFlags(rootCmd, &serverCfg)

}

var rootCmd = &cobra.Command{
	Use:   "frps",
	Short: "frps is the server of frp (https://github.com/imneov/kube-frp)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}

		var (
			svrCfg         *v1.ServerConfig
			isLegacyFormat bool
			err            error
		)

		if frpServerName != "" {
			svrCfg, err = InitManagerAndLoadConfig(frpServerName, kubeconfig)
			if err != nil {
				return err
			}
			svrCfg.EnableOperator = true
		} else if cfgFile != "" {
			svrCfg, isLegacyFormat, err = config.LoadServerConfig(cfgFile, strictConfigMode)
			if err != nil {
				return err
			}
			if isLegacyFormat {
				fmt.Printf("WARNING: ini format is deprecated and the support will be removed in the future, " +
					"please use yaml/json/toml format instead!\n")
			}
		} else {
			serverCfg.Complete()
			svrCfg = &serverCfg
		}

		warning, err := validation.ValidateServerConfig(svrCfg)
		if warning != nil {
			fmt.Printf("WARNING: %v\n", warning)
		}
		if err != nil {
			return err
		}

		if err := runServer(svrCfg); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	rootCmd.SetGlobalNormalizationFunc(config.WordSepNormalizeFunc)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServer(cfg *v1.ServerConfig) (err error) {
	log.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)

	// set controller-runtime global logger to enable client-go warnings
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	if cfgFile != "" {
		log.Infof("frps uses config file: %s", cfgFile)
	} else if frpServerName != "" {
		log.Infof("frps uses FRPServer CR as config: %s", frpServerName)
	} else {
		log.Infof("frps uses command line arguments for config")
	}

	svr, err := server.NewService(cfg)
	if err != nil {
		return err
	}
	log.Infof("frps started successfully")
	svr.Run(context.Background())
	return
}
