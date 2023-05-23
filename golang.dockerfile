ARG GOLANG_VERSION=alpine
FROM golang:${GOLANG_VERSION}

# Working directory
WORKDIR /app

# Copy everything at /app
COPY . /app

# Download packages and build the go app
RUN go get
RUN go build -o main .

# Expose port
EXPOSE 8080

# Execute the generated executable binary file from through command 
CMD ["./main"]
