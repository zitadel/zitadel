# Hướng dẫn triển khai Zitadel production (proxy domain, DB container riêng)

## 1. Yêu cầu
- Docker Compose v2 trở lên
- Đã có sẵn một container Postgres tên `DB` (user: `postgres`, pass: `postgres`)
- Đã tạo user `zitadel` (pass: `zitadel`) và database `zitadel` hoàn toàn mới, hoặc để Zitadel tự tạo
- Cả 2 container (`zitadel` và `DB`) cùng nằm trong 1 Docker network (ví dụ: `Gateway`)
- Đã có proxy (Traefik/Cloudflare) trỏ domain public về service Zitadel (ví dụ: `auth.joy.box`)

## 2. Cấu hình file

### `.env` mẫu (chuẩn production)
```env
ZITADEL_MASTERKEY=MasterkeyNeedsToHave32Characters

# Domain public truy cập qua proxy
ZITADEL_EXTERNALDOMAIN=auth.joy.box
ZITADEL_EXTERNALPORT=443
ZITADEL_TLS_ENABLED=false
ZITADEL_EXTERNALSECURE=true
ZITADEL_TLSMODE=external

# Thông tin database Postgres
ZITADEL_DATABASE_POSTGRES_HOST=DB
ZITADEL_DATABASE_POSTGRES_PORT=5432
ZITADEL_DATABASE_POSTGRES_DATABASE=zitadel
ZITADEL_DATABASE_POSTGRES_USER_USERNAME=zitadel
ZITADEL_DATABASE_POSTGRES_USER_PASSWORD=zitadel
ZITADEL_DATABASE_POSTGRES_USER_SSL_MODE=disable
ZITADEL_DATABASE_POSTGRES_ADMIN_USERNAME=postgres
ZITADEL_DATABASE_POSTGRES_ADMIN_PASSWORD=postgres
ZITADEL_DATABASE_POSTGRES_ADMIN_SSL_MODE=disable

# Logging
ZITADEL_LOGLEVEL=info
```

### `docker-compose.yaml` mẫu
```yaml
services:
  zitadel:
    restart: 'always'
    image: 'ghcr.io/zitadel/zitadel:latest'
    container_name: zitadel
    command: 'start-from-init --steps /zitadel-init-steps.yaml --masterkey "${ZITADEL_MASTERKEY}"'
    env_file:
      - .env
    depends_on: []
    ports:
      - '8080:8080'
    volumes:
      - ./zitadel-init-steps.yaml:/zitadel-init-steps.yaml:ro
    networks:
      - 'Gateway'
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.zitadel.rule=Host(`auth.joy.box`)"
      - "traefik.http.routers.zitadel.entrypoints=websecure"
      - "traefik.http.routers.zitadel.tls.certresolver=myresolver"
      - "traefik.http.services.zitadel.loadbalancer.server.port=8080"
      - "traefik.http.middlewares.csp.headers.contentSecurityPolicy=media-src 'none';script-src 'self' 'unsafe-eval';style-src 'self' 'unsafe-inline';img-src 'self' auth.joy.box blob: data:;frame-src 'none';frame-ancestors 'none';font-src 'self';manifest-src 'self';connect-src 'self' auth.joy.box;default-src 'none';object-src 'none'"
      - "traefik.http.routers.zitadel.middlewares=csp@docker"
networks:
  Gateway:
    external: true
```

## 3. Reset lại database khi cần khởi tạo lại từ đầu

1. Tắt container Zitadel:
   ```powershell
   docker stop zitadel
   ```
2. Xóa database `zitadel` trong container DB:
   ```powershell
   docker exec DB psql -U postgres -c "DROP DATABASE IF EXISTS zitadel;"
   docker exec DB psql -U postgres -c "CREATE DATABASE zitadel;"
   ```
3. Khởi động lại Zitadel:
   ```powershell
   docker start zitadel
   ```

## 4. Đăng nhập admin lần đầu
- Tài khoản admin sẽ lấy theo file `zitadel-init-steps.yaml` (ví dụ: email `example@zitadel.com`, mật khẩu `Password1!`)
- Đăng nhập tại: `https://auth.joy.box/ui/console/`

## 5. Tham khảo thêm
- [Zitadel Docs](https://docs.zitadel.com/)
- [Cấu hình steps.yaml](https://github.com/zitadel/zitadel/blob/main/cmd/setup/steps.yaml)
    env_file:
      - .env
    depends_on: []
    ports:
      - '8080:8080'

networks:
  zitadel:
```

## 3. Các bước triển khai

1. **Đảm bảo container DB đã chạy và có user/database đúng như trên.**
2. **Cập nhật file `.env` và `docker-compose.yaml` như mẫu trên.**
3. **Khởi động Zitadel:**
   ```powershell
   docker compose up -d
   ```
4. **Truy cập giao diện quản trị:**
   - http://localhost:8080/ui/console

## 4. Lưu ý
- Nếu từng chạy service `db` cũ trong compose, nên dọn dẹp orphan container:
  ```powershell
  docker compose up -d --remove-orphans
  ```
  hoặc
  ```powershell
  docker rm -f zitadel-db-1
  ```
- Nếu muốn tích hợp domain, TLS, hoặc DB ngoài host, chỉ cần sửa lại các biến trong `.env`.
- Nếu gặp lỗi migration, kiểm tra lại quyền user, tên host, và network giữa các container.

---
**Mọi thắc mắc hoặc cần tùy biến sâu hơn, liên hệ kỹ thuật để được hỗ trợ!**
