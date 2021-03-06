{
	"data_directory": "./data",
	"listeners": {
		"rest_port": 8001,
		"udp_port": 41244,
		"tcp_port": 7071
	},
	"receivers": {
		"udp_port": 41245,
		"udp_address": "0.0.0.0"
	},
	"logger": {
		"level": "info",
		"path": "/tmp/",
		"max_size": 134217728
	},
	"default_connector": "rest+ws",
	"connector": {
		"rest": {
			"host": "latest.streammyiot.com",
			"port": 443,
			"protocol": "https",
			"strictSSL": false,
			"timeout": 30000,
			"proxy": {
				"host": false,
				"port": false
			}
		},
		"ws": {
			"host": "latest.streammyiot.com",
			"port": 443,
			"minRetryTime": 2500,
			"maxRetryTime": 600000,
			"testTimeout": 40000,
			"pingPongIntervalMs": 30000,
			"enablePingPong": true,
			"secure": true,
			"proxy": {
				"host": false,
				"port": false
			}
		}
	}
}
