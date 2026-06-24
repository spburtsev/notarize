FROM node:24-alpine AS builder
RUN corepack enable
WORKDIR /app

COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml frontend/.npmrc ./
RUN pnpm install --frozen-lockfile

COPY openapi.yaml /openapi.yaml
COPY frontend/ ./
RUN pnpm gen:api && pnpm build

FROM node:24-alpine
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000
CMD ["node", "build"]
