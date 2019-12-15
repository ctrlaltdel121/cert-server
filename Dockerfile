# from OSS golang img
FROM golang:latest AS build

# Make a dir for my code
WORKDIR /go/cert-srv

# add the code
ADD . .

# go build, will pull in deps using modules
RUN go build -o cert-srv cmd/cert-srv/main.go

# make a dir to store the certs in
RUN mkdir /certs

# set env that my process requires
ENV PORT=8888 STORAGE_DIR=/certs

# Set it to run when docker container boots, and expose the set port
ENTRYPOINT ["/go/cert-srv/cert-srv"]
EXPOSE 8888


# Note: I tried building in golang:latest and then copying the binary to 
# scratch to save space, but i was having issues testing scratch on windows.