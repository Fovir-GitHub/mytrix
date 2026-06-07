package assets

import "embed"

//go:embed sql/schema
var SchemaFS embed.FS

//go:embed sql/migrations
var MigrationsFS embed.FS

const SchemaPath = "sql/schema"

const MigrationsPath = "sql/migrations"
