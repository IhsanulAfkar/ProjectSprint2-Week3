FROM golang:1.22


WORKDIR /Week3

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /halo_sus
EXPOSE 8080
CMD ["/halo_sus"]