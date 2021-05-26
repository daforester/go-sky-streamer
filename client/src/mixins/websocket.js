import { SocketRouter, CommandSocket } from '../libs/websocket';

export default {
  data() {
    return {
      socketRouter: new SocketRouter(),
      websocket: {},
    };
  },
  created() {
    this.websocket = new CommandSocket({ wsuri: this.wsUri });
    this.websocket.on('data', this.readWebSocket);

    this.socketRouter.addJSONHandler('ICE_DATA', (event) => {
      this.$emit('ICE_DATA', { Data: event.Data });
    });
  },
  destroyed() {
    if (this.socketRouter) {
      this.socketRouter = null;
    }
  },
  methods: {
    readWebSocket(p) {
      if (this.socketRouter) {
        this.socketRouter.readMessage(p);
      }
    },
  },
};
