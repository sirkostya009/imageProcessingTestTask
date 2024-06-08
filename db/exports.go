package db

import "embed"

//go:embed schema.sql
var Schema string

//go:embed migrations
var Migrations embed.FS
