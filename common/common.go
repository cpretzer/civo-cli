package common

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/tcnksm/go-latest"
)

var (
	// OutputFields for custom format output
	OutputFields string
	// OutputFormat for custom format output
	OutputFormat string
	// RegionSet picks the region to connect to, if you use this option it will use it over the default region
	RegionSet string
	// DefaultYes : automatic yes to prompts; assume \"yes\" as answer to all prompts and run non-interactively
	DefaultYes bool
	// PrettySet : Prints the json output in pretty format
	PrettySet bool
	// VersionCli is set from outside using ldflags
	VersionCli = "0.0.0"
	// CommitCli is set from outside using ldflags
	CommitCli = "none"
	// DateCli is set from outside using ldflags
	DateCli = "unknown"
)

// VersionCheck checks if there is a new version of the CLI
func VersionCheck() (res *latest.CheckResponse, skip bool) {
	githubTag := &latest.GithubTag{
		Owner:             "civo",
		Repository:        "cli",
		FixVersionStrFunc: latest.DeleteFrontV(),
	}
	res, err := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))
	if err != nil {
		if IsGHError(err) != nil {
			return nil, true
		}
		fmt.Printf("Checking for a newer version failed with %s \n", err)
		return nil, true
	}
	return res, false
}

//IsGHError checks if any error from github is returned
func IsGHError(err error) error {
	ghErr, ok := err.(*github.ErrorResponse)
	if ok {
		if ghErr.Response.StatusCode >= 400 && ghErr.Response.StatusCode < 500 {
			return errors.Wrap(err, `Failed to query the GitHub API for updates.`)
		}
		if ghErr.Response.StatusCode == http.StatusUnauthorized {
			return errors.Wrap(err, "Your Github token is invalid. Check the [github] section in ~/.gitconfig\n")
		} else {
			return errors.Wrap(err, "error finding latest release")
		}
	}
	return nil
}
