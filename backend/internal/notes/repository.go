package notes

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPublished() ([]ReadingNote, error) {
	var notes []ReadingNote
	err := r.db.
		Preload("Tags").
		Where("status = ?", "published").
		Order("published_at desc nulls last").
		Find(&notes).
		Error
	return notes, err
}

func (r *Repository) GetPublishedBySlug(slug string) (ReadingNote, error) {
	var note ReadingNote
	err := r.db.
		Preload("Tags").
		Where("slug = ? and status = ?", slug, "published").
		First(&note).
		Error
	return note, err
}

func (r *Repository) ListTags() ([]ReadingTag, error) {
	var tags []ReadingTag
	err := r.db.
		Joins("join reading_note_tags on reading_note_tags.tag_id = reading_tags.id").
		Joins("join reading_notes on reading_notes.id = reading_note_tags.note_id").
		Where("reading_notes.status = ?", "published").
		Group("reading_tags.id").
		Order("reading_tags.name asc").
		Find(&tags).
		Error
	return tags, err
}

func (r *Repository) Upsert(note ReadingNote, tagNames []string) (ReadingNote, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		tags, err := r.upsertTags(tx, tagNames)
		if err != nil {
			return err
		}

		var existing ReadingNote
		err = tx.Where("slug = ?", note.Slug).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(&note).Error; err != nil {
				return err
			}
			existing = note
		} else {
			note.ID = existing.ID
			if err := tx.Model(&existing).Updates(map[string]any{
				"title":         note.Title,
				"isbn":          note.ISBN,
				"douban_id":     note.DoubanID,
				"book_title":    note.BookTitle,
				"author":        note.Author,
				"translator":    note.Translator,
				"publisher":     note.Publisher,
				"producer":      note.Producer,
				"book_pub_date": note.BookPubDate,
				"book_pages":    note.BookPages,
				"cover_url":     note.CoverURL,
				"douban_url":    note.DoubanURL,
				"excerpt":       note.Excerpt,
				"content_md":    note.ContentMD,
				"content_html":  note.ContentHTML,
				"status":        note.Status,
				"rating":        note.Rating,
				"read_at":       note.ReadAt,
				"published_at":  note.PublishedAt,
			}).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&existing).Association("Tags").Replace(tags); err != nil {
			return err
		}

		return tx.Preload("Tags").First(&note, "slug = ?", note.Slug).Error
	})

	return note, err
}

func (r *Repository) upsertTags(tx *gorm.DB, tagNames []string) ([]ReadingTag, error) {
	tags := make([]ReadingTag, 0, len(tagNames))
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		tag := ReadingTag{Name: name}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&tag).Error; err != nil {
			return nil, err
		}
		if tag.ID == "" {
			if err := tx.Where("name = ?", name).First(&tag).Error; err != nil {
				return nil, err
			}
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
