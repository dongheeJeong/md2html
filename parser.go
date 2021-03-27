package main

import (
  "fmt"
  "flag"
  "bufio"
  "os"
  "log"
  "strings"
  "regexp"
)

type Config struct {
  fileName string
}


func main() {
  config := parseCliArgs()

  file, err := os.Open(config.fileName)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  codeBlockOpen := false


  for scanner.Scan() {

    lineConverted := false

    line := scanner.Text()
    line = strings.TrimSpace(line)


    if strings.HasPrefix(line, "* ") {
      handleOrderedList()
      continue
    }

    if strings.HasPrefix(line, "- ") {
      handleUnorderedList()
      continue
    }


    // codeBlock
    if strings.HasPrefix(line, "```") {
      lineConverted = true

      if codeBlockOpen {
	line = "</code></pre>"
        fmt.Printf("%s\n", line)
      } else {
	line =  "<pre><code>"
        fmt.Printf("%s", line) // newline을 출력하면 안된다.
      }
      codeBlockOpen = !codeBlockOpen
      continue
    }

    if codeBlockOpen {
      fmt.Printf("%s\n", line)
      continue
    }



    line, lineConverted = handlePrefix(line)


    regex_detected := false
    for {
      regex_detected, line = regex(line)
      if regex_detected == false {
	break;
      }
    }


    if lineConverted == false && line != "" {
      line = "<p>" + line
      line = AppendCloseTag(line, "</p>")
    }
    fmt.Printf("%s\n", line)
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

}
func handlePrefix2(line string) (newLine string, converted bool) {

  newLine = line

  if strings.HasPrefix(newLine, "# ") {
  }

  if strings.HasPrefix(newLine, "## ") {
  }

  if strings.HasPrefix(newLine, "### ") {
  }

  if strings.HasPrefix(newLine, "#### ") {
  }

  if strings.HasPrefix(newLine, "##### ") {
  }

  if strings.HasPrefix(newLine, "> ") {
  }

  if strings.HasPrefix(newLine, "---") {
  }

  type Mapping struct {
    mdToken   string
    openTag   string
    closeTag  string
  }

  var mappings = []Mapping {
    Mapping {
      mdToken: "# ",
      openTag: "<h1>",
      closeTag: "</h1>",
    },
    Mapping {
      mdToken: "## ",
      openTag: "<h2>",
      closeTag: "</h2>",
    },
    Mapping {
      mdToken: "### ",
      openTag: "<h3>",
      closeTag: "</h3>",
    },
    Mapping {
      mdToken: "#### ",
      openTag: "<h4>",
      closeTag: "</h4>",
    },
    Mapping {
      mdToken: "##### ",
      openTag: "<h5>",
      closeTag: "</h5>",
    },
    Mapping {
      mdToken: "> ",
      openTag: "<blockquote>",
      closeTag: "</blockquote>",
    },
    Mapping {
      mdToken: "---",
      openTag: "<hr>",
      closeTag: "",
    },
  }

  newLine = line

  // simple prefix
  for _, mapping := range mappings {
    if strings.HasPrefix(newLine, mapping.mdToken) {
      newLine = ReplacePrefix(newLine, mapping.mdToken, mapping.openTag)
      newLine = AppendCloseTag(newLine, mapping.closeTag)
      break
    }
  }

  converted = (line == newLine)

  return newLine, converted

}

func handlePrefix(line string) (newLine string, converted bool) {

  type Mapping struct {
    mdToken   string
    openTag   string
    closeTag  string
  }

  var mappings = []Mapping {
    Mapping {
      mdToken: "# ",
      openTag: "<h1>",
      closeTag: "</h1>",
    },
    Mapping {
      mdToken: "## ",
      openTag: "<h2>",
      closeTag: "</h2>",
    },
    Mapping {
      mdToken: "### ",
      openTag: "<h3>",
      closeTag: "</h3>",
    },
    Mapping {
      mdToken: "#### ",
      openTag: "<h4>",
      closeTag: "</h4>",
    },
    Mapping {
      mdToken: "##### ",
      openTag: "<h5>",
      closeTag: "</h5>",
    },
    Mapping {
      mdToken: "> ",
      openTag: "<blockquote>",
      closeTag: "</blockquote>",
    },
    Mapping {
      mdToken: "---",
      openTag: "<hr>",
      closeTag: "",
    },
  }

  newLine = line

  // simple prefix
  for _, mapping := range mappings {
    if strings.HasPrefix(newLine, mapping.mdToken) {
      newLine = ReplacePrefix(newLine, mapping.mdToken, mapping.openTag)
      newLine = AppendCloseTag(newLine, mapping.closeTag)
      break
    }
  }

  converted = (line == newLine)

  return newLine, converted

}

func ReplacePrefix(line string, prefix string, newPrefix string) (string) {

  // TODO: Check prefix really exists
  if prefix == "" {
    return newPrefix + line
  }

  startPos := len(prefix)
  endPos   := len(line)
  return newPrefix + string([]byte(line)[startPos:endPos])
}

func AppendCloseTag(line string, closeTag string) (newline string) {
  return line + closeTag
}

func parseCliArgs() (config Config) {

  flag.StringVar(&config.fileName, "fileName", "", "File name to parse")
  flag.StringVar(&config.fileName, "f", "", "File name to parse")
  flag.Parse()
  args := flag.Args()
  if len(args) == 0 {
    fmt.Printf("No file specified.")
  } else {
    config.fileName = args[0]
  }

  return config
}

func regex(line string) (detected bool, newLine string) {

  // TODO: _foo_ -> <em> / `foo` -> <code> / **foo** -> <strong> / ~~foo~~ -> <del></del>
  var re *regexp.Regexp

  // ![img_title_alt](img_url)
  re = regexp.MustCompile(`(.*)!\[(.+)\]\((.+)\)(.*)`)
  newLine = re.ReplaceAllString(line, "${1}<figure><a href=\"$3\"><img src=\"$3\" alt=\"$2\" title=\"$2\" width=\"100%\"></a><figcaption>$2</figcaption></figure>")
  if newLine != line {
    return true, newLine
  }

  // [replace](url)
  re = regexp.MustCompile(`(.*)\[(.+)\]\((.+)\)(.*)`)
  newLine = re.ReplaceAllString(line, "${1}<a target=\"_blank\" ref=\"noopener noreferrer\" href=\"${3}\">${2}</a>${4}")
  if newLine != line {
    return true, newLine
  }

  // `foo`
  re = regexp.MustCompile("(.*)`(.*)`(.*)")
  newLine = re.ReplaceAllString(line, "$1<code>$2</code>$3")
  if newLine != line {
    return true, newLine
  }

  // __foo__
  re = regexp.MustCompile("(.*)__(.*)__(.*)")
  newLine = re.ReplaceAllString(line, "$1<em>$2</em>$3")
  if newLine != line {
    return true, newLine
  }

  // **foo**
  re = regexp.MustCompile("(.*)\\*\\*(.*)\\*\\*(.*)")
  newLine = re.ReplaceAllString(line, "$1<strong>$2</strong>$3")
  if newLine != line {
    return true, newLine
  }

  // ~~foo~~
  re = regexp.MustCompile("(.*)~~(.*)~~(.*)")
  newLine = re.ReplaceAllString(line, "$1<del>$2</del>$3")
  if newLine != line {
    return true, newLine
  }

  newLine = line
  return false, line
}


func handleOrderedList(line string, scanner) {


}













































