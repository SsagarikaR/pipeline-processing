# ---- Build stage ----
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

ARG VITE_API_BASE_URL=http://localhost:8080/api/v1
ARG VITE_API_KEY=secret-pipeline-key
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
ENV VITE_API_KEY=$VITE_API_KEY

RUN npm run build

# ---- Serve stage ----
FROM node:20-alpine

RUN npm install -g serve && adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app/dist ./dist

USER appuser

# Ports below 1024 need root to bind, and appuser deliberately isn't
# root, so serve listens on 3000 instead of 80 - the host-side port
# (UI_PORT in docker-compose) is unaffected either way.
EXPOSE 3000

ENTRYPOINT ["serve", "-s", "dist", "-l", "3000"]
