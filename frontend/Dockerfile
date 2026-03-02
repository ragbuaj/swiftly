# Build stage
FROM node:20-alpine AS build-stage

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install

# Copy the rest of the source code
COPY . .

# Build the application
RUN pnpm run build

# Production stage
FROM nginx:stable-alpine AS production-stage

# Copy the build output to nginx's public directory
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Expose the frontend port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
