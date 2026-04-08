package service

type ShelterService struct {
	repo ShelterRepository
}

func NewShelterService(repo ShelterRepository) *ShelterService {
	return &ShelterService{repo: repo}
}
