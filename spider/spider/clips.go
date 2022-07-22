package spider

import (
	"fmt"
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
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	}).BindJSON(resp).Do()
	return resp, err
}
