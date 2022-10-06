package testapp

import "embed"

//go:embed webui/build/*
var EmbeddedFiles embed.FS
