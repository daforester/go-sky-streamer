<template>
  <div ref="videoStream"></div>
</template>
<script>

import websocket from '../mixins/websocket';

export default {
  name: 'Main',
  mixins: [
    websocket,
  ],
  data() {
    return {
      peer: null,
    };
  },
  mounted() {
    this.on('ICE_DATA', (event) => {
      try {
        this.peer.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(event.Data))));
      } catch (err) {
        console.log(err);
      }
    });
  },
  methods: {
    setupRTC() {
      const pc = new RTCPeerConnection({
        iceServers: [
          {
            urls: 'stun:stun.l.google.com:19302',
          },
        ],
      });
      this.peer = pc;

      pc.ontrack = (event) => {
        const el = document.createElement(event.track.kind);
        [el.srcObject] = event.streams;
        el.autoplay = true;
        el.controls = true;
        this.$refs.videoStream.appendChild(el);
      };

      pc.oniceconnectionstatechange = (event) => {
        console.log(event);
      };

      pc.onicecandidate = (event) => {
        if (event.candidate === null) {
          console.log(JSON.stringify(pc.localDescription));
        }
      };
    },
  },
};
</script>
<style>

</style>
