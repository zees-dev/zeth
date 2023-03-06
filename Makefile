all: clean test build run

clean:
	rm -rf zeth
	rm -rf Zeth/
	rm -rf ./**/.ethereum
	# cargo clean
.PHONY: clean

test:
	# TODO
.PHONY: test

db:
	docker run --rm -it --name surrealdb \
		-v ${PWD}/Zeth:/Zeth/ \
		-p 8000:8000 \
		surrealdb/surrealdb:latest \
		start --log trace --user admin --pass admin file://Zeth/zeth.db
.PHONY: client

client:
	cd client && yarn dev
.PHONY: client

server:
	cargo run
.PHONY: server

build:
	cargo build
.PHONY: build

run:
	cargo run
.PHONY: run

data:
	# TODO
	# ./scripts/register_nodes.sh
.PHONY: data


# run:
# 	cargo run

# watch:
# 	cargo watch --quiet --clear --exec 'run --quiet'

# build:
# 	cargo build

# build-release:
# 	cargo build --release

# test:
# 	cargo test

# docker-build:
# 	docker build -t zeeshans/slim:prayer-alarm-rust .

# docker-run:
# 	docker run --rm -it -e PULSE_SERVER=host.docker.internal -p 3000:3000 --mount type=bind,source=${HOME}/.config/pulse,target=/home/pulseaudio/.config/pulse zeeshans/prayer-alarm-rust

# docker-build-rpi:
# 	docker build -t zeeshans/slim:prayer-alarm-rust -f rpi.Dockerfile .

# docker-run-rpi:
# 	docker run --rm -it --device /dev/snd --security-opt seccomp=unconfined -p 3000:3000 zeeshans/slim:prayer-alarm-rust