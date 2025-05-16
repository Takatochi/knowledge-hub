package postgres

import "KnowledgeHub/internal/models"

type UserRepo struct {
	store *Repository
}

func (u UserRepo) CreateUser() error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) GetUserByID(_ uint) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) UpdateUser() error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) DeleteUser(_ int) error {
	//TODO implement me
	panic("implement me")
}
