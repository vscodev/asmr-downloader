package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/vscodev/asmr-downloader/api"
)

var (
	_rjCodeRegex = regexp.MustCompile(`^(?i)(RJ)?(\d{6,})$`)
)

func parseProductID(rjCode string) (string, bool) {
	matched := _rjCodeRegex.FindStringSubmatch(rjCode)
	if len(matched) < 3 {
		return "", false
	}

	return matched[2], true
}

func main() {
	c := api.NewClient()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nInput the RJ code and press Enter to continue: ")

		rjCode, _ := reader.ReadString('\n')
		id, ok := parseProductID(strings.TrimSpace(rjCode))
		if !ok {
			continue
		}

		fmt.Println()
		log.Printf("Getting tracks of work (RJ%s)", id)
		tracks, err := c.GetTracks(id)
		if err != nil {
			log.Printf("Can not get tracks of work (RJ%s) : %s", id, err.Error())
			continue
		}

		log.Printf("Downloading tracks of work (RJ%s)", id)
		if err = c.DownloadTracks(id, tracks); err == nil {
			log.Printf("Tracks of work (RJ%s) downloaded", id)
		}
	}
}
