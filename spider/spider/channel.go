package spider

import (
	"github.com/Sora233/DDC/spider/config"
	"github.com/guonaihong/gout"
)

const channelUrl = "https://api.matsuri.icu/channel"

type ChannelResp struct {
	Status int        `json:"status"`
	Data   []*VTBInfo `json:"data"`
}

type VTBInfo struct {
	Archive          bool   `json:"archive"`
	BilibiliLiveRoom int64  `json:"bilibili_live_room"`
	BilibiliUid      int64  `json:"bilibili_uid"`
	Face             string `json:"face"`
	Hidden           bool   `json:"hidden"`
	IsLive           bool   `json:"is_live"`
	LastDanmu        int64  `json:"last_danmu"`
	LastLive         int64  `json:"last_live"`
	Name             string `json:"name"`
	TotalClips       int    `json:"total_clips"`
	TotalDanmu       int    `json:"total_danmu"`
}

func GetChannel() (*ChannelResp, error) {
	checkLimit()
	var resp = new(ChannelResp)
	err := gout.GET(channelUrl).SetHeader(gout.H{
		"user-agent": config.Global.UserAgent,
	}).BindJSON(resp).Do()
	return resp, err
}
