package gh

import (
	"context"
	"github.com/satisfactorymodding/SMEI/config"
	"github.com/satisfactorymodding/SMEI/lib/cfmt"
	"github.com/satisfactorymodding/SMEI/lib/secret"

	"gg-scm.io/pkg/ghdevice"
	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var accessToken secret.String

func AuthedClient(ctx context.Context) (*github.Client, error) {
	token, err := GetToken()
	if err != nil {
		return nil, errors.Wrap(err, "could not get an auth accessToken")
	}

	return makeGithubClient(ctx, string(token)), nil
}

func GetToken() (secret.String, error) {
	if accessToken != "" {
		return accessToken, nil
	}

	token, err := config.GetSecretString(config.GHToken_key)
	if err != nil {
		return "", errors.Wrap(err, "error getting the gh token")
	}

	if token != "" {
		err = saveToken(token)
		if err != nil {
			return "", errors.Wrap(err, "error saving token")
		}
	}

	opt := ghdevice.Options{
		ClientID: viper.GetString("GH-client-id"),
		Prompter: prompter,
		Scopes:   []string{"repo"},
	}
	newToken, err := ghdevice.Flow(context.Background(), opt)
	token = secret.String(newToken)
	err = saveToken(token)
	if err != nil {
		return "", errors.Wrap(err, "error saving token")
	}
	return accessToken, err
}

func prompter(ctx context.Context, prompt ghdevice.Prompt) error {
	cfmt.Sequence.Printf("Please navigate to %v and enter the following code: %v\nAfter authenticating, the download will begin here after a few seconds.\n", prompt.VerificationURL, prompt.UserCode)
	return nil
}

func makeGithubClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func saveToken(token secret.String) error {
	accessToken = token
	return config.SetSecretString(config.GHToken_key, token)
}
