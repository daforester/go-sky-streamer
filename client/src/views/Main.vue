<template>
  <div ref="videoStream"></div>
</template>
<script>

import websocket from "../mixins/websocket";

export default {
  name: 'Main',
  mixins: [
    websocket,
  ],
  data() {
    return {
      peer: null,
    }
  },
  mounted() {
    this.on('ICE_DATA', (e) => {
      try {
        this.peer.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(e.Data))));
      } catch (e) {
        console.log(e);
      }
    });
  },
  methods: {
    setupRTC() {
      const pc = new RTCPeerConnection({
        iceServers: [
          {
            urls: 'stun:stun.l.google.com:19302',
          }
        ],
      });
      this.peer = pc;

      pc.ontrack = (e) => {
        const el = document.createElement(e.track.kind);
        el.srcObject = e.streams[0];
        el.autoplay = true;
        el.controls = true;
        this.$refs.videoStream.appendChild(el);
      }

      pc.oniceconnectionstatechange = (e) => {
        console.log(e);
      }

      pc.onicecandidate = (e) => {
        if (e.candidate === null) {
          console.log(JSON.stringify(pc.localDescription));
        }
      }
    }
  }
};
</script>
<style>

</style>
