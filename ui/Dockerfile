# Use an official Node.js LTS image as the base image
FROM node:20-alpine as builder

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY ./ui/package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the application code
COPY ./ui /app/

# Build the Next.js application
RUN npm run build

# Use a lightweight Node.js image for the production environment
FROM node:20-alpine

# Set the working directory in the container
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json



# Start the Next.js application
CMD ["npm", "start"]
