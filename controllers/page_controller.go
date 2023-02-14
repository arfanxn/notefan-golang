package controllers

import (
	"net/http"
	"notefan-golang/repositories"
)

type PageController struct {
	Repo *repositories.PageRepo
}

func NewPageController(repo *repositories.PageRepo) *PageController {
	return &PageController{
		Repo: repo,
	}
}

func (controller PageController) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: -
}
