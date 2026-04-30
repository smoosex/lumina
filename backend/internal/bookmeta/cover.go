package bookmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CoverImage struct {
	ContentType string
	Data        []byte
}

func (c *DoubanClient) DownloadCover(
	ctx context.Context,
	coverURL string,
	referer string,
) (CoverImage, error) {
	resp, err := c.fetchCover(ctx, coverURL, referer, "")
	if err != nil {
		return CoverImage{}, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		challengeBody, _ := io.ReadAll(io.LimitReader(resp.Body, 32*1024))
		cookie := doubanChallengeCookie(string(challengeBody))
		if cookie == "" {
			return CoverImage{}, fmt.Errorf("cover source did not return an image")
		}

		resp.Body.Close()
		resp, err = c.fetchCover(ctx, coverURL, referer, cookie)
		if err != nil {
			return CoverImage{}, err
		}
		defer resp.Body.Close()
		contentType = resp.Header.Get("Content-Type")
	}

	if !strings.HasPrefix(contentType, "image/") {
		return CoverImage{}, fmt.Errorf("cover source did not return an image")
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return CoverImage{}, err
	}

	return CoverImage{ContentType: contentType, Data: data}, nil
}

func (c *DoubanClient) fetchCover(
	ctx context.Context,
	coverURL string,
	referer string,
	cookie string,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, coverURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Referer", firstNonEmpty(referer, "https://book.douban.com/"))
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	return c.httpClient.Do(req)
}

func doubanChallengeCookie(body string) string {
	matches := regexp.MustCompile(`[A-Za-z_$][\w$]*:(\d+)`).FindAllStringSubmatch(body, -1)
	if len(matches) < 3 {
		return ""
	}

	status := 0
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err == nil {
			status += value
		}
	}

	sessionMatch := regexp.MustCompile(`\(t,(\d+)\);continue`).FindStringSubmatch(body)
	if len(sessionMatch) != 2 {
		return ""
	}

	return "__tst_status=" + strconv.Itoa(status) +
		"; EO_Bot_Ssid=" + sessionMatch[1]
}
