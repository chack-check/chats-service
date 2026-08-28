package entities

import (
	"fmt"
)

type User struct {
	id         int
	avatar     *SavedFile
	lastName   string
	firstName  string
	middleName *string
	username   string
}

func NewUser(
	id int,
	avatar *SavedFile,
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

func (user *User) GetID() int {
	return user.id
}

func (user *User) GetLastName() string {
	return user.lastName
}

func (user *User) GetFirstName() string {
	return user.firstName
}

func (user *User) GetMiddleName() *string {
	return user.middleName
}

func (user *User) GetUsername() string {
	return user.username
}

func (user *User) GetFullName() string {
	if user.middleName != nil {
		return fmt.Sprintf("%s %s %s", user.lastName, user.firstName, *user.middleName)
	} else {
		return fmt.Sprintf("%s %s", user.lastName, user.firstName)
	}
}

func (user *User) GetAvatar() *SavedFile {
	return user.avatar
}

type ActionUser struct {
	id         int
	lastName   string
	firstName  string
	middleName *string
	username   string
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

func (user *ActionUser) GetID() int {
	return user.id
}

func (user *ActionUser) GetFullName() string {
	if user.middleName != nil {
		return fmt.Sprintf("%s %s %s", user.lastName, user.firstName, *user.middleName)
	}

	return fmt.Sprintf("%s %s", user.lastName, user.firstName)
}

func (user *ActionUser) GetLastName() string {
	return user.lastName
}

func (user *ActionUser) GetFirstName() string {
	return user.firstName
}

func (user *ActionUser) GetMiddleName() *string {
	return user.middleName
}

func (user *ActionUser) GetUsername() string {
	return user.username
}
