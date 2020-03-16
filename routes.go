package main

import (
	"github.com/NOVAPokemon/utils"
)

const StatusName = "STATUS"
const RegisterName = "REGISTER"
const LoginName = "LOGIN"
const RefreshName = "REFRESH"

const GET = "GET"
const POST = "POST"

var Routes = utils.Routes{
	utils.Route{
		Name: StatusName,
		Method: GET,
		Pattern:"/",
		HandlerFunc: Status,
	},
	utils.Route{
		Name: RegisterName,
		Method: POST,
		Pattern: "/register",
		HandlerFunc: Register,
	},
	utils.Route{
		Name: LoginName,
		Method: POST,
		Pattern: "/login",
		HandlerFunc: Login,
	},
	utils.Route{
		Name: RefreshName,
		Method: POST,
		Pattern: "/refresh",
		HandlerFunc: Refresh,
	},
}
