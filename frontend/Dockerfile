# Build stage
FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps
COPY . .
RUN npm run build

# Production stage: serve with nginx
FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
# Provide nginx config inline (fix try_files directive)
RUN cat <<EOF > /etc/nginx/nginx.conf
worker_processes 1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  65;
    server {
        listen 3000;
        server_name localhost;
        root /usr/share/nginx/html;
        index index.html;
        location / {
            try_files $uri $uri/ /index.html;
        }
    }
}
EOF
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
