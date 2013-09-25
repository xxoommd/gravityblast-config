package config

import (
  "os"
  "io"
  "bufio"
  "strings"
  "regexp"
)

type Options map[string]string

type Sections map[string]Options

func parse(reader *bufio.Reader, mainSectionName string) (Sections, error) {
  sections    := make(Sections)
  section     := mainSectionName
  options     := make(Options)
  splitRegexp := regexp.MustCompile(`\s+`)

  for {
    line, err := reader.ReadString('\n')
    if err != nil && err == io.EOF {
      break
    } else if err != nil {
      return sections, err
    }

    line = strings.TrimSpace(line)

    if len(line) == 0 || line[0] == '#' || line[0] == ';' {
      continue
    }

    if line[0] == '[' && line[len(line) - 1] == ']' {
      sections[section] = options
      options = make(Options)
      section = line[1:(len(line) - 1)]
    } else {
      values := splitRegexp.Split(line, 2)
      if len(values) == 2 {
        options[values[0]] = values[1]
      }
    }
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