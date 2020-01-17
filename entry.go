package timetracking

import "context"

type Entry struct {
	ID       int    `json:"id"`
	DateTS   int64  `json:"date_ts"`
	Project  int    `json:"project"`
	User     string `json:"user"`
	Comment  string `json:"comment"`
	Duration int64  `json:"duration"`
}

type Repository interface {
	Add(ctx context.Context, e *Entry) (*Entry, error)
	Get(ctx context.Context, id int) (*Entry, error)
	Update(ctx context.Context, e *Entry) error
	Delete(ctx context.Context, id int) error
	All(ctx context.Context) ([]*Entry, error)
}
