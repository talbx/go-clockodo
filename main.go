/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/talbx/go-clockodo/cmd"
	"github.com/talbx/go-clockodo/pkg/util"
	"go.uber.org/zap/zapcore"
)

func main() {
	util.CreateSugaredLogger(zapcore.WarnLevel)
	cmd.Execute()
}
