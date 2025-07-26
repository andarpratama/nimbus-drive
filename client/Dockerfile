# Production Dockerfile for Vue.js (Vite) client using pm2
FROM node:20-slim AS build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Production image
FROM node:20-slim
WORKDIR /app

RUN npm install -g pm2
COPY --from=build /app/dist ./dist

EXPOSE 80
CMD ["pm2-runtime", "serve", "dist", "80", "--spa"] 