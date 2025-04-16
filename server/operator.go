// Copyright 2017 fatedier, fatedier@gmail.com
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

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatedier/frp/api/v1alpha1"
	configv1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/k8s"
	"github.com/fatedier/frp/pkg/metrics/mem"
	"github.com/fatedier/frp/pkg/util/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// func (svr *Service) updateProxyStats(proxyType string) (proxyInfos []*ProxyStatsInfo) {
// 	proxyStats := mem.StatsCollector.GetProxiesByType(proxyType)
// 	proxyInfos = make([]*ProxyStatsInfo, 0, len(proxyStats))
// 	for _, ps := range proxyStats {
// 		proxyInfo := &ProxyStatsInfo{}
// 		if pxy, ok := svr.pxyManager.GetByName(ps.Name); ok {
// 			content, err := json.Marshal(pxy.GetConfigurer())
// 			if err != nil {
// 				log.Warnf("marshal proxy [%s] conf info error: %v", ps.Name, err)
// 				continue
// 			}
// 			proxyInfo.Conf = getConfByType(ps.Type)
// 			if err = json.Unmarshal(content, &proxyInfo.Conf); err != nil {
// 				log.Warnf("unmarshal proxy [%s] conf info error: %v", ps.Name, err)
// 				continue
// 			}
// 			proxyInfo.Status = "online"
// 			if pxy.GetLoginMsg() != nil {
// 				proxyInfo.ClientVersion = pxy.GetLoginMsg().Version
// 			}
// 		} else {
// 			proxyInfo.Status = "offline"
// 		}
// 		proxyInfo.Name = ps.Name
// 		proxyInfo.TodayTrafficIn = ps.TodayTrafficIn
// 		proxyInfo.TodayTrafficOut = ps.TodayTrafficOut
// 		proxyInfo.CurConns = ps.CurConns
// 		proxyInfo.LastStartTime = ps.LastStartTime
// 		proxyInfo.LastCloseTime = ps.LastCloseTime
// 		proxyInfos = append(proxyInfos, proxyInfo)
// 	}
// 	return
// }

func (svr *Service) updateProxyStatsInterval() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			svr.updateProxyStats(svr.ctx)
		case <-svr.ctx.Done():
			return
		}
	}
}

func (svr *Service) updateProxyStats(ctx context.Context) {
	if !svr.cfg.EnableOperator {
		return
	}

	types := []string{"tcp", "udp", "serverinfo", "stcp", "sudp"}
	connectionStatus := make([]*v1alpha1.ConnectionStatus, 0)
	for _, proxyType := range types {
		proxyStats := mem.StatsCollector.GetProxiesByType(proxyType)
		for _, ps := range proxyStats {
			proxyInfo := &v1alpha1.ConnectionStatus{}
			if pxy, ok := svr.pxyManager.GetByName(ps.Name); ok {
				cfg := pxy.GetConfigurer()
				switch cfg := cfg.(type) {
				case *configv1.TCPProxyConfig:
					//{"name":"your_name_here","type":"tcp","transport":{"bandwidthLimit":"","bandwidthLimitMode":"client"},"loadBalancer":{"group":""},"healthCheck":{"type":"","intervalSeconds":0},"localIP":"127.0.0.1","plugin":null,"remotePort":1711}
					proxyInfo.ProxyConfig = v1alpha1.ProxyConfig{
						Name:       cfg.Name,
						Type:       cfg.Type,
						LocalIP:    cfg.LocalIP,
						RemotePort: int32(cfg.RemotePort),
					}

					proxyInfo.ProxyConfig.Transport = &v1alpha1.TransportConfig{
						BandwidthLimit:     cfg.Transport.BandwidthLimit.String(),
						BandwidthLimitMode: cfg.Transport.BandwidthLimitMode,
					}

					proxyInfo.ProxyConfig.LoadBalancer = &v1alpha1.LoadBalancerConfig{
						Group: cfg.LoadBalancer.Group,
					}

					proxyInfo.ProxyConfig.HealthCheck = &v1alpha1.HealthCheckConfig{
						Type:            cfg.HealthCheck.Type,
						IntervalSeconds: int32(cfg.HealthCheck.IntervalSeconds),
					}

				case *configv1.UDPProxyConfig:
					//{"name":"your_name_here","type":"tcp","transport":{"bandwidthLimit":"","bandwidthLimitMode":"client"},"loadBalancer":{"group":""},"healthCheck":{"type":"","intervalSeconds":0},"localIP":"127.0.0.1","plugin":null,"remotePort":1711}
					proxyInfo.ProxyConfig = v1alpha1.ProxyConfig{
						Name:       cfg.Name,
						Type:       cfg.Type,
						LocalIP:    cfg.LocalIP,
						RemotePort: int32(cfg.RemotePort),
					}

					proxyInfo.ProxyConfig.Transport = &v1alpha1.TransportConfig{
						BandwidthLimit:     cfg.Transport.BandwidthLimit.String(),
						BandwidthLimitMode: cfg.Transport.BandwidthLimitMode,
					}

					proxyInfo.ProxyConfig.LoadBalancer = &v1alpha1.LoadBalancerConfig{
						Group: cfg.LoadBalancer.Group,
					}

					proxyInfo.ProxyConfig.HealthCheck = &v1alpha1.HealthCheckConfig{
						Type:            cfg.HealthCheck.Type,
						IntervalSeconds: int32(cfg.HealthCheck.IntervalSeconds),
					}
					proxyInfo.ProxyType = "udp"
				// case *configv1.STCPProxyConfig:
				// 	proxyInfo.ProxyType = "stcp"
				// case *configv1.SUDPProxyConfig:
				// 	proxyInfo.ProxyType = "sudp"
				default:
					log.Warnf("unknown proxy type: %s", proxyType)
					continue
				}
				content, _ := json.Marshal(cfg)
				fmt.Println(string(content))
				proxyInfo.Status = "online"
				if pxy.GetLoginMsg() != nil {
					proxyInfo.ClientVersion = pxy.GetLoginMsg().Version
				}
			} else {
				proxyInfo.Status = "offline"
			}
			proxyInfo.ProxyName = ps.Name
			proxyInfo.ProxyType = ps.Type
			proxyInfo.CurrentConnections = int32(ps.CurConns)
			proxyInfo.StartTime = metav1.Time{Time: parseTime(ps.LastStartTime)}
			proxyInfo.LastCloseTime = metav1.Time{Time: parseTime(ps.LastCloseTime)}
			proxyInfo.TodayTrafficIn = ps.TodayTrafficIn
			proxyInfo.TodayTrafficOut = ps.TodayTrafficOut
			// RemoteAddr:         proxyInfo.RemoteAddr,
			// LocalAddr:          proxyInfo.LocalAddr,
			// BytesIn:            proxyInfo.BytesIn,
			// BytesOut:           proxyInfo.BytesOut,
			connectionStatus = append(connectionStatus, proxyInfo)
		}
	}

	err := k8s.UpdateProxyStats(ctx, connectionStatus)
	if err != nil {
		log.Warnf("update proxy stats error: %v", err)
	}

	return
}

func parseTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}
	// First try RFC3339 format
	startTime, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return startTime
	}

	// If that fails, try the custom format "MM-DD HH:MM:SS"
	startTime, err = time.Parse("01-02 15:04:05", timeStr)
	if err != nil {
		log.Warnf("parse proxy [%s] start time error: %v", timeStr, err)
		return time.Time{}
	}
	return startTime
}
