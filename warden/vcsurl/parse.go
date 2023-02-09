package vcsurl

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/exp/slices"
)

var hosts = []string{
	"github.com",
}

var protocols = [...]string{
	"git",
	"http",
	"https",
	"ssh",
}

type Repository struct {
	Host  string
	Owner string
	Name  string
}

func (this *Repository) ToHTTPS() string {
	return fmt.Sprintf("https://%s/%s/%s", this.Host, this.Owner, this.Name)
}

func (this *Repository) ToSSH() string {
	return fmt.Sprintf("git@%s:%s/%s.git", this.Host, this.Owner, this.Name)
}

func Parse(input string) (*Repository, error) {

	if strings.HasSuffix(input, ".git") {
		input = input[:len(input)-4]
	}

	if strings.HasPrefix(input, "git@") {
		input = strings.Replace(input, ":", "/", 1)
		input = strings.Replace(input, "git@", "https://", 1)
	}

	repoURL, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(hosts, repoURL.Host) {
		return nil, fmt.Errorf("%s is not a valid hostname.", repoURL.Host)
	}

	repoParts := strings.Split(repoURL.Path, "/")

	return &Repository{
		Host:  repoURL.Host,
		Owner: repoParts[1],
		Name:  repoParts[2],
	}, nil
}

func Validate(url string) bool {
	return false
}
