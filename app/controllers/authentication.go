package controllers

import (
	"fmt"
	"rala-blog/app/models"
	"rala-blog/app/routes"
	"rala-blog/app/services"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	*revel.Controller
}

// Renders the login page but only if we are not logged in
func (c Authentication) Login() revel.Result {
	if user := c.connected(); user != nil {
		return c.Redirect(routes.App.PostEntry())
	}
	return c.Render()
}

// Receives the POST request with the login information and checks whether they are valid
func (c Authentication) ReceiveLogin(username, password string) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			c.Session.SetNoExpiration()
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.App.PostEntry())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.Authentication.Login())
}

// Adds the user to the session
func (c Authentication) addUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

// Checks whether the user is logged it
func (c Authentication) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}

	return nil
}

// Retrieves the user object and stores it inside the session
func (c Authentication) getUser(username string) (user *models.User) {
	user = &models.User{}
	_, err := c.Session.GetInto("fulluser", user, false)

	if user.Name == username {
		return user
	}

	*user, err = services.GetSingleUserByName(username)
	if err != nil {
		fmt.Println("Failed to find user", "user", username, "error", err)
		return nil
	}
	c.Session["fulluser"] = user
	return
}

// Receives the GET request and logs the user out
func (c Authentication) Logout() revel.Result {
	// Deletes the session
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Home.Index())
}
