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
	_apiClient   = api.NewClient()
	_rjCodeRegex = regexp.MustCompile(`^(?i)(RJ)?(\d{6,})$`)
)

func parseProductID(rjCode string) (string, bool) {
	matched := _rjCodeRegex.FindStringSubmatch(rjCode)
	if len(matched) < 3 {
		return "", false
	}

	return matched[2], true
}

func downloadWork(id string) {
	fmt.Println()

	log.Printf("Getting tracks of work (RJ%s)", id)
	tracks, err := _apiClient.GetTracks(id)
	if err != nil {
		log.Printf("Can not get tracks of work (RJ%s) : %s", id, err.Error())
		return
	}

	log.Printf("Downloading tracks of work (RJ%s)", id)
	if err = _apiClient.DownloadTracks(id, tracks); err == nil {
		log.Printf("Tracks of work (RJ%s) downloaded", id)
	}
}

func downloadWorks(rjCodes []string) {
	for _, rjCode := range rjCodes {
		id, ok := parseProductID(rjCode)
		if !ok {
			continue
		}

		downloadWork(id)
	}
}

func main() {
	if len(os.Args) > 1 {
		downloadWorks(os.Args[1:])
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("\nInput the RJ codes separated by space and press Enter to continue: ")

			rjCodes, _ := reader.ReadString('\n')
			downloadWorks(strings.Fields(rjCodes))
		}
	}
}
