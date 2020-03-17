package app

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"github.com/banch0/mux/pkg/website/models"
	"net/http"
	"path/filepath"
)

func (receiver *server) handleBurgersList() func(http.ResponseWriter, *http.Request) {
	tpl, err := template.ParseFiles(filepath.Join(receiver.templatesPath, "index.gohtml"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := receiver.burgersSvc.BurgersList(context.Background())
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Burgers []models.Burger
		}{
			Title:   "McDonalds",
			Burgers: list,
		}

		err = tpl.Execute(writer, data)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (receiver *server) handleBurgersSave() func(responseWriter http.ResponseWriter, request *http.Request) {
	// POST
	return func(writer http.ResponseWriter, request *http.Request) {
		// TODO: save data in db

		// TODO: посмотреть, можно ли переделать на GET
		http.Redirect(writer, request, "/", http.StatusPermanentRedirect)
		return
	}
}

func (receiver *server) handleBurgersRemove() func(responseWriter http.ResponseWriter, request *http.Request) {
	// POST
	return func(writer http.ResponseWriter, request *http.Request) {
		// TODO: update removed = true in db
		// TODO: посмотреть, можно ли переделать на GET
		http.Redirect(writer, request, "/", http.StatusPermanentRedirect)
		return //ServiceBurgers
	}
}

func (receiver *server) handleFavicon() func(http.ResponseWriter, *http.Request) {
	// TODO: handle concurrency
	file, err := ioutil.ReadFile(filepath.Join(receiver.assetsPath, "favicon.ico"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(file)
		if err != nil {
			log.Print(err)
		}
	}
}
