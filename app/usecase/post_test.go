package usecase_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi"
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

	request, _ := http.NewRequest("POST", "/api/v1/posts/create", bytes.NewBuffer(jsonPost))
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

func TestGetRouter(t *testing.T) {
	inmemory := inmemory.NewInMemoryPostService()
	postUsecase := usecase.NewPost(inmemory)

	request, _ := http.NewRequest("GET", "/api/v1/posts", nil)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(postUsecase.Get)

	handler.ServeHTTP(response, request)

	// Check the status code is what we expect.
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestGetByIDRouter(t *testing.T) {
	// inmemory := inmemory.NewInMemoryPostService()
	// postUsecase := usecase.NewPost(inmemory)

	m := chi.NewRouter()
	//
	// request, _ := http.NewRequest("GET", "/api/v1/posts/1", nil)
	// response := httptest.NewRecorder()
	//
	// handler := http.HandlerFunc(postUsecase.GetByID)
	//
	// handler.ServeHTTP(response, request)
	m.Get("/post/{id}", pingOne)
	ts := httptest.NewServer(m)
	defer ts.Close()

	if _, body := testRequest(t, ts, "GET", "/post/123", nil); body != "ping one id: 123" {
		t.Fatalf(body)
	}
}
func pingOne(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("ping one id: %s", idParam)))
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
