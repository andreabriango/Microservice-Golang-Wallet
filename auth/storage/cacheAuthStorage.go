package storage

import (
	"errors"
	"log"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

//CheckExistingUserCache Check if the user Exits for Username or Email in Cache
func (s *AuthStorageService) CheckExistingUserCache(username string, email string) (bool, error) {
	if username == "" || len(username) <= 0 || email == "" || len(email) <= 0 {
		return false, errors.New("Review your Input")
	}

	val := s.rdb.Get(context.Background(), "RegisterUsername:"+username)

	if val.Err() != redis.Nil {
		return true, nil
	}

	val = s.rdb.Get(context.Background(), "RegisterEmail:"+email)

	if val.Err() != redis.Nil {
		return true, nil
	}

	return false, val.Err()
}

//SetRegisterCache Set in cache the user info
func (s *AuthStorageService) SetRegisterCache(username string, email string) {
	if username == "" || len(username) <= 0 || email == "" || len(email) <= 0 {
		log.Println("Review your Input")
		return
	}

	if err := s.rdb.Set(context.Background(), "RegisterUsername:"+username, "Exits", 0); err.Err() != nil {
		log.Println("Error in Auth Cache " + err.Err().Error())
	}
	if err := s.rdb.Set(context.Background(), "RegisterEmail:"+email, "Exits", 0); err.Err() != nil {
		log.Println("Error in Auth Cache " + err.Err().Error())
	}

	//Here must use User Service for set the user in cache
}

//ChangeRegisterCache Change in cache the username and email
func (s *AuthStorageService) ChangeRegisterCache(oldUsername string, newUsername string, oldEmail string, newEmail string) (bool, error) {
	if len(oldUsername) > 0 && len(newUsername) > 0 {
		status := s.rdb.Rename(context.Background(), "RegisterUsername:"+oldUsername, "RegisterUsername:"+newUsername)
		if status.Err() != nil {
			return false, status.Err()
		}
		return true, nil
	}

	if len(oldEmail) > 0 && len(newEmail) > 0 {
		status := s.rdb.Rename(context.Background(), "RegisterEmail:"+oldEmail, "RegisterEmail:"+newEmail)
		if status.Err() != nil {
			return false, status.Err()
		}
		return true, nil
	}

	return false, errors.New("Invalid Input")
}
