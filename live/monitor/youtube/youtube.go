package youtube

import (
	"fmt"
	"regexp"

	"github.com/fzxiao233/Vtb_Record/config"
	"github.com/fzxiao233/Vtb_Record/live/interfaces"
	"github.com/fzxiao233/Vtb_Record/live/monitor/base"
	. "github.com/fzxiao233/Vtb_Record/utils"
	"github.com/tidwall/gjson"
)

type yfConfig struct {
	IsLive bool
	Title  string
	Target string
}
type Youtube struct {
	base.BaseMonitor
	yfConfig
	usersConfig config.UsersConfig
}

func (y *Youtube) getVideoInfo(ctx *base.MonitorCtx, baseHost string, channelId string) error {
	url := baseHost + "/channel/" + channelId + "/streams"
	htmlBody, err := ctx.HttpGet(url, map[string]string{})
	if err != nil {
		return err
	}
	re, _ := regexp.Compile(`var\sytInitialData\s=\s*([^\n]+?});`)
	result := re.FindSubmatch(htmlBody)
	if len(result) < 2 {
		return fmt.Errorf("youtube cannot find js_data")
	}
	jsonYtConfig := result[1]
	contents := gjson.GetBytes(jsonYtConfig, `contents.twoColumnBrowseResultsRenderer.tabs.#(tabRenderer.title="Live").tabRenderer.content.richGridRenderer.contents`)
	firstLive := contents.Get(`#(richItemRenderer.content.videoRenderer.thumbnailOverlays.0.thumbnailOverlayTimeStatusRenderer.style="LIVE").richItemRenderer.content.videoRenderer`)
	if !firstLive.Exists() {
		return fmt.Errorf("no live exists")
	}
	y.IsLive = true
	y.Title = firstLive.Get("title.runs.0.text").String()
	y.Target = "https://www.youtube.com/watch?v=" + firstLive.Get("videoId").String()
	return nil
}

func (y *Youtube) CreateVideo(usersConfig config.UsersConfig) *interfaces.VideoInfo {
	if !y.yfConfig.IsLive {
		return &interfaces.VideoInfo{}
	}
	v := &interfaces.VideoInfo{
		Title:       y.Title,
		Date:        GetTimeNow(),
		Target:      y.Target,
		Provider:    "Youtube",
		UsersConfig: usersConfig,
	}
	return v
}
func (y *Youtube) CheckLive(usersConfig config.UsersConfig) bool {
	y.usersConfig = usersConfig
	err := y.getVideoInfo(base.GetCtx("Youtube"), "http://www.youtube.com", y.usersConfig.TargetId)
	if err != nil {
		y.IsLive = false
	}
	if !y.IsLive {
		base.NoLiving("Youtube", usersConfig.Name)
	}
	return y.yfConfig.IsLive
}
