const handlerMethod = {
  COMMAND_HANDLER: 0,
  JSON_HANDLER: 1,
};

export default class SocketRouter {
  #socketHandlers;

  constructor() {
    this.#socketHandlers = new Map();
  }

  addCommandHandler(command, ...handlers) {
    return this.addHandler(command, handlerMethod.COMMAND_HANDLER, ...handlers);
  }

  addHandler(command, method, ...handlers) {
    const uCommand = command.toUpperCase();

    const re = /^[A-Z-_]+$/;
    if (!uCommand.match(re)) {
      throw new Error('command may only contain A-Z, - or _');
    }

    handlers.forEach((handlerFunc) => {
      const h = {
        method,
        handlerFunc,
      };
      if (!this.#socketHandlers[uCommand]) {
        this.#socketHandlers[uCommand] = [h];
      } else {
        this.#socketHandlers[uCommand].push(h);
      }
    });

    return null;
  }

  addJSONHandler(command, ...handlers) {
    return this.addHandler(command, handlerMethod.JSON_HANDLER, ...handlers);
  }

  static parseData(input) {
    let r;
    let err;
    if (input.substr(0, 1) === '{') {
      [r, err] = SocketRouter.parseJSONData(input);
      if (!err) {
        return [handlerMethod.JSON_HANDLER, r];
      }
    }
    r = SocketRouter.parseCommandData(input);
    return [handlerMethod.COMMAND_HANDLER, r];
  }

  static parseJSONData(input) {
    try {
      const r = JSON.parse(input);
      return [r, null];
    } catch (e) {
      return [null, e];
    }
  }

  static parseCommandData(input) {
    const r = {
      Command: '',
      Data: '',
      Params: '',
    };
    const trimStr = input.trim();
    const i = trimStr.indexOf(' ');
    if (i === -1) {
      r.Command = trimStr;
      return r;
    }
    r.Command = trimStr.substr(0, i);
    if (trimStr.length > i + 1) {
      r.Data = trimStr.substr(i + 1);
    }
    return r;
  }

  readMessage(packet) {
    const [rType, r] = SocketRouter.parseData(packet.data);
    const handlers = this.#socketHandlers[r.Command];

    if (handlers === undefined || handlers === null || handlers.length === 0) {
      return;
    }

    handlers.forEach((h) => {
      if (h.method !== rType) {
        return;
      }
      const c = {
        Connection: packet.connection,
        RequestData: packet.data,
        Data: r.Data,
        Params: r.Params,
      };
      h.handlerFunc(c);
    });
  }
}
