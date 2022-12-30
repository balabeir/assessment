test:
	go clean -testcache && go test -v --tags=unit ./...

test-it:
	docker-compose -f docker-compose.test.yml down
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests