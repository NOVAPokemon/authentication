package main

import (
	"strings"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const RegisterName = "REGISTER"
const LoginName = "LOGIN"
const RefreshName = "REFRESH"

const POST = "POST"
const GET = "GET"

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(serviceName)),
	utils.Route{
		Name:        RegisterName,
		Method:      POST,
		Pattern:     api.RegisterPath,
		HandlerFunc: Register,
	},
	utils.Route{
		Name:        LoginName,
		Method:      POST,
		Pattern:     api.LoginPath,
		HandlerFunc: Login,
	},
	utils.Route{
		Name:        RefreshName,
		Method:      GET,
		Pattern:     api.RefreshPath,
		HandlerFunc: Refresh,
	},
}
