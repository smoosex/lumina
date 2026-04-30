package app

import (
	"github.com/smoosex/lumina/backend/internal/bookmeta"
	"github.com/smoosex/lumina/backend/internal/config"
	"github.com/smoosex/lumina/backend/internal/db"
	httpserver "github.com/smoosex/lumina/backend/internal/http"
	"github.com/smoosex/lumina/backend/internal/markdown"
	"github.com/smoosex/lumina/backend/internal/notes"
	"github.com/smoosex/lumina/backend/internal/uploadauth"
	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		config.Module,
		db.Module,
		bookmeta.Module,
		markdown.Module,
		uploadauth.Module,
		notes.Module,
		httpserver.Module,
	)
}
