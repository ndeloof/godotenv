package godotenv

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func Load(filenames ...string) (err error) {
	for _, filename := range filenames {
		err = loadFile(filename)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}

func loadFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	bufferSize := 20
	lines := make([]string, bufferSize)
	lineReader := bufio.NewReaderSize(file, bufferSize)
	for line, isPrefix, e := lineReader.ReadLine(); e == nil; line, isPrefix, e = lineReader.ReadLine() {
		fullLine := string(line)
		if isPrefix {
			for {
				line, isPrefix, _ = lineReader.ReadLine()
				fullLine += string(line)
				if !isPrefix {
					break
				}
			}
		}
		// add a line to the game/parse
		lines = append(lines, string(line))
	}

	for _, fullLine := range lines {
		key, value, err := parseLine(fullLine)

		if err == nil {
			os.Setenv(key, value)
		}
	}

	return
}

func parseLine(line string) (key string, value string, err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")
		return
	}

	splitString := strings.Split(line, "=")

	if len(splitString) != 2 {
		// try yaml mode!
		splitString = strings.Split(line, ":")
	}

	if len(splitString) != 2 {
		err = errors.New("Can't separate key from value")
		return
	}

	key = splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.Trim(key, " ")

	value = strings.Trim(splitString[1], " \"'")
	value = strings.Replace(value, "\\\"", "\"", -1)

	return
}