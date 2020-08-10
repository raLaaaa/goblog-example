package controllers

import "github.com/revel/revel"

func init() {
	revel.InterceptMethod(Authentication.AddUser, revel.BEFORE)
	revel.InterceptMethod(App.checkUser, revel.BEFORE)
}
