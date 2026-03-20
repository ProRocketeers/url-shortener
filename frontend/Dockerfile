# 🛠️ Použij oficiální Node.js image jako základ
FROM node:lts-alpine AS builder

WORKDIR /app

COPY package.json pnpm-lock.yaml ./

RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile

COPY . .

ARG NEXT_PUBLIC_BASE_PATH=/
ARG NEXT_PUBLIC_API_URL=https://url.rock8cloud.cz
ENV NEXT_PUBLIC_BASE_PATH=${NEXT_PUBLIC_BASE_PATH}
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}

RUN pnpm run build

# 🚀 Použij lehčí image pro běh aplikace
FROM node:lts-alpine AS runner

ARG NODE_ENV=production
ARG PORT=3000

WORKDIR /app

COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

RUN addgroup -S app && adduser -S app -G app \
	&& mkdir -p /app/temp \
	&& chown -R app:app /app/temp
USER app

ENV NODE_ENV=${NODE_ENV}
ENV PORT=${PORT}

EXPOSE ${PORT}

CMD ["node", "server.js"]
