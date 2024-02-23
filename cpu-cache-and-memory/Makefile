build-cacheline:
	go build -o cacheline/bin/shared-write cacheline/shared-write/main.go
	go build -o cacheline/bin/independent-write cacheline/independent-write/main.go 
	go build -o cacheline/bin/shared-read cacheline/shared-read/main.go
	go build -o cacheline/bin/independent-read cacheline/independent-read/main.go 

run-cacheline-write:
	./run.sh cacheline/bin/shared-write 
	./run.sh cacheline/bin/independent-write

run-cacheline-read:
	./run.sh cacheline/bin/shared-read 
	./run.sh cacheline/bin/independent-read

build-write-combining:
	go build -o write-combining/bin/one-byte-one-loop write-combining/one-byte-one-loop/main.go
	go build -o write-combining/bin/one-byte-two-loops write-combining/one-byte-two-loops/main.go
	go build -o write-combining/bin/64-bytes-one-loop write-combining/64-bytes-one-loop/main.go
	go build -o write-combining/bin/64-bytes-two-loops write-combining/64-bytes-two-loops/main.go

run-write-combining-one-byte:
	./run.sh write-combining/bin/one-byte-one-loop
	./run.sh write-combining/bin/one-byte-two-loops

run-write-combining-64-bytes:
	./run.sh write-combining/bin/64-bytes-one-loop
	./run.sh write-combining/bin/64-bytes-two-loops
