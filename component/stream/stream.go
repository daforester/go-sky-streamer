package stream

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"github.com/daforester/go-sky-streamer/component/capture"
	"github.com/pion/webrtc/v3"
	"io/ioutil"
)

type Stream struct {
	Connection *webrtc.PeerConnection
	VideoTrack *webrtc.TrackLocalStaticSample
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

	s.Connection = peerConnection

	return s
}

func (S *Stream) AddOffer(offer string) {
	var err error

	offer := webrtc.SessionDescription{}
	decode(offer, &offer)

	// Set the remote SessionDescription
	if err = S.Connection.SetRemoteDescription(offer); err != nil {
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