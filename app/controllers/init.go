package controllers

import "github.com/revel/revel"

// Interceptor methods which are called before corrosponding requests.
// For instance: All requests in the Authentication-Controller will call AddUser
// All requests in the App-Controller will call checkUser. This will realize our login.

func init() {
	revel.InterceptMethod(Authentication.addUser, revel.BEFORE)
	revel.InterceptMethod(App.checkUser, revel.BEFORE)
}
