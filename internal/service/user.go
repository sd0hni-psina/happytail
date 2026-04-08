package service

type UserService struct {
	repo      UserRepository
	jwtSecret string
}

func NewUserService(repo UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}
