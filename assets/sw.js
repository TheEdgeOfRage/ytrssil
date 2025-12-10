const CACHE_NAME = 'ytrssil-v1';
const ASSETS_TO_CACHE = [
	'/',
	'/assets/vendor/bootstrap.min.css',
	'/assets/vendor/bootstrap-icons.min.css',
	'/assets/vendor/bootstrap.bundle.min.js',
	'/assets/vendor/htmx.min.js',
	'/assets/vendor/fuse.js',
	'/assets/index.js',
	'/assets/ytrssil.png',
	'/assets/manifest.json'
];

self.addEventListener('install', (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then((cache) => {
			return cache.addAll(ASSETS_TO_CACHE);
		})
	);
	self.skipWaiting();
});

self.addEventListener('activate', (event) => {
	event.waitUntil(
		caches.keys().then((cacheNames) => {
			return Promise.all(
				cacheNames.map((cacheName) => {
					if (cacheName !== CACHE_NAME) {
						return caches.delete(cacheName);
					}
				})
			);
		})
	);
	self.clients.claim();
});

self.addEventListener('fetch', (event) => {
	event.respondWith(
		caches.match(event.request).then((response) => {
			return response || fetch(event.request);
		})
	);
});
