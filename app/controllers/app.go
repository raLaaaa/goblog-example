package controllers

import (
	"../models/"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	user User

	return c.Render()
}
