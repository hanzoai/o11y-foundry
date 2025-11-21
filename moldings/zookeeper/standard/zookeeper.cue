package zookeeper

import (
	"text/template"
	"list"
	"strings"
	"strconv"
)

// 1. Define the input variables and constraints
zoo_server_count: string @tag(ZOO_SERVER_COUNT)

// 2. Define the Zookeeper config file structure (zoo.cfg)
#ZooCfg: {
	dataDir:    string
	dataLogDir: string
	clientPort: number
	servers:    string
}

#ZooServerConfig: {
	serverID: string
}

listTest: [for i in list.Range(1, strconv.Atoi(zoo_server_count)+1, 1) {
	"zookeeper-\(i):2888:3888"
}]

test: strings.Join(listTest, ",")

serverConfig: #ZooCfg & {
	dataDir:    "/path/test"
	dataLogDir: "/path/logs"
	clientPort: 2181
	servers:    test
}

cfgTemplate: """
	dataDir={{.dataDir}}
	dataLogDir={{.dataLogDir}}
	clientPort={{.clientPort}}
	servers={{.servers}}
	"""

output: template.Execute(cfgTemplate, serverConfig)
