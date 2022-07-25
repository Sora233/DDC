package spider

import (
	"fmt"
	"github.com/Sora233/DDC/spider/config"
	"github.com/guonaihong/gout"
)

const clipsUrl = "https://api.matsuri.icu/channel/%v/clips"

type Clip struct {
	ID          string `json:"id"`
	BilibiliUid int64  `json:"bilibili_uid"`
	Cover       string `json:"cover"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	Title       string `json:"title"`
	TotalDanmu  int    `json:"total_danmu"`
	Views       int    `json:"views"`
}

type ClipsResp struct {
	Status int     `json:"status"`
	Clips  []*Clip `json:"data"`
}

func GetClips(uid int64) (*ClipsResp, error) {
	checkLimit()
	var resp = new(ClipsResp)
	err := gout.GET(fmt.Sprintf(clipsUrl, uid)).SetHeader(gout.H{
		"user-agent": config.Global.UserAgent,
	}).BindJSON(resp).Do()
	return resp, err
}
