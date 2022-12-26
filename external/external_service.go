package external

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/noornee/reddit-dl/handler"
	"github.com/noornee/reddit-dl/utility"
)

var temp_dir string = utility.CreateDir()

func Setup(media_url, audio_url, title string) {

	if audio_url != "" {
		status_code, mime := handler.GetHead(audio_url)

		if status_code == 200 && !strings.Contains(mime, "image") {
			aria2c(media_url, audio_url)
			ffmpeg(title)
			return
		}

	}

	aria2c_nos(media_url, title)

}

// download files[video,audio] with aria2c
func aria2c(media_url, audio_url string) {

	cmd := exec.Command("aria2c", "-d", temp_dir, "-Z", media_url, audio_url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

}

// download files[video/gif/image](files with no sound) with aria2c
func aria2c_nos(media_url, title string) {

	var cmd *exec.Cmd

	_, mime_type := handler.GetHead(media_url)

	switch mime_type {
	case "image/jpeg":
		cmd = exec.Command("aria2c", media_url, "-o", title+".jpg")
	case "image/png":
		cmd = exec.Command("aria2c", media_url, "-o", title+".png")
	case "image/gif":
		cmd = exec.Command("aria2c", media_url, "-o", title+".gif")
	default:
		cmd = exec.Command("aria2c", media_url, "-o", title+".mp4")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		utility.ErrorLog.Println(err)
	}

}

// merge downoladed files[video,audio] with ffmpeg
func ffmpeg(filename string) {

	filename = filename + ".mp4"

	files, err := ioutil.ReadDir(temp_dir)
	if err != nil {
		utility.ErrorLog.Println(err)
	}

	var aud, vid string

	for range files {
		vid = temp_dir + "/" + files[0].Name()
		aud = temp_dir + "/" + files[1].Name()
	}

	cmd := exec.Command("ffmpeg", "-y", "-v", "quiet", "-stats", "-i", vid, "-i", aud, filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	utility.InfoLog.Printf("Merging files into \t%s", filename)

	if err := cmd.Run(); err != nil {
		utility.ErrorLog.Println(err)
	}
	utility.InfoLog.Println("Done")

}
