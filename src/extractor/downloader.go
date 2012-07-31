package extractor

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type downloader struct {
    root   string
    client *http.Client
}

func newDownloader(root string, timeout time.Duration, deadline time.Time) downloader {
    return downloader{
        root:   root,
        client: newTimeoutDeadlineDialer(timeout, deadline),
    }
}

func (d *downloader) output(path string) string {
    return fmt.Sprintf("%s/%s", d.root, path)
}

func (d *downloader) downloadToFile(url, path string) error {
    resp, err := d.client.Get(url)
    if err != nil {
        return fmt.Errorf("downloader: HTTP request failed: %s", err)
    }
    defer resp.Body.Close()
    file, err := os.OpenFile(d.output(path), os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("downloader: File open failed: %s", err)
    }
    defer file.Close()

    written, err := io.Copy(file, resp.Body)
    if err != nil {
        return fmt.Errorf("downloader: Failed copying to file; %s", err)
    }

    if resp.ContentLength > 0 && written != resp.ContentLength {
        return fmt.Errorf("downloader: written != expected: %d != %d", written, resp.ContentLength)
    }

    return nil
}