# SvelteKit Frontend
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built application
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./

# Install production dependencies only
RUN npm install --omit=dev

# Create non-root user
RUN adduser -D appuser
USER appuser

EXPOSE 3000

ENV NODE_ENV=production
ENV PORT=3000

CMD ["node", "build"]
