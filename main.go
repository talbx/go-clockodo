/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/talbx/go-clockodo/cmd"
	"github.com/talbx/go-clockodo/pkg/util"
)

func main() {
	util.CreateSugaredLogger()
	cmd.Execute()
}
