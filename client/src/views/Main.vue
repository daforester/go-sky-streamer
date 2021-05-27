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
    this.setupRTC();
    this.$on('ICE_DATA', (event) => {
      try {
        console.log('ICE_DATA event');
        console.log(event);
        console.log(atob(event.Data.Offer));
        this.peer.setRemoteDescription(
          new RTCSessionDescription(JSON.parse(atob(event.Data.Offer))),
        );
      } catch (err) {
        console.log('err');
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
        console.log('ontrack');
        const el = document.createElement(event.track.kind);
        [el.srcObject] = event.streams;
        el.autoplay = true;
        el.controls = true;
        this.$refs.videoStream.appendChild(el);
      };

      pc.oniceconnectionstatechange = (event) => {
        console.log('oniceconnectionstatechange event');
        console.log(event);
      };

      pc.onicecandidate = (event) => {
        if (event.candidate === null) { // Finished Gathering
          console.log(JSON.stringify(pc.localDescription));
          const statusUpdate = {
            Command: 'GET_ICE',
            Data: {
              Offer: btoa(JSON.stringify(pc.localDescription)),
            },
          };

          this.websocket.send(JSON.stringify(statusUpdate));
        }
      };

      pc.addTransceiver('audio', { direction: 'sendrecv' });
      pc.addTransceiver('video', { direction: 'sendrecv' });

      pc.createOffer().then((localDescription) => {
        console.log('setLocalDescription');
        pc.setLocalDescription(localDescription);
      }).catch((err) => {
        console.log('err');
        console.log(err);
      });
    },
  },
};
</script>
<style>

</style>
