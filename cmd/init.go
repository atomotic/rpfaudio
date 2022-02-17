package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/atomotic/rpfaudio/utils"
	"github.com/facette/natsort"
	"github.com/gofrs/uuid"
	"github.com/readium/go-toolkit/pkg/manifest"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a manifest",
	Long:  `Write an audiobook manifest`,
	Run: func(cmd *cobra.Command, args []string) {
		m := manifest.Manifest{}
		m.Metadata.Type = "https://schema.org/Audiobook"
		m.Metadata.ConformsTo = append(m.Metadata.ConformsTo, manifest.ProfileAudiobook)
		u1 := uuid.Must(uuid.NewV4())
		m.Metadata.Identifier = fmt.Sprintf("urn:uuid:%s", u1)
		m.Metadata.Languages = append(m.Metadata.Languages, " ")

		emptyContributors := manifest.Contributor{LocalizedName: manifest.LocalizedString{}}
		m.Metadata.Authors = append(m.Metadata.Authors, emptyContributors)
		m.Metadata.Narrators = append(m.Metadata.Narrators, emptyContributors)
		m.Metadata.Publishers = append(m.Metadata.Publishers, emptyContributors)

		t := time.Now()
		m.Metadata.Published = &t
		m.Metadata.Modified = &t

		var linkList []manifest.Link
		var totalDuration float64
		tracks, err := filepath.Glob("*.mp3")
		if err != nil {
			log.Fatal(err)
		}

		if len(tracks) == 0 {
			fmt.Println("no mp3 tracks")
			os.Exit(1)
		}

		natsort.Sort(tracks)
		fmt.Println("# processing tracks")
		bar := progressbar.Default(int64(len(tracks)))
		for i, track := range tracks {

			duration, _ := utils.MP3Duration(track)
			mime, _ := utils.Mime(track)

			linkList = append(linkList, manifest.Link{
				Title:    fmt.Sprintf("Track %d", i+1),
				Href:     track,
				Type:     mime,
				Duration: duration})
			totalDuration = totalDuration + duration
			bar.Add(1)
		}

		m.ReadingOrder = linkList
		m.Metadata.Duration = &totalDuration

		var resources []manifest.Link

		if _, err := os.Stat("cover.jpg"); err == nil {
			resources = append(resources, manifest.Link{
				Title: "Cover",
				Href:  "cover.jpg",
				Type:  "image/jpg",
				Rels:  []string{"cover"},
			})

		}
		m.Resources = resources
		jsonmanifest, _ := m.MarshalJSON()
		_ = ioutil.WriteFile("manifest.json", jsonmanifest, 0644)
		fmt.Println("# edit manifest.json, then run `rpfaudio package`")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
