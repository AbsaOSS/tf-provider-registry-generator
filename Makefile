
.PHONY: mocks
mocks:
	go install github.com/golang/mock/mockgen@v1.5.0
	mockgen -source=internal/terraform/file.go -destination=internal/terraform/file_mock.go -package=terraform

