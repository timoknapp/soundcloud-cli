package sccli

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
)

// PrintTable prints a fomatted table of tracks
func PrintTable(tracks []soundcloud.Track) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "URL"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, track := range tracks {
		table.Append([]string{strconv.Itoa(track.ID), track.Title, track.PermalinkURL})
	}
	table.Render()
}
