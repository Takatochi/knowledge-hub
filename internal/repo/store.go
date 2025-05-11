package repo

// Repository implement from interface Store
type Store interface {
	User() UserRepository
	//... other entity
}

type UserRepository interface {
	CreateUser() error
	GetUserByID(id uint) error
	UpdateUser() error
	DeleteUser(id int) error
}
