package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/smoosex/lumina/backend/internal/config"
	"github.com/yuin/goldmark"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

var Module = fx.Module("markdown", fx.Provide(New))

type Service struct {
	cfg       config.MarkdownConfig
	renderer  goldmark.Markdown
	sanitizer *bluemonday.Policy
}

type Document struct {
	Meta        Frontmatter
	ContentMD   string
	ContentHTML string
}

type Frontmatter struct {
	Slug        string     `yaml:"slug"`
	Title       string     `yaml:"title"`
	ISBN        string     `yaml:"isbn"`
	Status      string     `yaml:"status"`
	Rating      *int       `yaml:"rating"`
	ReadAt      *Date      `yaml:"read_at"`
	PublishedAt *time.Time `yaml:"published_at"`
	Tags        []string   `yaml:"tags"`
	Excerpt     string     `yaml:"excerpt"`
}

type Date struct {
	time.Time
}

func New(cfg config.Config) *Service {
	return &Service{
		cfg:       cfg.Markdown,
		renderer:  goldmark.New(),
		sanitizer: bluemonday.UGCPolicy(),
	}
}

func (s *Service) Parse(raw []byte) (Document, error) {
	metaBytes, bodyBytes, err := splitFrontmatter(raw)
	if err != nil {
		return Document{}, err
	}

	var meta Frontmatter
	if err := yaml.Unmarshal(metaBytes, &meta); err != nil {
		return Document{}, fmt.Errorf("parse frontmatter: %w", err)
	}
	if err := meta.Validate(); err != nil {
		return Document{}, err
	}

	contentMD := strings.TrimSpace(string(bodyBytes))
	if contentMD == "" {
		return Document{}, errors.New("markdown body is required")
	}

	doc := Document{
		Meta:      meta,
		ContentMD: contentMD,
	}

	if s.cfg.RenderHTML {
		var out bytes.Buffer
		if err := s.renderer.Convert([]byte(contentMD), &out); err != nil {
			return Document{}, fmt.Errorf("render markdown: %w", err)
		}
		html := out.String()
		if s.cfg.SanitizeHTML {
			html = s.sanitizer.Sanitize(html)
		}
		doc.ContentHTML = html
	}

	return doc, nil
}

func (m Frontmatter) Validate() error {
	if strings.TrimSpace(m.Slug) == "" {
		return errors.New("slug is required")
	}
	if strings.TrimSpace(m.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(m.ISBN) == "" {
		return errors.New("isbn is required")
	}
	if strings.TrimSpace(m.Status) == "" {
		return errors.New("status is required")
	}
	if m.Status != "draft" && m.Status != "published" {
		return errors.New("status must be draft or published")
	}
	return nil
}

func splitFrontmatter(raw []byte) ([]byte, []byte, error) {
	text := strings.TrimLeft(string(raw), "\ufeff\r\n\t ")
	if !strings.HasPrefix(text, "---\n") && !strings.HasPrefix(text, "---\r\n") {
		return nil, nil, errors.New("frontmatter block is required")
	}

	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	parts := strings.SplitN(normalized, "\n---\n", 2)
	if len(parts) != 2 {
		return nil, nil, errors.New("frontmatter block is not closed")
	}

	meta := strings.TrimPrefix(parts[0], "---\n")
	return []byte(meta), []byte(parts[1]), nil
}

func (d *Date) UnmarshalYAML(value *yaml.Node) error {
	parsed, err := time.Parse("2006-01-02", value.Value)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}
