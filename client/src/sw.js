/* eslint-disable no-restricted-globals */
/* eslint-disable no-underscore-dangle */
/* eslint-disable no-console */

const dataCacheName = 'webphoneData-v1.0.0001';
const cacheName = 'webphoneCache';

let filesToCache = [
  '/',
];

self.addEventListener('install', (e) => {
  self.skipWaiting().then(() => {
    // No Action
  });

  if (typeof self.__precacheManifest !== 'undefined') {
    const manifestCache = [];
    self.__precacheManifest.forEach((v) => {
      if (v.url !== '/index.gohtml') {
        manifestCache.push(v.url);
      }
    });
    filesToCache = filesToCache.concat(manifestCache);
  } else {
    // console.log('[ServiceWorker] Failed to load Manifest');
  }

  e.waitUntil(
    caches.open(cacheName).then((cache) => {
      fetch('/').then((response) => {
        cache.put('sw-offline-content', response.clone()).then(() => {
          // No Action
        });
      });
      return cache.addAll(filesToCache);
    }).catch((err) => {
      console.log('[ServiceWorker] Install Fetch Error');
      console.log(err);
    }),
  );
});

self.addEventListener('activate', (e) => {
  e.waitUntil(
    caches.keys().then((keyList) => {
      keyList.map((key) => (
        (key !== cacheName && key !== dataCacheName) ? caches.delete(key) : null
      ));
    }),
  );
  /*
   * Fixes a corner case in which the app wasn't returning the latest data.
   * You can reproduce the corner case by commenting out the line below and
   * then doing the following steps: 1) load app for first time so that the
   * initial New York City data is shown 2) press the refresh button on the
   * app 3) go offline 4) reload the app. You expect to see the newer NYC
   * data, but you actually see the initial data. This happens because the
   * service worker is not yet activated. The code below essentially lets
   * you activate the service worker faster.
   */
  return self.clients.claim();
});

self.addEventListener('fetch', (e) => {
  e.respondWith(
    caches.open(dataCacheName).then((cache) => fetch(e.request).then((response) => {
      cache.put(e.request.url, response.clone()).then(() => {
        // No Action
      }).catch(() => {
        // Don't Care
      });
      return response;
    })).catch((err) => {
      console.log('[ServiceWorker] Fetch - API Failure');
      console.log(err);
    }),
  );
});

self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  // eslint-disable-next-line no-undef
  event.waitUntil(clients.matchAll({
    type: 'window',
  }).then((clientList) => {
    clientList.forEach((client) => {
      if ('focus' in client) {
        client.postMessage({ event: event.action });
        client.focus();
      }
    });
  }));
});
