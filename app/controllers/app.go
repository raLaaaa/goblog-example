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

// Method called by the interceptor to check whether we are logged in via session.
// Only applicable in this controller.
func (c App) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Authentication.Login())
	}
	return nil
}

// Method which receives the POST request of our PostEntry form.
func (c App) ReceiveEntry(name string, description string) revel.Result {

	// Validates the received fields
	c.Validation.Required(name).Message("An entry name is required!")
	c.Validation.Required(description).Message("An entry description is required!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.PostEntry)
	}

	// Creates a model with the received fields and the current time
	var entry models.BlogEntry
	entry.Name = name
	entry.Description = description
	entry.CreatedAt = time.Now()
	// Saves the created entry to the DB
	services.SaveToDatabase(entry)

	c.Flash.Success("Entry created!")
	return c.Redirect(App.PostEntry)
}

// Renders our PostEntry template
func (c App) PostEntry() revel.Result {
	return c.Render()
}
