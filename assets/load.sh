set -e


mkdir -p "assets/vendor/fonts"
curl -o "assets/vendor/bootstrap.min.css" https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/css/bootstrap.min.css
curl -o "assets/vendor/bootstrap.bundle.min.js" https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/js/bootstrap.bundle.min.js
curl -o "assets/vendor/bootstrap-icons.min.css" https://cdn.jsdelivr.net/npm/bootstrap-icons@1.13.1/font/bootstrap-icons.min.css
curl -o "assets/vendor/fonts/bootstrap-icons.woff2" "https://cdn.jsdelivr.net/npm/bootstrap-icons@1.13.1/font/fonts/bootstrap-icons.woff2?e34853135f9e39acf64315236852cd5a="
curl -o "assets/vendor/fuse.js" https://cdn.jsdelivr.net/npm/fuse.js@7.0.0
curl -o "assets/vendor/htmx.min.js" https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js
