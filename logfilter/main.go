package main

import (
	"bufio"

	// "flag"
	// // "flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

// var threshold = flag.Int("tl", 3000, "Time limit above which the queries have to be monitored")

var opts struct {
	Slow         int      `long:"slow" description:"Returns lines which are slower than the value in ms"`
	NonIndex     bool     `long:"ni" description:"Queries that aren't indexed"`
	KeysExamined int      `long:"keysExamined" description:"Returns Queries which have examined more keys than this"`
	Keywords     []string `short:"w" long:"keywords" description:"Keywords to track"`
	MultiIndex   bool     `long:"mi" description:"Queries using multiple indexes"`
}

func handleLine(line string) {
	if !strings.Contains(line, "COMMAND") {
		return
	}
	lexems := strings.Split(line, " ")
	tts := lexems[len(lexems)-1]
	tts = strings.Trim(tts, "\n")
	tts = tts[:len(tts)-2]
	tt, err := strconv.Atoi(tts)
	if err != nil {
		return
	}
	// if opts.Slow != 0 && tt > opts.Slow {
	// 	fmt.Print(line)
	// 	return
	// }
	ke, err := strconv.Atoi(getDataStartingWith(lexems, "keysExamined"))
	if err != nil {
		return
	}
	if checkKeyword(line) && slowQuery(tt) && (nonIndexed(line) || multiIndexed(line)) && keysExamined(ke) {
		fmt.Print(line)
	}
}

func checkKeyword(line string) bool {
	if len(opts.Keywords) == 0 {
		return true
	}
	for _, v := range opts.Keywords {
		if strings.Contains(line, v) {
			// fmt.Print(line)
			return true
		}
	}
	return false
}

func nonIndexed(line string) bool {
	if !opts.NonIndex {
		return true
	}
	return !strings.Contains(line, "IXSCAN")
}

func multiIndexed(line string) bool {
	if !opts.MultiIndex {
		return true
	}
	return strings.Count(line, "IXSCAN") > 1
}

func keysExamined(ke int) bool {
	if opts.KeysExamined == 0 {
		return true
	}
	return ke > opts.KeysExamined
}

func slowQuery(t int) bool {
	if opts.Slow == 0 {
		return true
	}
	return t > opts.Slow
}

func main() {
	// flag.Parse()
	_, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}
	if opts.NonIndex && opts.MultiIndex {
		panic(fmt.Errorf("ni and mi can't be set at the same time"))
	}
	// log.Printf("%+v\n", opts)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		handleLine(line)
		// break
	}
}

func getStringStartingWith(from []string, target string) string {
	for _, v := range from {
		if len(target) > len(v) {
			continue
		}
		if target == v[:len(target)] {
			return v
		}
	}
	return ""
}

func getDataStartingWith(from []string, target string) string {
	s := getStringStartingWith(from, target)
	if s == "" {
		return ""
	}
	target = target + ":"
	return strings.Replace(s, target, "", 1)
}
