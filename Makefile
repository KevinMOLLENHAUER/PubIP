build:
	@go build -o bin/main main.go

run:
	@go run main.go

containerize:
	@podman build -t kmollenhauer/pinger .
	
run-container: containerize
	@podman run --rm -p 9090:9090 -it --name pinger kmollenhauer/pinger
