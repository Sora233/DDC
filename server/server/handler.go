package server

import (
	"github.com/Sora233/DDC/db"
	"github.com/Sora233/DDC/server/middleware"
	"github.com/Sora233/DDC/server/model"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"time"
)

// @title DDC API
// @version 1.0
// @description DDC API
// @termsOfService http://swagger.io/terms/

// @contact.name Sora233
// @contact.url https://github.com/Sora233
// @contact.email sora@sora233.me

// @license.name MIT

// Ping
// @Summary      ping server
// @Description  ping server
// @Tags         api
// @Produce      json
// @Success      200 {object} model.PingResp
// @Router       /v1/ping [get]
func Ping(c *gin.Context) {
	var (
		req  = new(model.PingReq)
		resp = new(model.PingResp)
		log  = middleware.GetLog(c)
	)
	defer func() {
		middleware.UpdateLog(c, log)
	}()

	middleware.SetResp(c, resp)
	if err := c.ShouldBind(req); err != nil {
		resp.ErrPrint(model.DDCErrInvalidArgument, err)
		return
	}
	resp.Pong = "pong"
}

// GetSearchDanmu
// @Summary      查询弹幕
// @Description  查询弹幕
// @Tags         api
// @Produce      json
// @Param	     uid query int64 true "要查询的uid"
// @Success      200 {object} model.GetSearchDanmuResp
// @Router       /v1/danmu [get]
func GetSearchDanmu(c *gin.Context) {
	var (
		req  = new(model.GetSearchDanmuReq)
		resp = new(model.GetSearchDanmuResp)
		log  = middleware.GetLog(c)
	)
	defer func() {
		middleware.UpdateLog(c, log)
	}()

	middleware.SetResp(c, resp)
	if err := c.ShouldBind(req); err != nil {
		resp.ErrPrint(model.DDCErrInvalidArgument, err)
		return
	}

	var comments []*db.Comment

	err := db.G.Where(
		"user_id = ? and time >= ?",
		req.Uid,
		time.Now().Add(-time.Hour*24*3),
	).Order("time desc").Limit(30).Find(&comments).Error

	if err != nil {
		resp.Error(model.DDCErrInternal)
		return
	}

	if len(comments) == 0 {
		return
	}

	resp.Username = comments[0].Username

	vtbuids := lo.Map(comments, func(t *db.Comment, i int) int64 {
		return t.VTBUid
	})
	vtbuids = lo.Uniq(vtbuids)

	var vtbinfos []*db.VTBInfo
	if err := db.G.Where("bilibili_uid in ?", vtbuids).Find(&vtbinfos).Error; err != nil {
		resp.ErrPrint(model.DDCErrInternal, "search vtbinfo error")
		return
	}

	vtbinfoMap := lo.KeyBy(vtbinfos, func(v *db.VTBInfo) int64 {
		return v.BilibiliUid
	})

	for _, cm := range comments {
		resp.Damus = append(resp.Damus, &model.Danmu{
			VTBUid:  cm.VTBUid,
			VTBName: vtbinfoMap[cm.VTBUid].Name,
			Time:    cm.Time.Format("2006-01-02 15:04:05"),
			Text:    cm.Text,
		})
	}

}

// GetSearchMoney
// @Summary      查询消费
// @Description  查询消费
// @Tags         api
// @Produce      json
// @Param	     uid query int64 true "要查询的uid"
// @Success      200 {object}  model.GetSearchMoneyResp
// @Router       /v1/money [get]
func GetSearchMoney(c *gin.Context) {
	var (
		req  = new(model.GetSearchMoneyReq)
		resp = new(model.GetSearchMoneyResp)
		log  = middleware.GetLog(c)
	)
	defer func() {
		middleware.UpdateLog(c, log)
	}()

	middleware.SetResp(c, resp)
	if err := c.ShouldBind(req); err != nil {
		resp.ErrPrint(model.DDCErrInvalidArgument, err)
		return
	}

	var (
		superchats []*db.SuperChat
		gifts      []*db.Gift
	)

	err := db.G.Where(
		"user_id = ? and time >= ?",
		req.Uid,
		time.Now().Add(-time.Hour*24*3),
	).Order("time desc").Limit(30).Find(&superchats).Error

	if err != nil {
		resp.Error(model.DDCErrInternal)
		return
	}

	err = db.G.Where(
		"user_id = ? and time >= ?",
		req.Uid,
		time.Now().Add(-time.Hour*24*3),
	).Order("time desc").Limit(30).Find(&gifts).Error

	if err != nil {
		resp.Error(model.DDCErrInternal)
		return
	}

	if len(superchats) == 0 && len(gifts) == 0 {
		return
	}

	var vtbuids []int64

	for _, gift := range gifts {
		vtbuids = append(vtbuids, gift.VTBUid)
	}
	for _, sc := range superchats {
		vtbuids = append(vtbuids, sc.VTBUid)
	}
	vtbuids = lo.Uniq(vtbuids)

	var vtbinfos []*db.VTBInfo
	if err := db.G.Where("bilibili_uid in ?", vtbuids).
		Find(&vtbinfos).Error; err != nil {
		resp.ErrPrint(model.DDCErrInternal, "search vtbinfo error")
		return
	}

	vtbinfoMap := lo.KeyBy(vtbinfos, func(v *db.VTBInfo) int64 {
		return v.BilibiliUid
	})

	if len(superchats) > 0 {
		resp.Username = superchats[0].Username
		for _, sc := range superchats {
			resp.SuperChats = append(resp.SuperChats, &model.SuperChat{
				VTBUid:  sc.VTBUid,
				VTBName: vtbinfoMap[sc.VTBUid].Name,
				Time:    sc.Time.Format("2006-01-02 15:04:05"),
				Text:    sc.Text,
				Price:   sc.SuperChatPrice,
			})
		}
	}
	if len(gifts) > 0 {
		resp.Username = gifts[0].Username
		for _, gift := range gifts {
			resp.Gifts = append(resp.Gifts, &model.Gift{
				VTBUid:    gift.VTBUid,
				VTBName:   vtbinfoMap[gift.VTBUid].Name,
				Time:      gift.Time.Format("2006-01-02 15:04:05"),
				GiftName:  gift.GiftName,
				GiftNum:   gift.GiftNum,
				GiftPrice: gift.GiftPrice,
			})
		}
	}

}
