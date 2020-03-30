package main

import (
	"github.com/NOVAPokemon/authentication/exported"
	"github.com/NOVAPokemon/utils"
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
		Pattern:     exported.StatusPath,
		HandlerFunc: Status,
	},
	utils.Route{
		Name:        RegisterName,
		Method:      POST,
		Pattern:     exported.RegisterPath,
		HandlerFunc: Register,
	},
	utils.Route{
		Name:        LoginName,
		Method:      POST,
		Pattern:     exported.LoginPath,
		HandlerFunc: Login,
	},
	utils.Route{
		Name:        RefreshName,
		Method:      POST,
		Pattern:     exported.RefreshPath,
		HandlerFunc: Refresh,
	},
}
