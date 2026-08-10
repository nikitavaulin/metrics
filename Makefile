build:
		@go build ./cmd/server && \
		go build ./cmd/agent

test:
		@go test -v ./... -count 1

run:
		@go run ./cmd/server

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