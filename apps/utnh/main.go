/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/ffumaneri/github-cli/cmd"
	"github.com/ffumaneri/github-cli/ioc"
)

func main() {
	appContainer := &ioc.AppContainer{}
	cmd.Execute(appContainer)
}
