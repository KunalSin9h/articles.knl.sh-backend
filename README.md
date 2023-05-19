# articles.knl.sh-backend

backend for https://articles.kunalsin9h.dev

# Docker running

```bash
sudo docker run --name articles-back -p 5000:5000 --rm -d \
    -e PASSWORD=1234 \
    -e DB=./db/prod.db \
    -e PORT=5000 \
    --mount type=bind,source="$(pwd)"/db/prod.db,target=/app/db/prod.db \
    ghcr.io/kunalsin9h/articles-back:latest
```
