package vm

import (
	"log"
	"strconv"
	"strings"

	"arc/database"
)

type (
	// Runtime defines a virtual machine environment to run commands.
	Runtime struct {
		db           *database.Database
		library      Library
		libraryCache map[string]*LibraryFunction
	}
)

func getFunctionKey(command string, numberOfParameters int) string {
	if numberOfParameters <= 0 {
		return strings.ToUpper(command)
	} else {
		return strings.ToUpper(command) + "_" + strconv.FormatInt(int64(numberOfParameters), 10)
	}
}

func createLibraryCache(library Library) (cache map[string]*LibraryFunction) {
	var size = len(library)
	var functionKey string

	cache = make(map[string]*LibraryFunction, size)

	for index := 0; index < size; index++ {
		var function = library[index]
		functionKey = getFunctionKey(function.command, function.numberOfParameters)
		log.Printf("RTM: cached function %s", functionKey)
		cache[functionKey] = &function
	}

	return
}

// CreateRuntime creates a new runtime to run with the specified command library on the specified database.
func CreateRuntime(library Library, db *database.Database) (runtime *Runtime) {
	return &Runtime{
		db:           db,
		library:      library,
		libraryCache: createLibraryCache(library),
	}
}

// Execute executes a database command line and returns the result set (if any).
func (runtime *Runtime) Execute(line string) []string {
	var cmd = runtime.parseCommand(line)

	if cmd == nil {
		return invlaidCommandLineResult
	}

	var exists bool
	var function *LibraryFunction
	var functionKey = getFunctionKey(cmd.identifier, len(cmd.parameters))

	if function, exists = runtime.libraryCache[functionKey]; !exists {
		return unknownCommandResult
	}

	return function.call(runtime.db, cmd.parameters)
}

func (runtime *Runtime) parseCommand(line string) (cmd *command) {
	cmd = &command{
		identifier: "",
		parameters: make([]string, 0),
	}

	// FIXME: this is too hacky, but works for now: add a space at the end to make sure we interpret the last value found
	line = line + " "

	var inString bool
	var escapeChar bool
	var identifierOK bool
	var currentValue string

	for _, currentChar := range line {
		if escapeChar {
			if !identifierOK {
				return nil
			}

			currentValue += string(currentChar)
			continue
		}

		switch currentChar {
		case '"':
			if !identifierOK {
				return nil
			}

			if inString {
				cmd.parameters = append(cmd.parameters, currentValue)
				currentValue = ""
			}

			inString = !inString
		case ' ':
			if inString {
				currentValue += string(currentChar)
				continue
			}

			if currentValue == "" {
				continue
			}

			if !identifierOK {
				cmd.identifier = strings.ToUpper(currentValue)
				identifierOK = true
			} else {
				cmd.parameters = append(cmd.parameters, currentValue)
			}

			currentValue = ""
		case '\\':
			if !inString {
				return nil
			}

			escapeChar = true
		default:
			currentValue += string(currentChar)
		}
	}

	if inString || !identifierOK {
		return nil
	}

	return
}
