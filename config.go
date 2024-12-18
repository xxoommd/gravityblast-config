package config

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

type Options map[string]string

type Sections map[string]Options

var commentSplitRegexp = regexp.MustCompile(`[#;]`)

var keyValueSplitRegexp = regexp.MustCompile(`(\s*(:|=)\s*)|\s+`)

func cleanLine(line string) string {
	chunks := commentSplitRegexp.Split(line, 2)
	return strings.TrimSpace(chunks[0])
}

func parse(reader *bufio.Reader, mainSectionName string) (Sections, error) {
	sections := make(Sections)
	section := mainSectionName
	options := make(Options)

	processLine := func(line string) {
		line = cleanLine(line)
		if len(line) == 0 {
			return
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			sections[section] = options
			section = line[1:(len(line) - 1)]
			options = sections[section] // check if section already exists
			if options == nil {
				options = make(Options)
			}
		} else {
			values := keyValueSplitRegexp.Split(line, 2)
			key := values[0]
			value := ""
			if len(values) == 2 {
				value = values[1]
			}

			options[key] = value
		}
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					// Process the last line if it's not empty
					processLine(line)
				}
				break
			}
			return sections, err
		}

		processLine(line)
	}

	sections[section] = options

	return sections, nil
}

func ParseFile(path string, mainSectionName string) (Sections, error) {
	file, err := os.Open(path)
	if err != nil {
		return make(Sections), err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	return parse(reader, mainSectionName)
}
