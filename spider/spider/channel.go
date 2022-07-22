package spider

import (
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
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	}).BindJSON(resp).Do()
	return resp, err
}
