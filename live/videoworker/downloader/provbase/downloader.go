package provbase

/*
	Contains common functions & base types for downloaders
*/

import (
	"github.com/fzxiao233/Vtb_Record/live/interfaces"
	log "github.com/sirupsen/logrus"
)

type DownloadProvider interface {
	StartDownload(video *interfaces.VideoInfo, proxy string, cookie string, filepath string, dirpath string, retryCounter int) error
}
type Downloader struct {
	Prov DownloadProvider
}

func (d *Downloader) DownloadVideo(video *interfaces.VideoInfo, proxy string, cookie string, filePath string, dirpath string, retryCounter int) string {
	logger := log.WithField("video", video)
	logger.Infof("start to download")
	video.FilePath = filePath
	err := d.Prov.StartDownload(video, proxy, cookie, filePath, dirpath, retryCounter)
	logger.Infof("finished with status: %s", err)
	if err != nil {
		logger.Infof("download failed: %s", err)
		return ""
	}
	logger.Infof("%s download successfully", filePath)
	return filePath
}
