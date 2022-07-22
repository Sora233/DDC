package spider

import (
	"fmt"
	"github.com/guonaihong/gout"
)

const commentsUrl = "https://api.matsuri.icu/clip/%v/comments"

type Comment struct {
	Text string `json:"text"`

	SuperchatPrice int `json:"superchat_price"`

	GiftName  string  `json:"gift_name"`
	GiftNum   int     `json:"gift_num"`
	GiftPrice float64 `json:"gift_price"`

	Time     int64  `json:"time"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type CommentsResp struct {
	Status   int        `json:"status"`
	Comments []*Comment `json:"data"`
}

func GetComments(clipId string) (*CommentsResp, error) {
	checkLimit()
	var resp = new(CommentsResp)
	err := gout.GET(fmt.Sprintf(commentsUrl, clipId)).SetHeader(gout.H{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	}).BindJSON(resp).Do()
	return resp, err
}
