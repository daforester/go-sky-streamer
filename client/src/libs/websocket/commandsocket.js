import { EventEmitter } from 'events';

export default class CommandSocket extends EventEmitter {
  #buffer = [];

  #iAmConnecting = false;

  #iAmClosing = false;

  #socket = {};

  #timer = {
    ping: null,
    reconnect: null,
  };

  #wsuri = '';

  constructor(props) {
    super();

    if (typeof props.wsuri === 'string') {
      this.#wsuri = props.wsuri;
    } else {
      throw new Error('wsuri: string should be provided');
    }

    this.#setupSocket();
  }

  destroy() {
    this.#closeSocket();
    this.#clearTimers();
  }

  #clearTimers() {
    if (this.#timer.ping) {
      clearInterval(this.#timer.ping);
    }
    if (this.#timer.reconnect) {
      clearInterval(this.#timer.reconnect);
    }
  }

  #closeSocket() {
    this.#iAmClosing = true;
    if (this.#socket instanceof WebSocket) {
      this.#socket.close();
      this.#socket = {};
    }
  }

  #pingSocket() {
    if (!(this.#socket instanceof WebSocket)) {
      if (this.#timer.ping) {
        clearInterval(this.#timer.ping);
      }
      return;
    }

    this.#socket.send(`PING ${new Date().getTime()}`);
  }

  send(data) {
    if (this.#iAmConnecting) {
      this.#buffer.push(data);
      return;
    }
    if (typeof data === 'object') {
      this.#socket.send(JSON.stringify(data));
    } else if (typeof data === 'string') {
      this.#socket.send(data);
    } else if (data === null) {
      throw new Error('websocket.send: no data payload');
    } else {
      throw new Error('websocket.send: invalid data type');
    }
  }

  sendBuffer() {
    while (!this.#iAmConnecting && this.#buffer.length > 0) {
      const data = this.#buffer.shift();
      this.send(data);
    }
  }

  #setupSocket() {
    if (this.#iAmConnecting) {
      return;
    }
    this.#iAmConnecting = true;
    this.#socket = new WebSocket(this.#wsuri);

    this.#socket.onopen = () => {
      this.#iAmConnecting = false;
      this.sendBuffer();
      this.emit('sockOpen', {
        connection: this.#socket,
      });
      const commandSocket = this;
      this.#timer.ping = setInterval(() => {
        commandSocket.#pingSocket();
      }, 30000);
    };

    this.#socket.onclose = () => {
      this.#iAmConnecting = false;
      this.emit('sockClose', {
        connection: this.#socket,
      });
      this.#socket = {};
      if (!this.#iAmClosing) {
        this.#clearTimers();
        const commandSocket = this;
        this.#timer.reconnect = setTimeout(() => {
          commandSocket.#setupSocket();
        }, 5000);
      }
    };

    this.#socket.onerror = () => {
      this.#iAmConnecting = false;
      this.emit('sockError', {
        connection: this.#socket,
      });
      this.#socket = {};
      this.#clearTimers();
      const commandSocket = this;
      this.#timer.reconnect = setTimeout(() => {
        commandSocket.#setupSocket();
      }, 5000);
    };

    this.#socket.onmessage = (e) => {
      const msg = e.data;
      if (msg.substr(0, 4) === 'PING') {
        this.#socket.send(msg.replace('PING', 'PONG'));
        return;
      }
      if (msg.substr(0, 4) === 'PONG') {
        return;
      }

      this.emit('data', {
        connection: this,
        data: e.data,
      });
    };
  }
}
