package service

type AnimalService struct {
	repo AnimalRepository
}

func NewAnimalService(repo AnimalRepository) *AnimalService {
	return &AnimalService{repo: repo}
}
