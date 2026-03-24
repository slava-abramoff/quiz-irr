# ---------- build ----------
FROM node:20-alpine AS builder

WORKDIR /app

# принимаем build-time переменную от docker-compose
ARG VITE_API_BASE_URL
ARG QUIZ_FRONTEND
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV QUIZ_FRONTEND=${QUIZ_FRONTEND}

COPY package*.json ./
RUN npm install

COPY . .

# Vite увидит VITE_API_BASE_URL здесь
RUN npm run build

# ---------- serve ----------
FROM nginx:1.25-alpine

# чистим дефолтный конфиг
RUN rm /etc/nginx/conf.d/default.conf

# свой nginx конфиг
COPY nginx.conf /etc/nginx/conf.d/default.conf

# собранная статика
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 81

CMD ["nginx", "-g", "daemon off;"]
