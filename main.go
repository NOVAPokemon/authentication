package main

import (
	"os"

	"github.com/NOVAPokemon/utils"
	userdb "github.com/NOVAPokemon/utils/database/user"
	http "github.com/bruno-anjos/archimedesHTTPClient"
	cedUtils "github.com/bruno-anjos/cloud-edge-deployment/pkg/utils"
	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
)

const (
	host        = utils.ServeHost
	port        = utils.AuthenticationPort
	serviceName = "AUTHENTICATION"
)

func main() {
	flags := utils.ParseFlags(serverName)

	if !*flags.LogToStdout {
		utils.SetLogFile(serverName)
	}

	if !*flags.DelayedComms {
		commsManager = utils.CreateDefaultCommunicationManager()
	} else {
		commsManager = utils.CreateDefaultDelayedManager(false)
	}

	location, exists := os.LookupEnv("LOCATION")
	if !exists {
		log.Fatal("no location in environment")
	}

	var node string
	node, exists = os.LookupEnv(cedUtils.NodeIPEnvVarName)
	if !exists {
		log.Panicf("no NODE_IP env var")
	} else {
		log.Infof("Node IP: %s", node)
	}

	httpClient.InitArchimedesClient(node, http.DefaultArchimedesPort, s2.CellIDFromToken(location).LatLng())

	userdb.InitUsersDBClient(*flags.ArchimedesEnabled)
	utils.StartServer(serviceName, host, port, routes, commsManager)
}
