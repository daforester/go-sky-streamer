package capture

import (
	"errors"
	"fmt"
	"github.com/daforester/go-sky-streamer/component/stream"
	"github.com/pion/webrtc/v3/pkg/media"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Capture struct {
	Name string `json:"name"`
	DevicePath string `json:"device_path"`
	Framebuffer chan []byte
	Height uint
	InputFormat string
	Width uint
	on bool
	off chan bool
}

func (C *Capture) Start() error {
	C.DevicePath = "@device_pnp_\\\\?\\usb#vid_04ca&pid_707f&mi_00#6&306659e8&0&0000#{65e8773d-8f56-11d0-a3b9-00a0c9223196}\\global"
	C.InputFormat = "dshow"
	// C.InputFormat = "v4l2"
	C.Height = 1080
	C.Width = 1920

	cmd := exec.Command("ffmpeg", "-framerate", "30", "-f", C.InputFormat, "-input_format", "h264", "-video_size", fmt.Sprintf("%vx%v", C.Width, C.Height), "-i", c.DevicePath, "-c", "copy", "-f", "h264", "pipe:1")
	fmt.Println(cmd.Args)

	dataPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("could not create named pipe. ", err)
	}

	if err := execCmd(cmd); err != nil {
		return err
	}

	C.Framebuffer = make(chan []byte, 60)

	go func() {
		for {
			select {
			case <-C.off:
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					log.Println("failed to kill camera process. ", err)
				}
				return
			default:
				frameBytes := make([]byte, 600000)
				n, err := dataPipe.Read(frameBytes)
				if err != nil {
					log.Println("could not read pipe. ", err)
				}

				C.Framebuffer <- frameBytes[:n]
			}
		}
	}()

	C.off = make(chan bool)

	C.on = true
	return nil
}

func (C *Capture) AttachTrack(stream *stream.Stream) {
	go func() {
		for {
			select {
			case <-C.off:
				_ = stream.Connection.Close()
				return
			case f := <-C.Framebuffer:
				sample := media.Sample{
					Data:    f,
				}

				if err := stream.VideoTrack.WriteSample(sample); err != nil {
					log.Fatal("could not write rtp sample. ", err)
					return
				}
			}
		}
	}()
}

func execCmd(cmd *exec.Cmd) error {
	logFile, err := os.OpenFile("ffmpeg_log.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0775)
	if err != nil {
		return errors.New("could not create ffmpeg log. " + err.Error())
	}

	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
