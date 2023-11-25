package vm

import (
	"arc/database"
)

type (
	command struct {
		identifier string
		parameters []string
	}

	// Function defines the virtual machine library function interface.
	Function func(db *database.Database, parameters []string) []string

	// LibraryFunction holds the needed information for a library function to work on runtime.
	LibraryFunction struct {
		command            string
		numberOfParameters int
		call               Function
		help               string
	}

	// Library represents a set of library functions.
	Library []LibraryFunction
)

// StandardLibrary defines the standard function library.
var StandardLibrary = Library{
	{command: "SET", numberOfParameters: 2, call: stdSet, help: "SET key value"},
	{command: "SET", numberOfParameters: 4, call: stdSetEx, help: "SET key value EX seconds"},
	{command: "GET", numberOfParameters: 1, call: stdGet, help: "GET key"},
	{command: "DEL", numberOfParameters: -1, call: stdDel, help: "DEL key [key...]"},
	{command: "DBSIZE", numberOfParameters: 0, call: stdDbSize, help: "DBSIZE"},
	{command: "INCR", numberOfParameters: 1, call: stdIncr, help: "INCR key"},
	{command: "ZADD", numberOfParameters: -1, call: stdZadd, help: "ZADD key score member [score member...]"},
	{command: "ZCARD", numberOfParameters: 1, call: stdZcard, help: "ZCARD key"},
	{command: "ZRANK", numberOfParameters: 2, call: stdZrank, help: "ZRANK key member"},
	{command: "ZRANGE", numberOfParameters: 3, call: stdZrange, help: "ZRANGE key start stop"},
}

const (
	nilMessage                        = "(nil)"
	okMessage                         = "OK"
	unknownCommandErrorMessage        = "Error: unknown command or invalid parameters for command"
	invalidCommandLineErrorMessage    = "Error: invalid command line"
	invalidParametersErrorMessage     = "Error: invalid parameters"
	invalidParameterValueErrorMessage = "Error: invalid parameter value"
	invalidDataTypeErrorMessage       = "Error: invalid data type"
)

var (
	emptyResult                 = []string{}
	nilResult                   = []string{nilMessage}
	okResult                    = []string{okMessage}
	unknownCommandResult        = []string{unknownCommandErrorMessage}
	invlaidCommandLineResult    = []string{invalidCommandLineErrorMessage}
	invalidParametersResult     = []string{invalidParametersErrorMessage}
	invalidParameterValueResult = []string{invalidParameterValueErrorMessage}
	invalidDataTypeResult       = []string{invalidDataTypeErrorMessage}
)

// GetHelp returns the help string for the function.
func (function *LibraryFunction) GetHelp() string {
	return function.help
}
