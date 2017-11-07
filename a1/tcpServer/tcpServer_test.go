package tcpServer

import (
	"436bin/a1/config"
	"436bin/a1/model"
	"fmt"
	"testing"
)

func TestUpdateUserConfigCorrectlyUpdatesConfig(t *testing.T) {
	userName := "foo"
	user := &model.User{UserName: userName, Config: &config.ClientConfig{MessagePort: 1234}}
	updatedConfig := &config.ClientConfig{MessagePort: 9080}
	users = []*model.User{
		user,
	}

	// Execute block
	updateUserConfig(userName, updatedConfig)

	// Assert block
	if user.Config.MessagePort != 9080 {
		fmt.Println("Failed to update the user's config")
		t.Fail()
	}
}

func TestUpdateUserConfigCorrectlyUpdatesConfigAfterCreatingUser(t *testing.T) {
	userName := "foo"
	updatedConfig := &config.ClientConfig{MessagePort: 9080}
	users = []*model.User{}

	// Execute block
	updateUserConfig(userName, updatedConfig)

	// Assert block
	user := getUser(userName)
	if user.Config.MessagePort != 9080 {
		fmt.Println("Failed to update the user's config")
		t.Fail()
	}
}
