package controllers

import (
	"net/http"

	"github.com/notefan-golang/repositories"
)

type PageController struct {
	Repository *repositories.PageRepository
}

func NewPageController(repository *repositories.PageRepository) *PageController {
	return &PageController{
		Repository: repository,
	}
}

func (controller PageController) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: -
}
