package doctor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("doctor not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, page, pageSize int) (doctors []Doctor, total int, err error) {
	return s.repo.List(ctx, ListParams{Page: page, PageSize: pageSize})
}

func (s *Service) GetByID(ctx context.Context, id int64) (Doctor, error) {
	d, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Doctor{}, ErrNotFound
	}
	return d, err
}
