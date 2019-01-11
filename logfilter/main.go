package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var threshold = flag.Int("tl", 3000, "Time limit above which the queries have to be monitored")

func handleLine(line string) {
	if !strings.Contains(line, "COMMAND") {
		return
	}
	lexems := strings.Split(line, " ")
	// timestamp := lexems[0]
	tts := lexems[len(lexems)-1]
	tts = strings.Trim(tts, "\n")
	tts = tts[:len(tts)-2]
	// tt, err := strconv.Atoi(tts)
	// if err != nil {
	// 	continue
	// }
	// if tt > *threshold {
	// 	fmt.Println(timestamp, tts)
	// 	fmt.Println(line.Text)
	// 	break
	// }
	fmt.Println(getStringStartingWith(lexems, "protocol"))
}

func main() {
	flag.Parse()
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
