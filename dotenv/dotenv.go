package dotenv

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ParseDotenv(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "=")
		key, value := split[0], split[1]
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		os.Setenv(key, value)
	}
}
