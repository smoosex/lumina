package notes

import "time"

type ReadingNote struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Slug        string `gorm:"uniqueIndex;not null"`
	Title       string `gorm:"not null"`
	ISBN        string `gorm:"not null"`
	DoubanID    string
	BookTitle   string `gorm:"not null"`
	Author      string
	Translator  string
	Publisher   string
	Producer    string
	BookPubDate string
	BookPages   *int
	CoverURL    string
	DoubanURL   string
	Excerpt     string
	ContentMD   string `gorm:"column:content_md;not null"`
	ContentHTML string `gorm:"column:content_html"`
	Status      string `gorm:"not null;default:draft"`
	Rating      *int
	ReadAt      *time.Time
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tags        []ReadingTag `gorm:"many2many:reading_note_tags;joinForeignKey:NoteID;joinReferences:TagID"`
}

type ReadingTag struct {
	ID   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

func (ReadingNote) TableName() string {
	return "reading_notes"
}

func (ReadingTag) TableName() string {
	return "reading_tags"
}
