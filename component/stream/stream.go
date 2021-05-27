package stream

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/daforester/go-sky-streamer/component/capture"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"io/ioutil"
	"log"
)

type Stream struct {
	Connection *webrtc.PeerConnection
}

func (S Stream) New(capture *capture.Capture) *Stream {
	s := new(Stream)

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateConnected {
			localTrack, err := webrtc.NewTrackLocalStaticSample(
				webrtc.RTPCodecCapability{MimeType: "video/h264"},
				"video",
				"pion",
			)
			if err != nil {
				panic(err)
			}

			for {
				select {
				case <-capture.Off:
					_ = S.Connection.Close()
					return
				case f := <-capture.Framebuffer:
					sample := media.Sample{
						Data:    f,
					}

					if err := localTrack.WriteSample(sample); err != nil {
						log.Fatal("could not write rtp sample. ", err)
						return
					}
				}
			}
		}
	})

	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Println("ON TRACK")
		localTrack, err := webrtc.NewTrackLocalStaticSample(remoteTrack.Codec().RTPCodecCapability, "video", "pion")
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-capture.Off:
				_ = S.Connection.Close()
				return
			case f := <-capture.Framebuffer:
				sample := media.Sample{
					Data:    f,
				}

				if err := localTrack.WriteSample(sample); err != nil {
					log.Fatal("could not write rtp sample. ", err)
					return
				}
			}
		}
	})

	s.Connection = peerConnection

	return s
}

func (S *Stream) AddOffer(offer string) string {
	var err error

	sessionDescription := new(webrtc.SessionDescription)
	decode(offer, sessionDescription)

	// Set the remote SessionDescription
	if err = S.Connection.SetRemoteDescription(*sessionDescription); err != nil {
		panic(err)
	}

	// Create answer
	answer, err := S.Connection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(S.Connection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = S.Connection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	<-gatherComplete

	return encode(S.Connection.LocalDescription())
}

// Allows compressing offer/answer to bypass terminal input limits.
const compress = false

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if compress {
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}

func zip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}
	err = gz.Flush()
	if err != nil {
		panic(err)
	}
	err = gz.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}