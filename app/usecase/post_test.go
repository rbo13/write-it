package usecase_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbo13/write-it/app"
	"github.com/rbo13/write-it/app/inmemory"
	"github.com/rbo13/write-it/app/usecase"
)

func TestCreateRouter(t *testing.T) {

	inmemory := inmemory.NewInMemoryPostService()
	postUsecase := usecase.NewPost(inmemory)

	post := &app.Post{
		ID:        int64(1),
		CreatorID: int64(1),
		PostTitle: "Test Post Title",
		PostBody:  "Test Post Body",
	}

	jsonPost, _ := json.Marshal(post)

	request, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonPost))
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(postUsecase.Create)

	handler.ServeHTTP(response, request)

	// Check the status code is what we expect.
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
