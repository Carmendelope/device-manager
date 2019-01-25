/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/device-manager/cmd/device-manager/commands"
	"github.com/nalej/device-manager/version"
)

var MainVersion string

var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
