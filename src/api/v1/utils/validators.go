package utils

import (
	"fmt"

	"github.com/chack-check/chats-service/protousers"
)

func UserRequired(user *protousers.UserResponse) error {
	if user == nil || user.Id == 0 {
		return fmt.Errorf("Incorrect token")
	}

	return nil
}
