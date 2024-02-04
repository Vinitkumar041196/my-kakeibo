package repository

import (
	"fmt"
	"kakeibo-api/db"
	"kakeibo-api/domain"
	"kakeibo-api/utils"
	"strings"
	"time"
)

type userRepository struct {
	db *db.PostgresDB
}

func NewUserRepository(db *db.PostgresDB) domain.UserRepository {
	fmt.Println("Automigrate error:", db.GormDB.AutoMigrate(&domain.User{}))
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(user *domain.User) error {
	createdAt := utils.TimeNow().Format(time.RFC3339)
	isActive := true
	var err error

	if user.Password == "" {
		return fmt.Errorf("empty password")
	}

	hash, err := utils.GetHash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash

	user.MapMetaData(&createdAt, &createdAt, &isActive)

	tx := ur.db.GormDB.Begin()
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (ur *userRepository) UpdateUser(user *domain.User) error {
	updatedAt := utils.TimeNow().Format(time.RFC3339)
	isActive := true
	var err error

	if user.Password != "" {
		hash, err := utils.GetHash(user.Password)
		if err != nil {
			return err
		}
		user.Password = hash
	}

	user.MapMetaData(nil, &updatedAt, &isActive)

	fmt.Printf("%+v", user)
	tx := ur.db.GormDB.Begin()
	err = tx.Updates(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (ur *userRepository) GetUser(email string) (*domain.User, error) {
	user := domain.User{}
	err := ur.db.GormDB.Where("is_active = ?", true).Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) SearchUserByKey(key, value string) (*domain.User, error) {
	user := domain.User{}
	err := ur.db.GormDB.Where("is_active = ?", true).Where(fmt.Sprintf("lower(%s) = '%v'", key, strings.ToLower(value))).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) DeleteUser(email string) error {
	res := ur.db.GormDB.
		Table((domain.User{}).TableName()).
		Where("email = ?", email).
		Where("is_active = ?", true).
		Updates(map[string]interface{}{"is_active": false, "updated_at": utils.TimeNow().Format(time.RFC3339)})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}
