package downloader

import (
	"github.com/fzxiao233/Vtb_Record/live/videoworker/downloader/provbase"
	"github.com/fzxiao233/Vtb_Record/live/videoworker/downloader/provstreamlink"
	"github.com/fzxiao233/Vtb_Record/live/videoworker/downloader/provytdlp"
	log "github.com/sirupsen/logrus"
)

type Downloader = provbase.Downloader

func GetDownloader(providerName string) *Downloader {
	if providerName == "" || providerName == "streamlink" {
		return &Downloader{&provstreamlink.DownloaderStreamlink{}}
	} else if providerName == "ytdlp" {
		return &Downloader{&provytdlp.DownloaderStreamlink{}}
	} else {
		log.Fatalf("Unknown download provider %s", providerName)
		return nil
	}
}
