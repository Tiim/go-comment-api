package event

import (
	"fmt"
	"io"
	"os"
	"tiim/go-comment-api/model"
	"time"
)

type Store struct {
	store         model.Store
	eventhandlers []Handler
}

func NewEventStore(store model.Store, eventhandlers []Handler) *Store {
	return &Store{store: store, eventhandlers: eventhandlers}
}

func (s *Store) NewComment(c *model.Comment) error {
	for _, h := range s.eventhandlers {
		next, err := h.OnNewComment(c)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] on new comment %s: %s", h.Name(), c.Id, err)
		}
		if !next {
			return nil
		}
	}
	return s.store.NewComment(c)
}

func (s *Store) GetAllComments(since time.Time) ([]model.Comment, error) {
	return s.store.GetAllComments(since)
}

func (s *Store) GetCommentsForPost(page string, since time.Time) ([]model.Comment, error) {
	return s.store.GetCommentsForPost(page, since)
}

func (s *Store) DeleteComment(id string) error {
	comment, err := s.store.GetComment(id)

	if err != nil {
		return fmt.Errorf("error getting comment %s: %w", id, err)
	}

	for _, h := range s.eventhandlers {
		next, err := h.OnDeleteComment(&comment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] on delete comment %s: %s", h.Name(), id, err)
		}
		if !next {
			return nil
		}
	}
	return s.store.DeleteComment(id)
}

func (s *Store) GetComment(id string) (model.Comment, error) {
	return s.store.GetComment(id)
}

func (s *Store) Unsubscribe(secret string) (model.Comment, error) {
	v, ok := s.store.(model.SubscribtionStore)
	if ok {
		return v.Unsubscribe(secret)
	}
	return model.Comment{}, fmt.Errorf("store is not a subscribtion store")
}

func (s *Store) UnsubscribeAll(email string) ([]model.Comment, error) {
	v, ok := s.store.(model.SubscribtionStore)
	if ok {
		return v.UnsubscribeAll(email)
	}
	return nil, fmt.Errorf("store is not a subscribtion store")
}

func (s *Store) CleanUp() error {
	v, ok := s.store.(model.CleanupStore)
	if ok {
		return v.CleanUp()
	}
	return fmt.Errorf("store is not a cleanup store")
}

func (s *Store) Backup() (io.Reader, error) {
	v, ok := s.store.(model.BackupStore)
	if ok {
		return v.Backup()
	}
	return nil, fmt.Errorf("store is not a backup store")
}
