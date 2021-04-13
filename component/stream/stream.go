package stream

import "github.com/pion/webrtc/v3"

type Stream struct {
	Connection *webrtc.PeerConnection
	VideoTrack *webrtc.TrackLocalStaticSample
}
