package utils

import (
	"fmt"
	"log"
	"os/exec"
)

func GetCover(videoPath string, imgPath string, time string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", time, "-frames:v", "1", imgPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(fmt.Sprintf("cmd run failed: %v, output: %s", err, string(output)))
	}
	return err
}
