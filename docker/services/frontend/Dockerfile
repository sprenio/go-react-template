
# Stage 1: build
FROM node:24.11 AS build

WORKDIR /app

ARG APP_NAME
RUN test -n "$APP_NAME" || (echo "APP_NAME build argument is required" && exit 1)

COPY ./frontend/package*.json ./
RUN npm ci

ARG GIT_COMMIT
ARG GIT_TAG

ENV VITE_APP_NAME=$APP_NAME
ENV VITE_GIT_COMMIT=$GIT_COMMIT
ENV VITE_GIT_TAG=$GIT_TAG

COPY ./frontend/ ./
RUN npm run build


# Stage 2: runtime
FROM nginx:1.29.0-alpine

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
