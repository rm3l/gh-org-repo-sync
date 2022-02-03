package github

import "github.com/cli/go-gh"

func GetUser() (string, error) {
	client, err := gh.RESTClient(nil)
	if err != nil {
		return "", err
	}
	var response struct{ Login string }
	err = client.Get("user", &response)
	if err != nil {
		return "", nil
	}
	return response.Login, nil
}
