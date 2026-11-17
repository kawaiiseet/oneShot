package server

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type PersonalDataRepository interface {
	Create(ctx context.Context, p *PersonalData) error
	FindByLogin(ctx context.Context, login string) (*PersonalData, error)
	// FindMatch возвращает другого пользователя из той же страны с хотя бы одним общим интересом.
	FindMatch(ctx context.Context, country string, interests []string, login string) (*PersonalData, error)
}
