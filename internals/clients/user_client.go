package client

import (
	"fmt"
	"net/http"
	"user-favorites-service/internals/models"
)

type UserClient struct {
	URL string
}

func NewUserClient(host string) *UserClient {
	URL := host
	return &UserClient{URL: URL}
}

func (userClient *UserClient) CheckUserExist(userID int) (bool, error) {
	url := fmt.Sprintf("%s/users/%d", userClient.URL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, models.ErrorUserNotFound
	}

	return true, nil
}
