package services

import (
	"errors"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
)

func Register(userdto dto.RegisterDto) error {
	if userdto.Username == "" || userdto.Password == "" {
		return errors.New(models.UsernameOrPasswordCannotBeEmptyMessage)
	}

	userdto.Password = pkg.MD5Encrypt(userdto.Password)

	newuser := models.User{
		Username: userdto.Username,
		Password: userdto.Password,
		IsAdmin:  false,
	}

	user, err := repository.GetUserByUsername(userdto.Username)
	if err == nil && user != nil && user.ID != 0 {
		return errors.New(models.UsernameAlreadyExistsMessage)
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		return errors.New(models.GetAllUsersFailMessage)
	}
	if len(users) == 0 {
		newuser.IsAdmin = true
	}

	if err := repository.CreateUser(&newuser); err != nil {
		return errors.New(models.CreateUserFailMessage)
	}

	return nil
}

func Login(userdto dto.LoginDto) (string, error) {
	if userdto.Username == "" || userdto.Password == "" {
		return "", errors.New(models.UsernameOrPasswordCannotBeEmptyMessage)
	}

	userdto.Password = pkg.MD5Encrypt(userdto.Password)

	user, err := repository.GetUserByUsername(userdto.Username)
	if err != nil {
		return "", errors.New(models.UserNotFoundMessage)
	}

	if user.Password != userdto.Password {
		return "", errors.New(models.PasswordIncorrectMessage)
	}

	token, err := pkg.GenerateToken(pkg.CreateClaims(*user))
	if err != nil {
		return "", errors.New(models.GenerateTokenFailMessage)
	}

	return token, nil
}

func GetStatus() (models.Status, error) {
	sysuser, err := repository.GetSysAdmin()
	if err != nil {
		return models.Status{}, errors.New(models.UserNotFoundMessage)
	}

	var users []models.UserStatus
	allusers, err := repository.GetAllUsers()
	if err != nil {
		return models.Status{}, errors.New(models.GetAllUsersFailMessage)
	}
	for _, user := range allusers {
		users = append(users, models.UserStatus{
			UserID:   user.ID,
			UserName: user.Username,
			IsAdmin:  user.IsAdmin,
		})
	}

	status := models.Status{}

	messages, err := repository.GetAllMessages(true)
	if err != nil {
		return status, errors.New(models.GetAllMessagesFailMessage)
	}

	status.SysAdminID = sysuser.ID
	status.Username = sysuser.Username
	status.Users = users
	status.TotalMessages = len(messages)

	return status, nil
}

func GetUserByID(userID uint) (*models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New(models.UserNotFoundMessage)
	}
	return user, nil
}

func IsUserAdmin(userID uint) (bool, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return false, errors.New(models.UserNotFoundMessage)
	}
	return user.IsAdmin, nil
}

func UpdateUser(user *models.User, userdto dto.UserInfoDto) error {
	if user.Username == userdto.Username {
		return nil
	}

	if userdto.Username == "" {
		return errors.New(models.UsernameCannotBeEmptyMessage)
	}

	user.Username = userdto.Username
	if err := repository.UpdateUser(user); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func ChangePassword(user *models.User, userdto dto.UserInfoDto) error {
	if user.Password == pkg.MD5Encrypt(userdto.Password) {
		return errors.New(models.PasswordCannotBeSameAsBeforeMessage)
	}

	if userdto.Password == "" {
		return errors.New(models.PasswordCannotBeEmptyMessage)
	}

	user.Password = pkg.MD5Encrypt(userdto.Password)
	if err := repository.UpdateUser(user); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func UpdateUserAdmin(userID uint) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.IsAdmin = !user.IsAdmin

	if err := repository.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func UpdateSetting(settingID uint, updates map[string]interface{}) error {
	// 实现更新逻辑
	// ...
	return nil
}