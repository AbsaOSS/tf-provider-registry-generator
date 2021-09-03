
.PHONY: mocks
mocks:
	go install github.com/golang/mock/mockgen@v1.5.0
	mockgen -source=internal/storage/storage.go -destination=internal/storage/storage_mock.go -package=storage

