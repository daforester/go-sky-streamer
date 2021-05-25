import createPersistedState from 'vuex-persistedstate';
import isElectron from './libs/isElectron';

export const config = {
  mutations: {
    audioSettings(state, payload) {
      if (typeof payload === 'object') {
        state.audioSettings = payload;
      } else {
        throw new Error('audioSettings: Payload is not object');
      }
    },
    autoLogin(state, payload) {
      if (typeof payload === 'boolean') {
        state.autoLogin = payload;
      } else {
        throw new Error('autoLogin: Payload is not boolean');
      }
    },
    callHistory(state, payload) {
      // 30 Days
      if (Array.isArray(payload)) {
        state.callHistory = payload;
      } else if (typeof payload === 'object') {
        state.callHistory.push(payload);
      } else {
        throw new Error('callHistory: Payload is not array or object');
      }
    },
    clearLocalStorage(state, payload) {
      if (typeof payload === 'boolean') {
        if (!isElectron()) {
          state.clearLocalStorage = payload;
          sessionStorage.setItem('clearLocalStorage', payload);
        } else {
          state.clearLocalStorage = false;
          sessionStorage.setItem('clearLocalStorage', false);
        }
      } else {
        throw new Error('clearLocalStorage: Payload is not boolean');
      }
    },
    currentPanel(state, payload) {
      if (typeof payload === 'string') {
        state.currentPanel = payload;
      } else {
        throw new Error('clearLocalStorage: Payload is not string');
      }
    },
    expandedAddressBooks(state, payload) {
      let id = 0;
      if (Array.isArray(payload)) {
        state.expandedAddressBooks = payload;
        return;
      }
      if (typeof payload === 'object') {
        id = payload.ID;
      } else if (typeof payload === 'number') {
        id = payload;
      } else {
        throw new Error('expandedAddressBooks: Payload must be array, address book object or number');
      }
      if (id > 0) {
        const index = state.expandedAddressBooks.indexOf(id);
        if (index >= 0) {
          state.expandedAddressBooks.splice(index, 1);
        } else {
          state.expandedAddressBooks.push(payload);
        }
      } else {
        throw new Error('expandedAddressBooks: Payload ID must be greater than zero');
      }
    },
    rememberLogin(state, payload) {
      if (typeof payload === 'boolean') {
        state.rememberLogin = payload;
      } else {
        throw new Error('rememberLogin: Payload is not boolean');
      }
    },
    storedCredentials(state, payload) {
      if (typeof payload === 'object') {
        if (payload != null) {
          if (!payload.domain) {
            throw new Error('storedCredentials: Object missing domain');
          } else if (!payload.email) {
            throw new Error('storedCredentials: Object missing email');
          } else if (!payload.password) {
            throw new Error('storedCredentials: Object missing password');
          }
        }
        state.storedCredentials = payload;
      } else {
        throw new Error('storedCredentials: Payload is not object');
      }
    },
    managedTokens(state, payload) {
      state.managedTokens = payload;
    },
  },
  plugins: [
    createPersistedState({
      storage: window.localStorage,
    }),
  ],
  state: {
    audioSettings: {
      headset: {
        speaker: 'communications',
        microphone: 'communications',
        ringer: 'default',
      },
      speakerphone: {
        speaker: 'default',
        microphone: 'default',
        ringer: 'default',
      },
      mode: 'headset',
      volume: 50,
    },
    autoLogin: false,
    callHistory: [],
    clearLocalStorage: !isElectron(),
    currentPanel: 'dialpad',
    expandedAddressBooks: [],
    rememberLogin: false,
    storedCredentials: {
      domain: '', email: '', password: '',
    },
    managedTokens: [],
    version: 1.1,
  },
};

export default config;
