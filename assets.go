package assets

import "embed"

//go:embed views/* static/* public/*
var EmbeddedFiles embed.FS
