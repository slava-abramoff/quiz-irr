# ---------- build ----------
FROM node:20-alpine AS builder

WORKDIR /app

# принимаем build-time переменную от docker-compose
ARG VITE_API_BASE
ENV VITE_API_BASE=${VITE_API_BASE}

COPY package*.json ./
RUN npm install

COPY . .

# Vite увидит VITE_API_BASE здесь
RUN npm run build

# ---------- serve ----------
FROM nginx:1.25-alpine

# чистим дефолтный конфиг
RUN rm /etc/nginx/conf.d/default.conf

# свой nginx конфиг
COPY nginx.conf /etc/nginx/conf.d/default.conf

# собранная статика
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
