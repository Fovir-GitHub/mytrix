package assets

import "embed"

//go:embed sql/schema
var SchemaFS embed.FS

const SchemaPath = "sql/schema"
