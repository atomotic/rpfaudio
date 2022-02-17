# rpfaudio

A small opinionated cli to create [Readium Audiobooks](https://readium.org/webpub-manifest/profiles/audiobook.html) from a directory containing mp3 files.


## install

    go install github.com/atomotic/rpfaudio@latest

## run

Generate a readium manifest (with empty metadata) inside a directory containing mp3 files.  
 If `cover.jpg` available is added to resources.

    rpfaudio init

After manually editing manifest.json package the audiobook

    rpfaudio package