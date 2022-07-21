package handlers

import (
	"cde/internal/entities"
	"cde/internal/storage"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Service struct {
	Storage      *storage.MongoStorage
	MongoStorage *storage.MongoStorage
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u entities.MongoUser
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		usedID, err := s.MongoStorage.Save(&u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		response := CrateResponse{ID: usedID}
		responseBody, _ := json.Marshal(response)
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBody)

		defer r.Body.Close()
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u MakeFriendsRequest
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		ans, err := s.MongoStorage.Add(u.SourceId, u.TargetId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ans))

		defer r.Body.Close()
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) DelieteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u DelieteUserRequest
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		ans, err := s.MongoStorage.DelieteU(u.TargetId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ans))

		defer r.Body.Close()
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parts := strings.Split(r.URL.Path, "/")
		userID := parts[len(parts)-1]

		friends, err := s.MongoStorage.GetU(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		responseBody, err := json.Marshal(friends)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)

		defer r.Body.Close()
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) NewAge(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		parts := strings.Split(r.URL.Path, "/")
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u NewAgeRequest
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		userID := parts[len(parts)-1]
		userAge, err := s.MongoStorage.NewAgeU(u.NewAge, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userAge))

		defer r.Body.Close()
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		user, err := s.MongoStorage.AllUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		responseBody, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)

		defer r.Body.Close()
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
