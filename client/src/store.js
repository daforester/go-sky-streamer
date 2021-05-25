import createPersistedState from 'vuex-persistedstate';

export const config = {
  mutations: {

  },
  plugins: [
    createPersistedState({
      storage: window.localStorage,
    }),
  ],
  state: {

  },
};

export default config;
