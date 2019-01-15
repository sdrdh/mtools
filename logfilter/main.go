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
	Fast         int      `long:"fast" description:"Returns lines which are faster than the value in ms"`
	NonIndex     bool     `long:"ni" description:"Queries that aren't indexed"`
	KeysExamined int      `long:"keysExamined" description:"Returns Queries which have examined more keys than this"`
	Keywords     []string `short:"w" long:"keywords" description:"Keywords to track"`
	MultiIndex   bool     `long:"mi" description:"Queries using multiple indexes"`
	DB           []string `long:"database" description:"Databases to track"`
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
	if checkKeyword(line) &&
		timeQuery(tt) &&
		(nonIndexed(line) || multiIndexed(line)) &&
		keysExamined(ke) &&
		toTrackDatabase(getDatabase(lexems)) {
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

func timeQuery(t int) bool {
	if opts.Slow == 0 && opts.Fast == 0 {
		return true
	}
	if opts.Slow == 0 {
		return opts.Fast > t
	}
	if opts.Fast == 0 {
		return t > opts.Slow
	}
	return t > opts.Slow && opts.Fast > t
}

func getDatabase(lexems []string) string {
	for i, v := range lexems {
		if v == "$db:" {
			return strings.Replace(lexems[i+1], `"`, "", -1)
		}
	}
	return ""
}

func toTrackDatabase(db string) bool {
	if len(opts.DB) == 0 {
		return true
	}
	for _, v := range opts.DB {
		if v == db {
			return true
		}
	}
	return false
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
