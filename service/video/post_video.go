package video

import "tiktok/models"

type PostVideoFlow struct {
	videoName string
	coverName string
	title     string
	userId    int64
	video     *models.Video
}

// PostVideo 投稿视频
func PostVideo(userId int64, videoName, coverName, title string) error {
	return NewPostVideoFlow(userId, videoName, coverName, title).Do()
}

func NewPostVideoFlow(userId int64, videoName, coverName, title string) *PostVideoFlow {
	return &PostVideoFlow{userId: userId, videoName: videoName, coverName: coverName, title: title}
}

func (p *PostVideoFlow) Do() error {
	p.prepareParam()

	if err := p.publish(); err != nil {
		return err
	}

	return nil
}

func (p *PostVideoFlow) prepareParam() {
	// TODO
}

func (p *PostVideoFlow) publish() error {
	video := &models.Video{
		UserInfoId: p.userId,
		PlayUrl:    p.videoName,
		CoverUrl:   p.coverName,
		Title:      p.title,
	}
	return models.NewVideoDAO().AddVideo(video)
}
