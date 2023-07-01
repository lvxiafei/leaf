

build:
	go generate .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build

transfer:
	curl --upload-file  ./leaf https://transfer.sh/leaf  -w '\n\n'

bashupload:
	curl bashupload.com -T ./leaf

download:
	wget -O leaf ${url}

upload-json:
	curl --upload-file  ./output.json https://transfer.sh/output.json  -w '\n\n'

copy:
	curl https://transfer.sh/h3KMD3cy9v/output.json | pbcopy
	cat your_file.txt | yank

download-copy:
	curl https://transfer.sh/h3KMD3cy9v/output.json | tee output.json | pbcopy

docker-run:
	docker run --privileged --rm -t --pid=host -v /sys/kernel/debug/:/sys/kernel/debug/ cilium/pwru pwru --output-tuple 'host 1.1.1.1'
	docker run --privileged --rm -t --pid=host -v /sys/kernel/debug/:/sys/kernel/debug/ xx/leaf leaf -H 'host 1.1.1.1'

nerdctl-run:
	nerdctl  run --privileged --rm -it --pid=host -v /sys/kernel/debug/:/sys/kernel/debug/ cilium/pwru pwru ---output-tuple 'host 1.1.1.1'

docker-interactive:
	docker run --privileged --rm -it --pid=host -v /sys/kernel/debug/:/sys/kernel/debug/ cilium/pwru