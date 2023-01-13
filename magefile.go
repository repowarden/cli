//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Install() error {

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	return sh.Run("go", "install", "./warden")
}

func Remove() error {

	return sh.Run("go", "clean", "-i", "github.com/repowarden/cli/warden")
}

func Test() error {

	return sh.Run("gotestsum", "./...")
}

func TestCI() error {

	return sh.Run("gotestsum", "--junitfile=junit/unit-tests.xml", "--", "-coverprofile=coverage.txt", "-covermode=atomic", "./...")
}
