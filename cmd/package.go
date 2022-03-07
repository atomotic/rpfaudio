package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/facette/natsort"
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the audiobook",
	Long:  `Package the audiobook`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("# create RPF audiobook")
		audiobook, err := os.Create("new.audiobook")
		if err != nil {
			log.Fatal(err)
		}
		defer audiobook.Close()

		zipwriter := zip.NewWriter(audiobook)
		defer zipwriter.Close()

		// write manifest
		fmt.Println("# add manifest")
		zmanifest, err := zipwriter.Create("manifest.json")
		if err != nil {
			log.Fatal(err)
		}

		manifest, _ := os.Open("manifest.json")
		defer manifest.Close()
		_, err = io.Copy(zmanifest, manifest)
		if err != nil {
			log.Fatal(err)
		}

		// write cover
		if _, err := os.Stat("cover.jpg"); err == nil {
			fmt.Println("# add cover")
			zcover, err := zipwriter.Create("cover.jpg")
			if err != nil {
				log.Fatal(err)
			}

			cover, _ := os.Open("cover.jpg")
			defer manifest.Close()
			_, err = io.Copy(zcover, cover)
			if err != nil {
				log.Fatal(err)
			}
		}

		// write tracks (STORED)
		tracks, err := filepath.Glob("*.mp3")
		if err != nil {
			log.Fatal(err)
		}
		natsort.Sort(tracks)

		fmt.Println("# add tracks")
		for _, track := range tracks {

			mp3, _ := os.Open(track)
			defer mp3.Close()

			mp3info, _ := mp3.Stat()
			header, err := zip.FileInfoHeader(mp3info)
			if err != nil {
				log.Fatal(err)
			}
			header.Method = zip.Store
			zmp3, err := zipwriter.CreateHeader(header)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(zmp3, mp3)
			if err != nil {
				log.Fatal(err)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
}
