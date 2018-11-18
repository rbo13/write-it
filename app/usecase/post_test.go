package usecase_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

	wantJSONResponse := `{"id":1,"creator_id":1,"post_title":"Test Post Title","post_body":"Test Post Body","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":null}`
	gotJSONResponse := response.Body.String()

	wantResponse := &app.Post{}
	gotResponse := &app.Post{}

	if err := json.NewDecoder(strings.NewReader(wantJSONResponse)).Decode(wantResponse); err != nil {
		log.Fatalln(err)
	}

	if err := json.NewDecoder(strings.NewReader(gotJSONResponse)).Decode(gotResponse); err != nil {
		log.Fatalln(err)
	}

	if !reflect.DeepEqual(gotResponse, wantResponse) {
		t.Errorf("Want: %s, but got: %s instead.\n", wantResponse.String(), gotResponse.String())
	}
}
