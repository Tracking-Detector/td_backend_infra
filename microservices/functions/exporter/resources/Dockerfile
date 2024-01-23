FROM node:20-alpine

# Create a non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser
USER appuser

WORKDIR /app

COPY exporter.js /app/

CMD ["node", "exporter.js"]