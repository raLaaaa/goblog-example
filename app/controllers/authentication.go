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

func (c Authentication) Login() revel.Result {
	if user := c.connected(); user != nil {
		return c.Redirect(routes.App.PostEntry())
	}
	return c.Render()
}

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

func (c Authentication) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c Authentication) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}

	return nil
}

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

func (c Authentication) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Home.Index())
}
