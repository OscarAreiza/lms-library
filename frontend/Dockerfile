# ---- Build stage ----
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# Baked in at build time — Vite inlines env vars into the static bundle.
# Defaults to going through nginx at a relative path once ADR-003's reverse proxy
# is set up (see infra/); override for this docker-compose's direct-port topology.
ARG VITE_API_BASE_URL=http://localhost:8080/api/v1
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

# ---- Run stage ----
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
