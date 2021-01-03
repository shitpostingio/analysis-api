# First stage: build the executable.
FROM golang:buster AS builder

# It is important that these ARG's are defined after the FROM statement
ARG SSH_PRIV="nothing"
ARG SSH_PUB="nothing"
ARG GOSUMDB=off

# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'analysis:x:65534:65534:analysis:/:' > /user/passwd && \
    echo 'analysis:x:65534:' > /user/group

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/shitpostingio/analysis-api

# Import the code from the context.
COPY .  .

# Build the executable
RUN go install

# Final stage: the running container.
FROM debian:latest

RUN apt update; \
    apt install -y ca-certificates curl

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Copy the built executable
COPY --from=builder /go/bin/analysis-api /home/analysis/server

# Give the right permissions
RUN chown -R analysis /home/analysis

# Expose the port
EXPOSE 9999

# Set the workdir
WORKDIR /home/analysis

# Perform any further action as an unprivileged user.
USER analysis:analysis

# Run the compiled binary.
CMD ["./server"]




