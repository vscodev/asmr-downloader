package api

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/vscodev/asmr-downloader/model"
)

type Client struct {
	inner *http.Client
}

func NewClient() *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Jar: jar,
		//Timeout: 30 * time.Second,
	}

	return &Client{inner: httpClient}
}

func (c *Client) GetTracks(id string) ([]*model.Track, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.asmr.one/api/tracks/"+id, nil)
	resp, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrTracksNotFound
	}

	var tracks []*model.Track
	if err = json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (c *Client) downloadFile(name string, url string) error {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (c *Client) downloadTrack(parent string, track *model.Track) error {
	currentPath := filepath.Join(parent, track.Title)
	if track.IsFolder() {
		_ = os.MkdirAll(currentPath, 0755)
		for _, child := range track.Children {
			if err := c.downloadTrack(currentPath, child); err != nil {
				log.Printf("Can not download track (%s) : %s", filepath.Join(currentPath, child.Title), err.Error())
			}
		}
		return nil
	} else {
		log.Printf("Downloading track (%s)", currentPath)
		return c.downloadFile(currentPath, track.MediaDownloadURL)
	}
}

func (c *Client) DownloadTracks(id string, tracks []*model.Track) error {
	basePath := "RJ" + id
	for _, track := range tracks {
		if err := c.downloadTrack(basePath, track); err != nil {
			log.Printf("Can not download track (%s) : %s", filepath.Join(basePath, track.Title), err.Error())
		}
	}
	return nil
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Referer", "https://www.asmr.one")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	return c.inner.Do(req)
}
