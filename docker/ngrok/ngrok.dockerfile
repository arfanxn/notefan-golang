ARG NGROK_VERSION=alpine
FROM ngrok/ngrok:${NGROK_VERSION}

# Copy Ngrok configuration file to the container and set the Environment variable
COPY ./ngrok.yml /etc/ngrok.yml
ENV NGROK_CONFIG=/etc/ngrok.yml

# Expose port 4040
EXPOSE 4040

# TODO: automate this in the container
# ENTRYPOINT ["ngrok"]
# CMD ["http golang:8080"]