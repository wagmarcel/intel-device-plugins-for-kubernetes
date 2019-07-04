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
	"flag"
	"fmt"
	"os"
	"time"
	"encoding/json"
	"io/ioutil"
	"strconv"
	 
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
	"github.com/intel/intel-device-plugins-for-kubernetes/pkg/debug"
	dpapi "github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin"
)

const (
	//sysfsDrmDirectory = "/sys/class/testdevice"
	//devfsDriDirectory = "/dev/test"
	vendorString      = "0x1234"

	// Device plugin settings.
	namespace  = "org.industry-business-network"
	deviceType = "preconfigured"
)

type devicePlugin struct {
	sysfsDir string
	devfsDir string
	sharedDevNum int
}

const deviceConfigFileName = "/etc/deviceconfig.json"

func newDevicePlugin(sharedDevNum int) *devicePlugin {
	return &devicePlugin{
		sharedDevNum:     sharedDevNum,
	}
}

//scanning is currently mocked by a config file in "/etc/deviceconfig"
// should be: {"name": name, "id": id, "hostPath": hostPath,
// "containerPath": containerPath, "permissions": permissions }
func (dp *devicePlugin) Scan(notifier dpapi.Notifier) error {
        devTree := dpapi.NewDeviceTree()

	deviceConfigFile, err := os.Open(deviceConfigFileName)
	if err != nil {
	   panic(err)
	}
	defer deviceConfigFile.Close()
	byteValue, _ := ioutil.ReadAll(deviceConfigFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	fmt.Print("Name: ", result["name"])

	for i := 0; i < dp.sharedDevNum; i++ {
	    id := result["id"].(string) + strconv.Itoa(i)
	    devTree.AddDevice(result["name"].(string), id, dpapi.DeviceInfo{
	        State:  pluginapi.Healthy,
	    	Nodes: []pluginapi.DeviceSpec{
	       	       {
               	       HostPath:      result["hostPath"].(string),
               	       ContainerPath: result["containerPath"].(string),
       	       	       Permissions:   result["permission"].(string),
               	       },
	   	},
       	    })
	}
	notifier.Notify(devTree)
	return nil;
}


func main() {
	var sharedDevNum int
	var debugEnabled bool
	var devicePluginPath string

	flag.IntVar(&sharedDevNum, "shared-dev-num", 1, "number of containers sharing the same GPU device")
	flag.BoolVar(&debugEnabled, "debug", false, "enable debug output")
	flag.StringVar(&devicePluginPath, "device-plugin-path", pluginapi.DevicePluginPath,  "device plugin")
	flag.Parse()

	if debugEnabled {
		debug.Activate()
	}

	if sharedDevNum < 1 {
		fmt.Println("The number of containers sharing the same GPU must greater than zero")
		os.Exit(1)
	}

	fmt.Println("IoT device plugin started")

	plugin := newDevicePlugin(sharedDevNum)
	manager := dpapi.NewManager(namespace, plugin, devicePluginPath)
	manager.Run()
	for {
		fmt.Println("Infinite Loop 1")
	   	time.Sleep(time.Second * 5)
	}
}
