package controllers

import (
	"rala-blog/app/models"
	"rala-blog/app/routes"
	"rala-blog/app/services"
	"time"

	"github.com/revel/revel"
)

type App struct {
	Authentication
	*revel.Controller
}

func (c App) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Authentication.Login())
	}
	return nil
}

func (c App) ReceiveEntry(name string, description string) revel.Result {

	c.Validation.Required(name).Message("An entry name is required!")
	c.Validation.Required(description).Message("An entry description is required!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.PostEntry)
	}

	var entry models.BlogEntry
	entry.Name = name
	entry.Description = description
	entry.CreatedAt = time.Now()
	services.SaveToDatabase(entry)

	c.Flash.Success("Entry created!")
	return c.Redirect(App.PostEntry)
}

func (c App) PostEntry() revel.Result {
	return c.Render()
}
