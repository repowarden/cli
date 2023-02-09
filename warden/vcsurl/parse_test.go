package vcsurl

import (
	"testing"
)

func TestParse(t *testing.T) {

	tcs := []struct {
		url string
		org string
	}{
		{url: "http://github.com/felicianotech/sonar", org: "felicianotech"},
		{url: "https://github.com/felicianotech/sonar", org: "felicianotech"},
		{url: "https://github.com/felicianotech/sonar.git", org: "felicianotech"},
		{url: "git@github.com:felicianotech/sonar.git", org: "felicianotech"},
	}

	for i, tc := range tcs {

		repo, err := Parse(tc.url)
		if err != nil {
			t.Fatal("The URL should have parsed successfully but it didn't.")
		}

		if repo.Owner != tc.org {
			t.Errorf("URL %d: Want org name '%s', got '%s'", i+1, tc.org, repo.Owner)
		}
	}
}

func TestToHTTPS(t *testing.T) {

	tcs := []struct {
		start string
		end   string
	}{
		{
			start: "https://github.com/felicianotech/sonar.git",
			end:   "https://github.com/felicianotech/sonar",
		},
		{
			start: "https://github.com/felicianotech/sonar",
			end:   "https://github.com/felicianotech/sonar",
		},
		{
			start: "git@github.com:felicianotech/sonar.git",
			end:   "https://github.com/felicianotech/sonar",
		},
	}

	for i, tc := range tcs {

		repo, err := Parse(tc.start)
		if err != nil {
			t.Fatal("The URL should have parsed successfully but it didn't.")
		}

		if repo.ToHTTPS() != tc.end {
			t.Errorf("URL %d: Want the generated URL as '%s',\n got '%s'", i+1, tc.end, repo.ToHTTPS())
		}
	}
}

func TestToSSH(t *testing.T) {

	tcs := []struct {
		start string
		end   string
	}{
		{
			start: "https://github.com/felicianotech/sonar.git",
			end:   "git@github.com:felicianotech/sonar.git",
		},
		{
			start: "https://github.com/felicianotech/sonar",
			end:   "git@github.com:felicianotech/sonar.git",
		},
		{
			start: "git@github.com:felicianotech/sonar.git",
			end:   "git@github.com:felicianotech/sonar.git",
		},
	}

	for i, tc := range tcs {

		repo, err := Parse(tc.start)
		if err != nil {
			t.Fatal("The URL should have parsed successfully but it didn't.")
		}

		if repo.ToSSH() != tc.end {
			t.Errorf("URL %d: Want the generated URL as '%s',\n got '%s'", i+1, tc.end, repo.ToHTTPS())
		}
	}
}
