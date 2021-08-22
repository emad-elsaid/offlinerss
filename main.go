package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	RootDirName    = "rss"
	InboxDirName   = "INBOX"
	MetaDirName    = ".meta"
	ConfigFileName = "config.json"
)

type Config struct {
	URLs []string `json:"urls"`
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Can't user home directory\n\t%s", err)
	}

	destination, err := ensurePath(home, RootDirName)
	if err != nil {
		log.Fatal(err)
	}

	inbox, err := ensurePath(destination, InboxDirName)
	if err != nil {
		log.Fatal(err)
	}

	metaDir, err := ensurePath(destination, MetaDirName)
	if err != nil {
		log.Fatal(err)
	}

	configFile := filepath.Join(destination, ConfigFileName)
	configContent, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Can't read %s\n\t%s", configFile, err)
	}

	var config Config
	if err = json.Unmarshal(configContent, &config); err != nil {
		log.Fatalf("Can't parse %s as JSON\n\t%s", configFile, err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(config.URLs))

	for _, url := range config.URLs {
		go func(url string) {
			err := refreshRSS(config, destination, metaDir, inbox, url)
			if err != nil {
				log.Printf("Error refreshing RSS: %s\n\t%s", url, err)
			}

			wg.Done()
		}(url)
	}

	wg.Wait()
}

func ensurePath(paths ...string) (string, error) {
	path := filepath.Join(paths...)
	err := os.MkdirAll(path, fs.ModeDir|fs.ModePerm)
	if err != nil {
		err = fmt.Errorf("Can't ensure the existence of %s directory\n\t%w", path, err)
	}

	return path, err
}

func refreshRSS(config Config, destination, metaDir, inbox, url string) error {
	log.Printf("[GET] %s", url)

	urlDigest := fmt.Sprintf("%x", sha1.Sum([]byte(url)))
	meta, items, err := fetchRSS(url, config)
	if err != nil {
		return fmt.Errorf("Error fetching: %s", err)
	}

	// Write meta file
	itemMeta := filepath.Join(metaDir, urlDigest+".rss")
	err = os.WriteFile(itemMeta, meta, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error writing meta file: %s", err)
	}

	// Write items files if they don't already exist in any subdirectory
	for _, item := range items {
		err = writeItem(destination, inbox, urlDigest, item)
		if err != nil {
			return fmt.Errorf("Error writing item file: %s", err)
		}
	}

	return nil
}

var itemsTags = []struct {
	start []byte
	end   []byte
}{
	{start: []byte("<entry>"), end: []byte("</entry>")},
	{start: []byte("<entry "), end: []byte("</entry>")}, // some RSS have attributes for each entry
	{start: []byte("<item>"), end: []byte("</item>")},
	{start: []byte("<item "), end: []byte("</item>")}, // same for item tags
}

func fetchRSS(url string, config Config) ([]byte, [][]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("GET: %s\n\tError: Status %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Extract Meta
	meta := bytes.NewBuffer([]byte{})
	for _, tag := range itemsTags {
		start := bytes.Index(body, tag.start)
		if start <= 0 {
			continue
		}

		end := bytes.LastIndex(body, tag.end)
		meta.Write(body[:start-1])
		meta.Write(body[end+len(tag.end):])
	}

	// Extract items/entries
	items := [][]byte{}
	for _, tag := range itemsTags {
		tagBody := body
		for {
			start := bytes.Index(tagBody, tag.start)
			if start <= 0 {
				break
			}

			end := bytes.Index(tagBody, tag.end)
			if end <= start {
				break
			}

			item := tagBody[start : end+len(tag.end)]
			items = append(items, item)

			tagBody = tagBody[end+len(tag.end):]
		}
	}

	return meta.Bytes(), items, nil
}

func writeItem(destination, inbox, urlDigest string, item []byte) error {
	itemDigest := fmt.Sprintf("%x", sha1.Sum(item))
	itemFile := urlDigest + "-" + itemDigest + ".rss"
	itemPath := filepath.Join(inbox, itemFile)
	itemGlobPath := filepath.Join(destination, "**", itemFile)

	if matches, err := filepath.Glob(itemGlobPath); len(matches) > 0 || err != nil {
		return nil
	}

	return os.WriteFile(itemPath, item, os.ModePerm)
}
