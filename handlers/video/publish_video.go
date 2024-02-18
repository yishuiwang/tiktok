package video

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"path/filepath"
	"tiktok/config"
	"tiktok/models"
	"tiktok/service/video"
	"time"
)

var (
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
)

// PublishVideoHandler 发布视频，并截取一帧画面作为封面
func PublishVideoHandler(c *gin.Context) {
	rawId, _ := c.Get("user_id")

	userId, ok := rawId.(int64)
	if !ok {
		PublishVideoError(c, "userId parse error")
		return
	}

	title := c.PostForm("title")

	form, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}

	files := form.File["data"]
	for _, file := range files {
		suffix := filepath.Ext(file.Filename)
		if _, ok := videoIndexMap[suffix]; !ok {
			PublishVideoError(c, "unsupported video format")
			continue
		}
		name := fmt.Sprintf("%d_%d", userId, time.Now().Unix())
		filename := name + suffix
		savePath := filepath.Join(config.GetConfig("video.save_path"), filename)
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}

		coverPath := filepath.Join(config.GetConfig("video.save_path"), "cover_"+name+".jpg")
		cmd := exec.Command("ffmpeg", "-i", savePath, "-ss", "00:00:01", "-vframes", "1", coverPath)
		err = cmd.Run()
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		err = video.PostVideo(userId, filename, "cover_"+name, title)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		PublishVideoOk(c, file.Filename+" upload success")
	}
}

func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.CommonResponse{StatusCode: 1,
		StatusMsg: msg})
}

func PublishVideoOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.CommonResponse{StatusCode: 0, StatusMsg: msg})
}
