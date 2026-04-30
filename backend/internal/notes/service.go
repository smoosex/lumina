package notes

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/smoosex/lumina/backend/internal/bookmeta"
	"github.com/smoosex/lumina/backend/internal/markdown"
)

type Service struct {
	repo     *Repository
	markdown *markdown.Service
	books    bookmeta.Client
}

func NewService(
	repo *Repository,
	markdown *markdown.Service,
	books bookmeta.Client,
) *Service {
	return &Service{repo: repo, markdown: markdown, books: books}
}

func (s *Service) ListPublished() ([]ReadingNote, error) {
	return s.repo.ListPublished()
}

func (s *Service) GetPublishedBySlug(slug string) (ReadingNote, error) {
	return s.repo.GetPublishedBySlug(slug)
}

func (s *Service) ListTags() ([]ReadingTag, error) {
	return s.repo.ListTags()
}

func (s *Service) Upload(slug string, body []byte) (ReadingNote, error) {
	doc, err := s.markdown.Parse(body)
	if err != nil {
		return ReadingNote{}, err
	}
	if doc.Meta.Slug != slug {
		return ReadingNote{}, errors.New("path slug does not match frontmatter slug")
	}

	book, err := s.books.FetchByISBN(context.Background(), doc.Meta.ISBN)
	if err != nil {
		return ReadingNote{}, err
	}

	coverURL := ""
	if book.CoverURL != "" {
		coverURL, err = s.saveCover(context.Background(), doc.Meta.Slug, book)
		if err != nil {
			return ReadingNote{}, err
		}
	}

	note := ReadingNote{
		Slug:        doc.Meta.Slug,
		Title:       doc.Meta.Title,
		ISBN:        book.ISBN,
		DoubanID:    book.DoubanID,
		BookTitle:   book.Title,
		Author:      book.Author,
		Translator:  book.Translator,
		Publisher:   book.Publisher,
		Producer:    book.Producer,
		BookPubDate: book.PublishedAt,
		BookPages:   book.Pages,
		CoverURL:    coverURL,
		DoubanURL:   book.URL,
		Excerpt:     doc.Meta.Excerpt,
		ContentMD:   doc.ContentMD,
		ContentHTML: doc.ContentHTML,
		Status:      doc.Meta.Status,
		Rating:      doc.Meta.Rating,
		PublishedAt: doc.Meta.PublishedAt,
	}
	if doc.Meta.ReadAt != nil {
		note.ReadAt = &doc.Meta.ReadAt.Time
	}

	return s.repo.Upsert(note, doc.Meta.Tags)
}

func (s *Service) saveCover(
	ctx context.Context,
	slug string,
	book bookmeta.Metadata,
) (string, error) {
	cover, err := s.books.DownloadCover(ctx, book.CoverURL, book.URL)
	if err != nil {
		return "", err
	}

	extension := ".jpg"
	if strings.Contains(cover.ContentType, "png") {
		extension = ".png"
	} else if strings.Contains(cover.ContentType, "webp") {
		extension = ".webp"
	}

	if err := os.MkdirAll(filepath.Join("data", "covers"), 0o755); err != nil {
		return "", err
	}

	fileName := slug + extension
	path := filepath.Join("data", "covers", fileName)
	if err := os.WriteFile(path, cover.Data, 0o644); err != nil {
		return "", err
	}

	return "/covers/" + fileName, nil
}
