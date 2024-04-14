package Database

import (
	"errors"
	"log"
	"messaging/Model"

	"gorm.io/gorm"
)

type userRepo struct{}

func UserRepo() *userRepo {
	return &userRepo{}
}

func (*userRepo) CreateUser(user Model.User) (Model.User, error) {
	if err := DBConnection.Create(&user).Error; err != nil {
		return Model.User{}, err
	}
	return user, nil
}

func (*userRepo) GetUser(name string) (Model.User, error) {
	// Query to retrieve user by name or email
	var user Model.User
	if err := DBConnection.Where("username = ?", name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Model.User{}, errors.New("no record found")
		} else {
			return Model.User{}, err
		}
	}

	log.Printf("User found: ID=%d, Username=%s\n", user.ID, user.Username)
	return user, nil
}
