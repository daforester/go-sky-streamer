import { SocketRouter, CommandSocket } from '../libs/websocket';

export default {
    data() {
        return {
            socketRouter: new SocketRouter(),
            websocket: {},
        };
    },
    created() {
        this.websocket = new CommandSocket({ wsuri: this.wsuri });
        this.websocket.on('data', this.readWebSocket);

        this.socketRouter.addJSONHandler('ICE_DATA', (e) => {
            this.$emit('ICE_DATA', { Data: e.Data });
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
