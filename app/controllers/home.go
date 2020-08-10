package controllers

import (
	"rala-blog/app/services"

	"github.com/revel/revel"
)

type Home struct {
	*revel.Controller
}

func (c Home) Index() revel.Result {
	entries := services.GetAllEntries()
	return c.Render(entries)
}
