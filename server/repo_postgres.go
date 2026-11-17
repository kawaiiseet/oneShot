package server

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

func (r *PostgresRepo) Create(ctx context.Context, p *PersonalData) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO personal_data (name, login, age, country, interests, likes, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, p.Name, p.Login, p.Age, p.Country, p.Interests, p.Likes, passwordHashing(p.Password)).Scan(&p.ID)
}

func (r *PostgresRepo) FindByLogin(ctx context.Context, login string) (*PersonalData, error) {
	var p PersonalData
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, login, age, country, interests, likes
		FROM personal_data
		WHERE login = $1
	`, login).Scan(&p.ID, &p.Name, &p.Login, &p.Age, &p.Country, &p.Interests, &p.Likes)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresRepo) FindMatch(
	ctx context.Context,
	country string,
	interests []string,
	excludeLogin string,
) (*PersonalData, error) {
	var p PersonalData
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, login, age, country, interests, likes
		FROM personal_data
		WHERE country = $1
		  AND interests && $2::text[]          -- хотя бы один общий интерес
		  AND $3 = ANY(likes)
		  AND ($3 = '' OR login <> $3)         -- не возвращаем самого себя
		ORDER BY (
			SELECT count(*)
			FROM unnest(interests) AS i
			WHERE i = ANY($2::text[])
		) DESC                                 -- больше общих интересов — выше
		LIMIT 1
	`, country, interests, excludeLogin).Scan(
		&p.ID, &p.Name, &p.Login, &p.Age, &p.Country, &p.Interests, &p.Likes,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
