package postgres

import (
	"KnowledgeHub/internal/repo"
	"KnowledgeHub/pkg/postgres"
)

type Repository struct {
	db             *postgres.Postgres
	userRepository *UserRepo
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) User() repo.UserRepository {
	if r.userRepository != nil {
		return r.userRepository
	}

	r.userRepository = &UserRepo{
		store: r,
	}

	return r.userRepository
}

//... other
