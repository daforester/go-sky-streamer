import get from 'lodash.get';
import Vue from 'vue';
import Vuex from 'vuex';
import VueRouter from 'vue-router';
import { routes } from 'routes';
import { config } from 'store';
import WebApp from 'views/WebApp.vue';

export default function SkyStreamer(options = {}, componentConfig = {}) {
  Vue.use(Vuex);
  Vue.use(VueRouter);
  Vue.mixin({
    data() {
      return {
        // eslint-disable-next-line no-undef
        apiUri: componentConfig.apiUri || API_URL,
        // eslint-disable-next-line no-undef
        baseUri: componentConfig.baseUri || BASE_URL,
        // eslint-disable-next-line no-undef
        wsUri: componentConfig.wsUri || WS_URL,
      };
    },
  });

  const store = new Vuex.Store(config);
  const router = new VueRouter({
    // eslint-disable-next-line no-undef
    mode: 'history',
    routes,
  });

  router.beforeEach((to, from, next) => {
    let code = 200;
    const selector = document.querySelector("meta[name='http.status']");

    if (typeof selector !== 'undefined' && selector !== null) {
      const content = selector.getAttribute('content');
      if (typeof content !== 'undefined' && content !== null && content !== '') {
        code = parseInt(content, 10);
      }
    }

    if (to.path !== '/error' && code >= 400) {
      next('/error');
    } else {
      next();
    }
  });

  return new Vue({
    render: (h) => h(WebApp, { props: { config: componentConfig } }),
    router,
    store,
  }).$mount(get(options, 'el', '#webapp'));
}

if ('serviceWorker' in navigator) {
  navigator.serviceWorker
    // eslint-disable-next-line no-undef
    .register('/service-worker.js')
    .then((registration) => {
      registration.onupdatefound = () => {
        if (navigator.serviceWorker.controller) {
          const installingWorker = registration.installing;
          installingWorker.onstatechange = () => {
            switch (installingWorker.state) {
              case 'installed':
                // New SW Installed
                break;
              case 'redundant':
                // SW Became Redundant
                break;
              default:
                // Other State Change
                break;
            }
          };
        }
      };
    })
    .catch((e) => {
      Console.log('Service Worker Registration Failed');
      Console.log(e);
    });
}
