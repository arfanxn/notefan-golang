package controllers

import (
	"encoding/json"
	"net/http"
	"notion-golang/helper"
	"notion-golang/repositories"
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
	pages := controller.Repo.Get(r.Context())

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(pages)
	helper.LogFatalIfError(err)
}
