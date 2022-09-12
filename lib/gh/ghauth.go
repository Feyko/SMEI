package gh

import (
	"context"
	"fmt"
	"gg-scm.io/pkg/ghdevice"
	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func AuthedClient(ctx context.Context) (*github.Client, error) {
	token, err := GetToken()
	if err != nil {
		return nil, errors.Wrap(err, "could not get an auth token")
	}

	return makeGithubClient(ctx, token), nil
}

func GetToken() (string, error) {
	opt := ghdevice.Options{
		ClientID: "0e4260b720ae65240864",
		Prompter: prompter,
		Scopes:   []string{"repo"},
	}
	return ghdevice.Flow(context.Background(), opt)
}

func prompter(ctx context.Context, prompt ghdevice.Prompt) error {
	fmt.Printf("Please navigate to %v and enter the following code: %v\n", prompt.VerificationURL, prompt.UserCode)
	return nil
}

func makeGithubClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
