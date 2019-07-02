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
	//"github.com/pkg/errors"

	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	"github.com/intel/intel-device-plugins-for-kubernetes/pkg/debug"
	dpapi "github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin"
)

const (
	sysfsDrmDirectory = "/sys/class/testdevice"
	devfsDriDirectory = "/dev/test"
	gpuDeviceRE       = `^card[0-9]+$`
	controlDeviceRE   = `^controlD[0-9]+$`
	vendorString      = "0x1234"

	// Device plugin settings.
	namespace  = "org.ibn4_0"
	deviceType = "sensor"
)

type devicePlugin struct {
	sysfsDir string
	devfsDir string

	sharedDevNum int
}

func newDevicePlugin(sysfsDir, devfsDir string, sharedDevNum int) *devicePlugin {
	return &devicePlugin{
		sysfsDir:         sysfsDir,
		devfsDir:         devfsDir,
		sharedDevNum:     sharedDevNum,
	}
}

func (dp *devicePlugin) Scan(notifier dpapi.Notifier) error {
        devTree := dpapi.NewDeviceTree()
	devTree.AddDevice("testdevice", "id", dpapi.DeviceInfo{
	    State:  pluginapi.Healthy,
	    Nodes: []pluginapi.DeviceSpec{
	       {
               HostPath:      "/dev/test",
               ContainerPath: "/dev/test",
       	       Permissions:   "rw",
               },
	   },
        })
	notifier.Notify(devTree)
	return nil;
}


func main() {
	var sharedDevNum int
	var debugEnabled bool

	flag.IntVar(&sharedDevNum, "shared-dev-num", 1, "number of containers sharing the same GPU device")
	flag.BoolVar(&debugEnabled, "debug", true, "enable debug output")
	flag.Parse()

	if debugEnabled {
		debug.Activate()
	}

	if sharedDevNum < 1 {
		fmt.Println("The number of containers sharing the same GPU must greater than zero")
		os.Exit(1)
	}

	fmt.Println("IoT device plugin started")

	plugin := newDevicePlugin(sysfsDrmDirectory, devfsDriDirectory, sharedDevNum)
	manager := dpapi.NewManager(namespace, plugin)
	manager.Run()
}
