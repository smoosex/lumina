package notes

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublic(router *gin.RouterGroup) {
	router.GET("/reading-notes", h.list)
	router.GET("/reading-notes/:slug", h.detail)
	router.GET("/reading-tags", h.tags)
}

func (h *Handler) RegisterAdmin(router *gin.RouterGroup) {
	router.PUT("/reading-notes/:slug", h.upload)
}

func (h *Handler) list(c *gin.Context) {
	items, err := h.service.ListPublished()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load notes"})
		return
	}

	response := make([]noteListItem, 0, len(items))
	for _, item := range items {
		response = append(response, toListItem(item))
	}

	c.JSON(http.StatusOK, gin.H{"items": response})
}

func (h *Handler) detail(c *gin.Context) {
	note, err := h.service.GetPublishedBySlug(c.Param("slug"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load note"})
		return
	}

	c.JSON(http.StatusOK, toDetail(note))
}

func (h *Handler) tags(c *gin.Context) {
	tags, err := h.service.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tags"})
		return
	}

	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	c.JSON(http.StatusOK, gin.H{"items": names})
}

func (h *Handler) upload(c *gin.Context) {
	body, err := readUploadBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.Upload(c.Param("slug"), body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toDetail(note))
}

func readUploadBody(c *gin.Context) ([]byte, error) {
	contentType := c.GetHeader("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil && strings.EqualFold(mediaType, "multipart/form-data") {
		file, err := c.FormFile("file")
		if err != nil {
			return nil, errors.New("markdown file field is required")
		}

		opened, err := file.Open()
		if err != nil {
			return nil, errors.New("failed to open markdown file")
		}
		defer opened.Close()

		body, err := io.ReadAll(opened)
		if err != nil {
			return nil, errors.New("failed to read markdown file")
		}
		return body, nil
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.New("failed to read markdown")
	}
	return body, nil
}

type noteListItem struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	ISBN        string   `json:"isbn"`
	DoubanID    string   `json:"doubanId,omitempty"`
	BookTitle   string   `json:"bookTitle"`
	Author      string   `json:"author,omitempty"`
	Translator  string   `json:"translator,omitempty"`
	Publisher   string   `json:"publisher,omitempty"`
	Producer    string   `json:"producer,omitempty"`
	BookPubDate string   `json:"bookPubDate,omitempty"`
	BookPages   *int     `json:"bookPages,omitempty"`
	CoverURL    string   `json:"coverUrl,omitempty"`
	DoubanURL   string   `json:"doubanUrl,omitempty"`
	Excerpt     string   `json:"excerpt,omitempty"`
	Tags        []string `json:"tags"`
	Rating      *int     `json:"rating,omitempty"`
	ReadAt      *string  `json:"readAt,omitempty"`
	PublishedAt *string  `json:"publishedAt,omitempty"`
}

type noteDetail struct {
	noteListItem
	ContentHTML string `json:"contentHtml"`
	ContentMD   string `json:"contentMd,omitempty"`
}

func toListItem(note ReadingNote) noteListItem {
	item := noteListItem{
		Slug:        note.Slug,
		Title:       note.Title,
		ISBN:        note.ISBN,
		DoubanID:    note.DoubanID,
		BookTitle:   note.BookTitle,
		Author:      note.Author,
		Translator:  note.Translator,
		Publisher:   note.Publisher,
		Producer:    note.Producer,
		BookPubDate: note.BookPubDate,
		BookPages:   note.BookPages,
		CoverURL:    note.CoverURL,
		DoubanURL:   note.DoubanURL,
		Excerpt:     note.Excerpt,
		Tags:        tagNames(note.Tags),
		Rating:      note.Rating,
	}
	if note.ReadAt != nil {
		value := note.ReadAt.Format("2006-01-02")
		item.ReadAt = &value
	}
	if note.PublishedAt != nil {
		value := note.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		item.PublishedAt = &value
	}
	return item
}

func toDetail(note ReadingNote) noteDetail {
	return noteDetail{
		noteListItem: toListItem(note),
		ContentHTML:  note.ContentHTML,
		ContentMD:    note.ContentMD,
	}
}

func tagNames(tags []ReadingTag) []string {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	return names
}
