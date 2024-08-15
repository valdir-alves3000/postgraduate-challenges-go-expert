package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDBUser() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateUser(t *testing.T) {
	db, err := getDBUser()
	if err != nil {
		t.Error(err)
	}
	user, _ := entity.NewUser("John", "john@doe.com", "123456")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)
	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)
}

func TestFindByEmail(t *testing.T) {
	db, err := getDBUser()
	if err != nil {
		t.Error(err)
	}
	user, _ := entity.NewUser("John", "john@doe.com", "123456")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	userFound, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)
}
