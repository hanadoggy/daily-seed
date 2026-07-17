FROM node:24-alpine

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

WORKDIR /app

# Install dependencies first
COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

# Copy application code
COPY . .

# Expose Vite's default port
EXPOSE 5173

# Start development server
CMD ["pnpm", "run", "dev"]
