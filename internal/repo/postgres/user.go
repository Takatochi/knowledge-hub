package postgres

type UserRepo struct {
	store *Repository
}

func (u UserRepo) CreateUser() error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) GetUserByID(id uint) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) UpdateUser() error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepo) DeleteUser(id int) error {
	//TODO implement me
	panic("implement me")
}
