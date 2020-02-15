package controllers

import (
	"fmt"
	"github.com/madjlzz/madlens/models"
	"github.com/madjlzz/madlens/views"
	"log"
	"net/http"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs: gs,
	}
}

type Galleries struct {
	New *views.View
	gs models.GalleryService
}

type GalleryForm struct {
	Title     string `schema:"title"`
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	gallery := models.Gallery{
		Title:     form.Title,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}