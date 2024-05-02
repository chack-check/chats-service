package users

import (
	"fmt"

	"github.com/chack-check/chats-service/domain/files"
)

type User struct {
	id         int
	avatar     *files.SavedFile
	lastName   string
	firstName  string
	middleName *string
	username   string
}

func (model *User) GetId() int {
	return model.id
}

func (model *User) GetLastName() string {
	return model.lastName
}

func (model *User) GetFirstName() string {
	return model.firstName
}

func (model *User) GetMiddleName() *string {
	return model.middleName
}

func (model *User) GetUsername() string {
	return model.username
}

func (model *User) GetFullName() string {
	if model.middleName != nil {
		return fmt.Sprintf("%s %s %s", model.lastName, model.firstName, *model.middleName)
	} else {
		return fmt.Sprintf("%s %s", model.lastName, model.firstName)
	}
}

func (model *User) GetAvatar() *files.SavedFile {
	return model.avatar
}

type ActionUser struct {
	id         int
	lastName   string
	firstName  string
	middleName *string
	username   string
}

func (model *ActionUser) GetId() int {
	return model.id
}

func (model *ActionUser) GetFullName() string {
	if model.middleName != nil {
		return fmt.Sprintf("%s %s %s", model.lastName, model.firstName, *model.middleName)
	}

	return fmt.Sprintf("%s %s", model.lastName, model.firstName)
}

func (model *ActionUser) GetLastName() string {
	return model.lastName
}

func (model *ActionUser) GetFirstName() string {
	return model.firstName
}

func (model *ActionUser) GetMiddleName() *string {
	return model.middleName
}

func (model *ActionUser) GetUsername() string {
	return model.username
}

func NewActionUser(id int, lastName string, firstName string, middleName *string, username string) ActionUser {
	return ActionUser{
		id:         id,
		lastName:   lastName,
		firstName:  firstName,
		middleName: middleName,
		username:   username,
	}
}

func NewUser(
	id int,
	avatar *files.SavedFile,
	lastName string,
	firstName string,
	middleName *string,
	username string,
) User {
	return User{
		id:         id,
		avatar:     avatar,
		lastName:   lastName,
		firstName:  firstName,
		middleName: middleName,
		username:   username,
	}
}
