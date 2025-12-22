package services

import (
	"errors"
	"sync"

	"go-microservice/models"
)

// UserService manages user data in-memory.
type UserService struct {
	mu     sync.RWMutex
	users  map[int]models.User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]models.User),
		nextID: 1,
	}
}

func (s *UserService) Create(user models.User) models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = s.nextID
	s.nextID++
	s.users[user.ID] = user
	return user
}

func (s *UserService) GetAll() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result
}

func (s *UserService) GetByID(id int) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Update(id int, updated models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	updated.ID = id
	s.users[id] = updated
	return updated, nil
}

func (s *UserService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(s.users, id)
	return nil
}
