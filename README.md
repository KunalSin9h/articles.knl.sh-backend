# aniketsingh.co.in
backend for aniketsingh.co.in

# Docker running


```bash
sudo docker run --name asback -p 5005:5005 --rm -d \
    -e PASSWORD=1234 \
    -e DB=./db/prod.db \
    -e PORT=5005 \
    --mount type=bind,source="$(pwd)"/db/prod.db,target=/app/db/prod.db \
    kunalsin9h/asback:latest
```