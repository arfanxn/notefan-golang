ARG NGINX_VERSION=stable
FROM nginx:${NGINX_VERSION}

# Copy Nginx configuration file to the container
COPY nginx.conf /etc/nginx/conf.d

# Expose port 80
EXPOSE 80