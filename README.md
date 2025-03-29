Shorten a URL (POST request to /shorten)

1. curl -X POST -d "url=https://example.com" http://<host>:<port>/shorten
    Response: Short URL: http://<host>:<port>/abc123

2. Visit the short URL
  Open http://<host>:<port>/abc123 → Redirects to https://example.com
