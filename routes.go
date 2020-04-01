package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const StatusName = "STATUS"
const RegisterName = "REGISTER"
const LoginName = "LOGIN"
const RefreshName = "REFRESH"

const GET = "GET"
const POST = "POST"

var routes = utils.Routes{
	utils.Route{
		Name:        StatusName,
		Method:      GET,
		Pattern:     api.StatusPath,
		HandlerFunc: Status,
	},
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
		Method:      POST,
		Pattern:     api.RefreshPath,
		HandlerFunc: Refresh,
	},
}
