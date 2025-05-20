package sccli

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/timoknapp/soundcloud-cli/pkg/soundcloud"
)

// PrintTable prints a fomatted table of tracks
func PrintTable(tracks []soundcloud.Track) {
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader([]string{"ID", "Title", "URL"}),
	)
	
	for _, track := range tracks {
		table.Append([]string{strconv.Itoa(track.ID), track.Title, track.PermalinkURL})
	}
	table.Render()
}
