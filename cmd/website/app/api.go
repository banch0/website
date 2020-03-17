package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/banch0/mux/pkg/website/models"

	"github.com/banch0/mux/cmd/website/app/dto"
)

// Delete One Burger
func (s *server) handleDeleteBurgers() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(path.Base(req.URL.Path))
		if err != nil {
			return
		}
		log.Println("ID::", id)

		err = s.burgersSvc.BurgersDelete(req.Context(), id)
		if err != nil {
			http.Error(
				res,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
		log.Println("deleted")
		res.WriteHeader(http.StatusNoContent)
	}
}

// Get One Burger By ID
func (s *server) getBurgerByID() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Println("get by id")
		id, err := strconv.Atoi(path.Base(req.URL.Path))
		if err != nil {
			http.Error(
				res,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
		log.Println("Get by ID::", id)

		burger, err := s.burgersSvc.GetBurgerByID(req.Context(), id)
		if err != nil {
			log.Println(err)
		}
		log.Println(burger)

		encoded, err := json.Marshal(burger)
		if err != nil {
			log.Println("Marshal Error: ", encoded)
		}
		res.Header().Set("Content-Type", "application/json")
		_, err = res.Write(encoded)
		if err != nil {
			http.Error(
				res,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
	}
}

// Get All Burgers
func (s *server) handleAllBurgers() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Println("get all handlers")
		burgers, err := s.burgersSvc.BurgersList(req.Context())
		if err != nil {
			http.Error(
				res,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		DTOs := make([]dto.BurgerDTO, len(burgers))
		for ix, burger := range burgers {
			DTOs[ix] = dto.BurgerDTO{
				ID:    burger.ID,
				Name:  burger.Name,
				Price: burger.Price,
			}
		}

		encoded, err := json.Marshal(DTOs)
		if err != nil {
			log.Println(errors.New("Marshal Error: "), err)
		}
		res.Header().Set("Content-Type", "application/json")
		_, err = res.Write(encoded)
		if err != nil {
			http.Error(
				res,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
	}
}

// Create Burger
func (s *server) handleBurgerSave() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(
				res,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
			return
		}

		// TODO: http.MaxBytesReader()
		var DTOs dto.BurgerDTO
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("ioutil error", err)
			return
		}

		err = json.Unmarshal(body, &DTOs)
		if err != nil {
			log.Println("unmarshal error")
			http.Error(
				res,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
			return
		}

		burger := models.Burger{
			ID:      DTOs.ID,
			Name:    DTOs.Name,
			Price:   DTOs.Price,
			Removed: false,
		}

		err = s.burgersSvc.SaveBurger(req.Context(), burger)
		if err != nil {
			log.Println(err)
		}
		log.Println(burger)
		res.WriteHeader(http.StatusNoContent)
	}
}

// Update Burger
func (s *server) handleBurgerUpdate() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(path.Base(req.URL.Path))
		if err != nil {
			return
		}
		log.Println("ID::", id)

		var model *models.Burger

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
		}
		defer req.Body.Close()

		err = json.Unmarshal(body, &model)
		if err != nil {
			log.Println("Marshal Error: ", err)
			http.Error(
				res,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
			return
		}

		burger := models.Burger{
			ID:      model.ID,
			Name:    model.Name,
			Price:   model.Price,
			Removed: model.Removed,
		}

		err = s.burgersSvc.UpdateByID(req.Context(), burger, id)
		if err != nil {
			log.Println(err)
		}
		log.Println(burger)
		res.WriteHeader(http.StatusOK)
	}
}
