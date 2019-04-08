package app

import (
  "net/http"
)

// Handler is an interface that defines the basic operations of every http request.
type Handler interface {
  Create(w http.ResponseWriter, r *http.Request)
  Get(w http.ResponseWriter, r *http.Request)
  GetByID(w http.ResponseWriter, r *http.Request)
  Update(w http.ResponseWriter, r *http.Request)
  Delete(w http.ResponseWriter, r *http.Request)
}

// UserHandler implements the Handler interface with some user related methods.
type UserHandler interface {
  Handler

  Login(w http.ResponseWriter, r *http.Request)
  GetUserPosts(w http.ResponseWriter, r *http.Request)
}
