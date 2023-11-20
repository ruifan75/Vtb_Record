package provytdlp

import (
	"strconv"

	"github.com/fzxiao233/Vtb_Record/config"
	"github.com/fzxiao233/Vtb_Record/live/interfaces"
	"github.com/fzxiao233/Vtb_Record/live/videoworker/downloader/provbase"
	"github.com/fzxiao233/Vtb_Record/utils"
	log "github.com/sirupsen/logrus"
)

// func addStreamlinkProxy(co []string, proxy string) []string {
// 	co = append(co, "--http-proxy", "socks5://"+proxy)
// 	return co
// }

type DownloaderStreamlink struct {
	provbase.Downloader
}

func (d *DownloaderStreamlink) StartDownload(video *interfaces.VideoInfo, proxy string, cookie string, filepath string, dirpath string, retryCounter int) error {
	var arg []string
	outputOption := "[%(upload_date)s] %(title)s (%(id)s)[%(uploader)s].%(ext)s"
	if retryCounter > 0 {
		outputOption = "[%(upload_date)s] %(title)s (%(id)s)[%(uploader)s]_" + strconv.Itoa(retryCounter) + ".%(ext)s"
	}
	arg = append(arg, video.Target)
	arg = append(arg, []string{"-i", "--add-metadata", "--embed-thumbnail", "--live-from-start",
		"-o", outputOption, "-P", dirpath}...)
	if config.Config.YtdlpCookies != "" {
		arg = append(arg, []string{"--cookies-from-browser", config.Config.YtdlpCookies}...)
	}
	logger := log.WithField("video", video)
	logger.Infof("start to download %s, command yt-dlp %s", filepath, arg)
	err := utils.ExecShell("yt-dlp", arg...)
	return err
}
