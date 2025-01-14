//-----------------------------------------------------------------------------
// Copyright (c) 2020 Detlef Stern
//
// This file is part of zettelstore.
//
// Zettelstore is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// Zettelstore is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Zettelstore. If not, see <http://www.gnu.org/licenses/>.
//-----------------------------------------------------------------------------

package cmd

import (
	"flag"
	"sort"

	"zettelstore.de/z/domain"
)

// Command stores information about commands / sub-commands.
type Command struct {
	Name  string              // command name as it appears on the command line
	Func  CommandFunc         // function that executes a command
	Flags func(*flag.FlagSet) // function to set up flag.FlagSet
	flags *flag.FlagSet       // flags that belong to the command
}

// CommandFunc is the function that executes the command.
// It accepts meta data as configuration data and returns the exit code and an
// error.
type CommandFunc func(*domain.Meta) (int, error)

// GetFlags return the flag.FlagSet defined for the command.
func (c *Command) GetFlags() *flag.FlagSet { return c.flags }

var commands = make(map[string]Command)

// RegisterCommand registers the given command.
func RegisterCommand(cmd Command) {
	if cmd.Name == "" || cmd.Func == nil {
		panic("Required command values missing")
	}
	if _, ok := commands[cmd.Name]; ok {
		panic("Command already registered: " + cmd.Name)
	}
	cmd.flags = flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	if cmd.Flags != nil {
		cmd.Flags(cmd.flags)
	}
	commands[cmd.Name] = cmd
}

// Get returns the command identified by the given name and a bool to signal success.
func Get(name string) (Command, bool) {
	cmd, ok := commands[name]
	return cmd, ok
}

// List returns a sorted list of all registered command names.
func List() []string {
	result := make([]string, 0, len(commands))
	for name := range commands {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
