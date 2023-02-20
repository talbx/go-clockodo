package util

import "github.com/spf13/pflag"

var flagSet *pflag.FlagSet

func StoreFlags(flags *pflag.FlagSet){
	flagSet = flags
}

func GetFlags() *pflag.FlagSet{
	return flagSet
}