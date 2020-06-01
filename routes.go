package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"strings"
)

const RegisterName = "REGISTER"
const LoginName = "LOGIN"
const RefreshName = "REFRESH"

const POST = "POST"

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(fmt.Sprintf("/%s", serviceName))),
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
