package spider

import (
	"context"
	"errors"
	"github.com/Sora233/DDC/db"
	"github.com/Sora233/DDC/spider/config"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"sync"
	"time"
)

var limiter = rate.NewLimiter(rate.Every(time.Second*5), 1)

func checkLimit() {
	_ = limiter.Wait(context.Background())
}

func checkTime(msec int64) bool {
	return !time.UnixMilli(msec).Add(time.Hour * 24 * 1).Before(time.Now())
}

func Run() {
	limiter.SetLimit(rate.Every(time.Second * time.Duration(config.Global.RequestLimit)))
	var vtbinfoChan = make(chan *VTBInfo, 4)
	var clipChan = make(chan *Clip, 4)
	var clipWg sync.WaitGroup
	var commentWg sync.WaitGroup

	var clipProcessorNum = config.Global.ClipProcessorNum
	var commentProcessorNum = config.Global.CommentProcessorNum

	for i := 0; i < clipProcessorNum; i++ {
		clipWg.Add(1)
		go ClipProcessor(&clipWg, vtbinfoChan, clipChan)
	}

	for i := 0; i < commentProcessorNum; i++ {
		commentWg.Add(1)
		go CommentProcessor(&commentWg, clipChan)
	}

	channelResp, err := GetChannel()
	if err != nil {
		logrus.Fatalf("GetChannel error %v", err)
	}
	if channelResp.Status != 0 {
		logrus.Fatalf("ChannelReps status %v", channelResp.Status)
	}

	lo.ForEach(channelResp.Data, func(info *VTBInfo, idx int) {
		var dbvtbid int64
		err := db.G.Model((*db.VTBInfo)(nil)).
			Select("id").
			Where("bilibili_uid = ?", info.BilibiliUid).
			Find(&dbvtbid).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dbvtb := &db.VTBInfo{
				Archive:          info.Archive,
				BilibiliLiveRoom: info.BilibiliLiveRoom,
				BilibiliUid:      info.BilibiliUid,
				Face:             info.Face,
				Hidden:           info.Hidden,
				IsLive:           info.IsLive,
				LastDanmu:        info.LastDanmu,
				LastLive:         info.LastLive,
				Name:             info.Name,
				TotalClips:       info.TotalClips,
				TotalDanmu:       info.TotalDanmu,
			}
			db.G.Save(dbvtb)
		} else if err != nil {
			logrus.Errorf("check VTBInfo error %v", err)
		}
	})

	datas := lo.Filter(channelResp.Data, func(info *VTBInfo, idx int) bool {
		log := logrus.WithField("uid", info.BilibiliUid).
			WithField("name", info.Name)
		if info.Archive {
			log.Debug("skip achive")
			return false
		}
		if !checkTime(info.LastLive) {
			log.Debug("skip not active")
			return false
		}
		if info.IsLive {
			log.Debug("skip living")
			return false
		}
		var err error
		var dbclip *db.Clip
		err = db.G.Order("end_time desc").Limit(1).
			Find(&dbclip, "vtb_uid = ?", info.BilibiliUid).Error
		if err == nil && dbclip.StartTime.Equal(time.UnixMilli(info.LastLive)) {
			log.Debug("skip not new clips")
			return false
		}

		var dbvtbid int64
		err = db.G.Model((*db.VTBInfo)(nil)).
			Select("id").
			Where("bilibili_uid = ?", info.BilibiliUid).
			Find(&dbvtbid).Error
		if err != nil {
			log.Errorf("get VTBInfo error %v", err)
		} else {
			dbvtb := &db.VTBInfo{
				ID:               dbvtbid,
				Archive:          info.Archive,
				BilibiliLiveRoom: info.BilibiliLiveRoom,
				BilibiliUid:      info.BilibiliUid,
				Face:             info.Face,
				Hidden:           info.Hidden,
				IsLive:           info.IsLive,
				LastDanmu:        info.LastDanmu,
				LastLive:         info.LastLive,
				Name:             info.Name,
				TotalClips:       info.TotalClips,
				TotalDanmu:       info.TotalDanmu,
			}
			db.G.Save(dbvtb)
		}
		return true
	})

	if len(datas) == 0 {
		close(vtbinfoChan)
		clipWg.Wait()
		close(clipChan)
		commentWg.Wait()
		return
	}

	logrus.Infof("total idx: %v", len(datas))
	for _, info := range datas {
		vtbinfoChan <- info
	}
	close(vtbinfoChan)
	clipWg.Wait()
	logrus.Info("all vtb done")
	close(clipChan)
	commentWg.Wait()
	logrus.Infof("all comment done")
}

func ClipProcessor(wg *sync.WaitGroup, c <-chan *VTBInfo, clipChan chan<- *Clip) {
	defer wg.Done()

	for info := range c {
		log := logrus.WithField("uid", info.BilibiliUid).
			WithField("name", info.Name)
		log.Infof("process start")
		clipsResp, err := GetClips(info.BilibiliUid)
		if err != nil {
			log.Errorf("GetClips error %v", err)
			continue
		}
		if clipsResp.Status != 0 {
			log.Errorf("ClipsResp status %v", clipsResp.Status)
			continue
		}
		slices.SortFunc(clipsResp.Clips, func(a, b *Clip) bool {
			return a.StartTime > b.StartTime
		})
		for _, clip := range clipsResp.Clips {
			log := log.WithField("clip_id", clip.ID).
				WithField("start_time", time.UnixMilli(clip.StartTime))
			if !checkTime(clip.StartTime) {
				log.Debug("skip old clip")
				break
			}
			if clip.EndTime == 0 {
				log.Debug("skip living clip")
				continue
			}

			var exist bool
			if err := db.G.Model((*db.Clip)(nil)).
				Select("count(*) > 0").Where("clip_id = ?", clip.ID).Find(&exist).Error; err != nil {
				log.Errorf("check clip cache error %v", err)
			} else if exist {
				log.Debug("skip cache hit")
				continue
			}

			clipChan <- clip
		}
	}

}

func CommentProcessor(wg *sync.WaitGroup, c <-chan *Clip) {
	defer wg.Done()

	for clip := range c {
		log := logrus.WithField("uid", clip.BilibiliUid).
			WithField("clip_id", clip.ID).
			WithField("start_time", time.UnixMilli(clip.StartTime))

		commentsResp, err := GetComments(clip.ID)
		if err != nil {
			log.Errorf("GetComments error %v", err)
			continue
		}
		if commentsResp.Status != 0 {
			log.Errorf("CommentsResp status %v", commentsResp.Status)
			continue
		}
		var cms []*db.Comment
		var sc []*db.SuperChat
		var gift []*db.Gift
		for _, comment := range commentsResp.Comments {
			if comment.SuperchatPrice != 0 {
				sc = append(sc, &db.SuperChat{
					VTBUid:         clip.BilibiliUid,
					ClipID:         clip.ID,
					UserID:         comment.UserID,
					Username:       comment.Username,
					Time:           time.UnixMilli(comment.Time),
					Text:           comment.Text,
					SuperChatPrice: comment.SuperchatPrice,
				})
			} else if len(comment.GiftName) > 0 {
				gift = append(gift, &db.Gift{
					VTBUid:    clip.BilibiliUid,
					ClipID:    clip.ID,
					UserID:    comment.UserID,
					Username:  comment.Username,
					Time:      time.UnixMilli(comment.Time),
					GiftName:  comment.GiftName,
					GiftNum:   comment.GiftNum,
					GiftPrice: comment.GiftPrice,
				})
			} else {
				cms = append(cms, &db.Comment{
					VTBUid:   clip.BilibiliUid,
					ClipID:   clip.ID,
					UserID:   comment.UserID,
					Username: comment.Username,
					Time:     time.UnixMilli(comment.Time),
					Text:     comment.Text,
				})
			}
		}
		if err := db.G.Create(cms).Error; err != nil {
			log.Errorf("create comments error %v", err)
		}
		if err := db.G.Create(gift).Error; err != nil {
			log.Errorf("create gift error %v", err)
		}
		if err := db.G.Create(sc).Error; err != nil {
			log.Errorf("create sc error %v", err)
		}

		if err := db.G.Create(&db.Clip{
			ClipID:     clip.ID,
			VTBUid:     clip.BilibiliUid,
			Cover:      clip.Cover,
			StartTime:  time.UnixMilli(clip.StartTime),
			EndTime:    time.UnixMilli(clip.EndTime),
			Title:      clip.Title,
			TotalDanmu: clip.TotalDanmu,
			Views:      clip.Views,
		}).Error; err != nil {
			log.Errorf("create clip error %v", err)

		}

		log.Infof("clip done")
	}
}
