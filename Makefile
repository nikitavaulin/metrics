ifneq ($(wildcard .env),)
    include .env
    export $(shell sed 's/=.*//' .env)
endif

build:
		@go build ./cmd/server && \
		go build ./cmd/agent

test:
		@go test -v ./... -count 1

run:
		@DATABASE_DSN="$(DATABASE_DSN)" go run ./cmd/server

migrate-up:
		@migrate -database "$(DATABASE_DSN)" -path ./migrations up 

migrate-down:
		@migrate -database "$(DATABASE_DSN)" -path ./migrations down 1

migrate-create: ## make migrate-create NAME=add_example
ifeq ($(strip $(NAME)),)
	$(error Укажите NAME: make migrate-create NAME=add_example)
endif
	migrate create -ext sql -dir ./migrations -seq $(NAME)

autotest-1:
		@./metricstest -test.v -test.run=^TestIteration1$ -binary-path=./server.exe

autotest-2a:
		@./metricstest -test.v -test.run=^TestIteration2A$ \
			-source-path=. \
			-agent-binary-path=./agent.exe

autotest-2b:
		@./metricstest -test.v -test.run=^TestIteration2B$ \
				-source-path=. \
				-agent-binary-path=./agent.exe

autotest-3a:
		@./metricstest -test.v -test.run=^TestIteration3A$ \
				-source-path=. \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe

autotest-3b:
		@./metricstest -test.v -test.run=^TestIteration3B$ \
				-source-path=. \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe

autotest-4:
		@./metricstest -test.v -test.run=^TestIteration4$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-server-port=8080 \
				-source-path=.

autotest-5:
		@ADDRESS="localhost:5050" \
		./metricstest -test.v -test.run=^TestIteration5$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-server-port=5050 \
				-source-path=.

autotest-6:
		@ADDRESS="localhost:5050" \
		./metricstest -test.v -test.run=^TestIteration6$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-server-port=5050 \
				-source-path=.

autotest-7:
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration7$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-server-port=8080 \
				-source-path=.

autotest-8:
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration8$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-server-port=8080 \
				-source-path=.

autotest-9:
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration9$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-file-storage-path=./metrics_temp.json \
				-server-port=8080 \
				-source-path=.

autotest-10a: build
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration10A$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-database-dsn='$(DATABASE_DSN)' \
				-server-port=8080 \
				-source-path=.

autotest-10b: build
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration10B$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-database-dsn='$(DATABASE_DSN)' \
				-server-port=8080 \
				-source-path=.

autotest-11: build
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration11$$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-database-dsn='$(DATABASE_DSN)' \
				-server-port=8080 \
				-source-path=.

.PHONY: autotest-12 autotest-13
autotest-12: build
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration12$$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-database-dsn='$(DATABASE_DSN)' \
				-server-port=8080 \
				-source-path=.

autotest-13: build
		@ADDRESS="localhost:8080" \
		./metricstest -test.v -test.run=^TestIteration13$$ \
				-agent-binary-path=./agent.exe \
				-binary-path=./server.exe \
				-database-dsn='$(DATABASE_DSN)' \
				-server-port=8080 \
				-source-path=.
