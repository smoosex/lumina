package bookmeta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	ISBN        string
	DoubanID    string
	Title       string
	Author      string
	Translator  string
	Publisher   string
	Producer    string
	PublishedAt string
	Pages       *int
	CoverURL    string
	URL         string
}

type Client interface {
	FetchByISBN(ctx context.Context, isbn string) (Metadata, error)
	DownloadCover(ctx context.Context, coverURL string, referer string) (CoverImage, error)
}

type DoubanClient struct {
	httpClient *http.Client
}

func NewDoubanClient() Client {
	return &DoubanClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *DoubanClient) FetchByISBN(ctx context.Context, isbn string) (Metadata, error) {
	isbn = normalizeISBN(isbn)
	if isbn == "" {
		return Metadata{}, errors.New("isbn is required")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://book.douban.com/isbn/%s/", isbn),
		nil,
	)
	if err != nil {
		return Metadata{}, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Metadata{}, fmt.Errorf("fetch douban book: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Metadata{}, fmt.Errorf("douban returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return Metadata{}, err
	}

	meta, err := parseDoubanHTML(string(body))
	if err != nil {
		return Metadata{}, err
	}
	meta.ISBN = isbn
	if meta.URL == "" {
		meta.URL = resp.Request.URL.String()
	}
	if meta.DoubanID == "" {
		meta.DoubanID = doubanIDFromURL(meta.URL)
	}
	return meta, nil
}

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124 Safari/537.36"

type jsonLDBook struct {
	Name   string `json:"name"`
	ISBN   string `json:"isbn"`
	URL    string `json:"url"`
	SameAs string `json:"sameAs"`
	Author []struct {
		Name string `json:"name"`
	} `json:"author"`
}

func parseDoubanHTML(raw string) (Metadata, error) {
	var meta Metadata

	if match := regexp.MustCompile(`(?s)<script type="application/ld\+json">\s*(.*?)\s*</script>`).FindStringSubmatch(raw); len(match) == 2 {
		var book jsonLDBook
		if err := json.Unmarshal([]byte(match[1]), &book); err == nil {
			meta.Title = strings.TrimSpace(book.Name)
			meta.ISBN = strings.TrimSpace(book.ISBN)
			meta.URL = firstNonEmpty(book.URL, book.SameAs)
			if len(book.Author) > 0 {
				authors := make([]string, 0, len(book.Author))
				for _, author := range book.Author {
					if name := strings.TrimSpace(author.Name); name != "" {
						authors = append(authors, name)
					}
				}
				meta.Author = strings.Join(authors, " / ")
			}
		}
	}

	meta.Title = firstNonEmpty(meta.Title, htmlMeta(raw, `og:title`))
	meta.CoverURL = htmlMeta(raw, `og:image`)
	meta.URL = firstNonEmpty(meta.URL, htmlMeta(raw, `og:url`))
	meta.DoubanID = doubanIDFromURL(meta.URL)

	info := bookInfo(raw)
	meta.Author = firstNonEmpty(meta.Author, info["作者"])
	meta.Translator = info["译者"]
	meta.Publisher = info["出版社"]
	meta.Producer = info["出品方"]
	meta.PublishedAt = info["出版年"]
	if pages := parsePages(info["页数"]); pages != nil {
		meta.Pages = pages
	}

	if meta.Title == "" {
		return Metadata{}, errors.New("douban book title not found")
	}
	return meta, nil
}

func htmlMeta(raw string, property string) string {
	pattern := fmt.Sprintf(`<meta property="%s" content="([^"]*)"`, regexp.QuoteMeta(property))
	match := regexp.MustCompile(pattern).FindStringSubmatch(raw)
	if len(match) != 2 {
		return ""
	}
	return strings.TrimSpace(html.UnescapeString(match[1]))
}

func bookInfo(raw string) map[string]string {
	values := map[string]string{}
	match := regexp.MustCompile(`(?s)<div id="info"[^>]*>(.*?)</div>`).FindStringSubmatch(raw)
	if len(match) != 2 {
		return values
	}

	text := regexp.MustCompile(`(?i)<br\s*/?>`).ReplaceAllString(match[1], "\n")
	text = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(text, "")
	text = html.UnescapeString(text)

	lines := make([]string, 0)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if !strings.Contains(line, ":") && !strings.Contains(line, "：") {
			continue
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			key, value, ok = strings.Cut(line, "：")
		}
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if value == "" && i+1 < len(lines) {
			value = strings.TrimSpace(lines[i+1])
		}
		if key != "" && value != "" {
			values[key] = value
		}
	}

	return values
}

func doubanIDFromURL(value string) string {
	match := regexp.MustCompile(`/subject/(\d+)/?`).FindStringSubmatch(value)
	if len(match) != 2 {
		return ""
	}
	return match[1]
}

func normalizeISBN(isbn string) string {
	return strings.NewReplacer("-", "", " ", "").Replace(strings.TrimSpace(isbn))
}

func parsePages(value string) *int {
	match := regexp.MustCompile(`\d+`).FindString(value)
	if match == "" {
		return nil
	}
	pages, err := strconv.Atoi(match)
	if err != nil {
		return nil
	}
	return &pages
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
