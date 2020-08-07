package main

import (
	"github.com/NOVAPokemon/utils"
	userdb "github.com/NOVAPokemon/utils/database/user"
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
		locationTag := utils.GetLocationTag(utils.DefaultLocationTagsFilename, serverName)
		commsManager = utils.CreateDefaultDelayedManager(locationTag, false)
	}

	userdb.InitUsersDBClient(*flags.ArchimedesEnabled)
	utils.StartServer(serviceName, host, port, routes, commsManager)
}
