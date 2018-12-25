package inmemory_test

import (
	"testing"

	"github.com/rbo13/write-it/app"
	"github.com/rbo13/write-it/app/persistence/inmemory"
)

func TestInMemoryStore(t *testing.T) {

	postInmemory := inmemory.NewInMemoryPostService()

	t.Run("TestInMemoryCreatePost", func(t *testing.T) {
		post := &app.Post{
			ID:        int64(1),
			CreatorID: int64(1),
			PostTitle: "Test Post Title",
			PostBody:  "Test Post Body",
		}

		if err := postInmemory.CreatePost(post); err != nil {
			t.Errorf("Error occurred due to: %v", err)
		}

		post2 := &app.Post{
			ID:        int64(2),
			CreatorID: int64(1),
			PostTitle: "Test Post Title 2",
			PostBody:  "Test Post Body 2",
		}

		if err := postInmemory.CreatePost(post2); err != nil {
			t.Errorf("Error occurred due to: %v", err)
		}
	})

	t.Run("TestInMemoryGetPost", func(t *testing.T) {
		postID := int64(1)

		gotPost, err := postInmemory.Post(postID)

		if err != nil {
			t.Errorf("Error due to: %v", err)
		}

		if gotPost == nil {
			t.Errorf("Expecting: %v, but got: %v instead", gotPost, err)
		}

		t.Log(gotPost)
	})

	t.Run("TestInMemoryUpdatePost", func(t *testing.T) {
		post := &app.Post{
			ID:        int64(1),
			CreatorID: int64(1),
			PostTitle: "Test Update Post Title",
			PostBody:  "Test Update Post Body",
		}

		err := postInmemory.UpdatePost(post)

		if err != nil {
			t.Errorf("Error due to: %v", err)
		}
	})

	t.Run("TestInMemoryDeletePost", func(t *testing.T) {
		postID := int64(1)

		err := postInmemory.DeletePost(postID)

		if err != nil {
			t.Errorf("Error due to: %v", err)
		}
	})

	t.Run("TestInMemoryGetPosts", func(t *testing.T) {
		gotPosts, err := postInmemory.Posts()

		if err != nil {
			t.Errorf("Error due to: %v", err)
		}

		if gotPosts == nil {
			t.Errorf("Expecting: %v, but got: %v instead", gotPosts, err)
		}

		t.Log(gotPosts)
	})
}
