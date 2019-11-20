package vcsurl

import (
	"fmt"
	"strings"
)

var protocols = [...]string{
	"git",
	"https",
	"ssh",
}

type Repository struct {
	Host  string
	Owner string
	Name  string
}

func (this *Repository) ToHTTPS() string {
	return fmt.Sprintf("https://%s/%s/%s.git", this.Host, this.Owner, this.Name)
}

func (this *Repository) ToSSH() string {
	return fmt.Sprintf("git@%s:%s/%s.git", this.Host, this.Owner, this.Name)
}

func Parse(input string) (*Repository, error) {

	var repo Repository

	if strings.HasSuffix(input, ".git") {
		input = input[:len(input)-4]
	}

	if strings.HasPrefix(input, "git@") {

		input = input[4:]
		repo.Host = strings.Split(input, ":")[0]
		input = input[len(repo.Host)+1:]
	}

	if strings.HasPrefix(input, "https://") {

		input = input[8:]
		repo.Host = strings.Split(input, "/")[0]
		input = input[len(repo.Host)+1:]
	}

	repoParts := strings.Split(input, "/")
	repo.Owner = repoParts[0]
	repo.Name = repoParts[1]

	return &repo, nil
}

func Validate(url string) bool {
	return false
}
