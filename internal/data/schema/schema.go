package schema

import "embed"

// Migrations isi dari folder migrations yang di-embed ke dalam binary.
//
//go:embed migrations/*.sql seeds/*.sql
var Migrations embed.FS
