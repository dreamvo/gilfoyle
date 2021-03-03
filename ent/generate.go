//go:generate go run github.com/facebook/ent/cmd/ent generate ./schema
package ent

import (
	_ "github.com/facebook/ent/entc/gen"
	_ "github.com/mattn/go-runewidth"
	_ "github.com/olekukonko/tablewriter"
	_ "github.com/russross/blackfriday/v2"
	_ "github.com/shurcooL/sanitized_anchor_name"
)
