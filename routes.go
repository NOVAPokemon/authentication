package main

import (
	"strings"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const registerName = "REGISTER"
const loginName = "LOGIN"
const refreshName = "REFRESH"

const post = "POST"
const get = "GET"

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(serviceName)),
	utils.Route{
		Name:        registerName,
		Method:      post,
		Pattern:     api.RegisterPath,
		HandlerFunc: register,
	},
	utils.Route{
		Name:        loginName,
		Method:      post,
		Pattern:     api.LoginPath,
		HandlerFunc: login,
	},
	utils.Route{
		Name:        refreshName,
		Method:      get,
		Pattern:     api.RefreshPath,
		HandlerFunc: refresh,
	},
}
