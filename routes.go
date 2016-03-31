package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		TodoCreate,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		TodoShow,
	},
	Route{
		"CreateUser",
		"POST",
		"/users/create",
		CreateUser,
	},
	/**
		не понятно пока нужен ли метод и обработка его внутри хендлера
	**/
	/**	Route{
		"CreateUser",
		"GET",
		"/users/create",
		CreateUser,
	},**/
}
