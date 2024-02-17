FROM golang:1.20-alpine

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./
COPY go.sum ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY *.go ./
COPY ./tools ./tools
COPY ./database ./database

# compile application
RUN go build -o /godocker

# command to be used to execute when the image is used to start a container
CMD [ "/godocker" ]