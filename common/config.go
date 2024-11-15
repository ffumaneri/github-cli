package common

import (
	"fmt"
	"github.com/google/go-github/v65/github"
	"github.com/spf13/viper"
)

func GetConfig() (token string, owner string, ok bool) {
	ok = true
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	token, ok = viper.Get("TOKEN").(string)
	if !ok {
		return "", "", false
	}
	owner, ok = viper.Get("OWNER").(string)
	if !ok {
		return "", "", false
	}

	return
}

func GetClient() (*github.Client, string) {
	token, owner, ok := GetConfig()
	if !ok {
		return nil, ""
	}
	// Create Git hub client
	client := github.NewClient(nil).WithAuthToken(token)
	return client, owner
}
