// Copyright 2017 Intel Corporation. All Rights Reserved.
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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/intel/intel-device-plugins-for-kubernetes/pkg/debug"
	dpapi "github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

const (
	//sysfsDrmDirectory = "/sys/class/testdevice"
	//devfsDriDirectory = "/dev/test"
	//vendorString = "0x1234"

	// Device plugin settings.
	namespace  = "oisp.net"
	//deviceType = "preconfigured"
)

type devicePlugin struct {
	sysfsDir     string
	devfsDir     string
	sharedDevNum int
}

const deviceConfigFileName = "/etc/oisp/deviceconfig.json"

func newDevicePlugin(sharedDevNum int) *devicePlugin {
	return &devicePlugin{
		sharedDevNum: sharedDevNum,
	}
}

//scanning is currently mocked by a either config file in "/etc/oisp/deviceconfig" or env variable PLUGIN_CONFIG
// should be: {"name": name, "id": id, "hostPath": hostPath,
// "containerPath": containerPath, "permissions": permissions, "gid": gid }
func (dp *devicePlugin) Scan(notifier dpapi.Notifier) error {
	devTree := dpapi.NewDeviceTree()

	// There are two potential sources for deviceConfig:
	// (1) config file in /etc/oisp/deviceconfig.json
	// (2) env variable PLUGIN_CONFIG
	var byteValue []byte
	deviceConfigFile, err := os.Open(deviceConfigFileName)
	defer deviceConfigFile.Close()
	if err != nil {
		log.Info("No file found in /etc/oisp, now trying env variable K8S_PLUGIN_CONFIG. Err from reading file: ", err)
		if byteValue = []byte(os.Getenv("K8S_PLUGIN_CONFIG")); len(byteValue) == 0 {
			panic("No configuration found.")
		}
	} else {
		byteValue, _ = ioutil.ReadAll(deviceConfigFile)
	}
	var result []map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	log.Info("Found resources: ", result)
	log.Info("Number of resources: ", len(result))

	for i := 0; i < len(result); i++ {
		id := result[i]["id"].(string)
		devTree.AddDevice(result[i]["name"].(string), id, dpapi.DeviceInfo{
			State: pluginapi.Healthy,
			Nodes: []pluginapi.DeviceSpec{
				{
					HostPath:      result[i]["hostPath"].(string),
					ContainerPath: result[i]["containerPath"].(string),
					Permissions:   result[i]["permission"].(string),
				},
			},
		})
		gid := int(result[i]["gid"].(float64))
		hostPath := result[i]["hostPath"].(string)
		log.Info("Changing device " + hostPath + " gid to ", gid)
		info, err := os.Stat(hostPath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Warn("File does not exist on Host.")
			}
			panic("Cannot manage host device")
		} else {
			var uid int;
			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				uid = int(stat.Uid)
				log.Info("UID of device " + hostPath + ": ", uid)
				err = os.Chown(hostPath, uid, gid)
				if err != nil {
					log.Error(err)
					log.Warn("No permission to update gid!")
				}
			} else {
				log.Warn("No permissions to update gid!")
			}
		}
	}

	//change gid of device fix
	notifier.Notify(devTree)
	return nil
}

func main() {
	var sharedDevNum int
	var debugEnabled bool
	var devicePluginPath string

	flag.BoolVar(&debugEnabled, "debug", false, "enable debug output")
	flag.StringVar(&devicePluginPath, "device-plugin-path", pluginapi.DevicePluginPath, "device plugin")
	flag.Parse()

	if debugEnabled {
		debug.Activate()
	}

	sharedDevNum = 1

	fmt.Println("IoT device plugin started")

	plugin := newDevicePlugin(sharedDevNum)
	manager := dpapi.NewManager(namespace, plugin, devicePluginPath, devicePluginPath+"kubelet.sock")
	manager.Run()
	for {
		fmt.Println("Heartbeat")
		time.Sleep(time.Second * 10)
	}
}
