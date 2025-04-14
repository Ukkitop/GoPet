package handlers

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"realAPI/api/apiUtils"
	"realAPI/api/types"
	"realAPI/internal/database"
	"realAPI/internal/database/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var p types.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		log.Error(err)
	}

	userError := errors.New("Invalid credentials")

	if err != nil {
		apiUtils.RequestErrorHandler(w, userError)
	}

	dbCon, dbErr := database.GetDatabaseConnection()

	if dbErr != nil {
		log.Error(dbErr)
		return
	}

	user := models.User{Username: p.Username}

	result := dbCon.First(&user)

	if result.Error != nil {
		log.Error(result.Error)
		apiUtils.InternalErrorHandler(w)
	}

	if p.Password == user.Password {
		token, err := apiUtils.CreateToken(user.Username)

		if err != nil {
			log.Error(err)
			apiUtils.InternalErrorHandler(w)
		}

		apiUtils.WriteSuccess(w, types.LoginResponse{Token: token})
	} else {
		apiUtils.RequestErrorHandler(w, userError)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	apiUtils.WriteSuccess(w, r)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var p types.RegisterRequest

	userError := json.NewDecoder(r.Body).Decode(&p)

	if userError != nil {
		apiUtils.RequestErrorHandler(w, userError)
	}

	userModel := models.User{Username: p.Username, Email: p.Email, Password: p.Password}

	dbCon, dbErr := database.GetDatabaseConnection()

	if dbErr != nil {
		log.Error(dbErr)
		return
	}

	result := dbCon.Create(&userModel)

	if result.Error != nil {
		log.Error(result.Error)
		apiUtils.InternalErrorHandler(w)
	}
}
