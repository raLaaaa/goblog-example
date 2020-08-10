package controllers

import (
	"rala-blog/app/services"

	"github.com/revel/revel"
)

type Home struct {
	*revel.Controller
}

// Fetches all blog entries via the service and renders it inside our Index.html template
func (c Home) Index() revel.Result {
	entries := services.GetAllEntries()
	return c.Render(entries)
}
