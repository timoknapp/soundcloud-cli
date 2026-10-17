package sccli

import (
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
)

var headerColor = color.New(color.Bold, color.FgCyan).SprintFunc()

// PrintTable prints a fomatted table of tracks
func PrintTable(tracks []soundcloud.Track) {
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRendition(tw.Rendition{Borders: tw.BorderNone}),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithHeaderFilter(tw.CellFilter{Global: func(cells []string) []string {
			for i, c := range cells {
				cells[i] = headerColor(c)
			}
			return cells
		}}),
	)
	table.Header("ID", "Title", "URL")
	for _, track := range tracks {
		table.Append(strconv.Itoa(track.ID), track.Title, track.PermalinkURL)
	}
	table.Render()
}
