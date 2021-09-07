
.PHONY: mocks
gpg:
	test ! -f ~/.gnupg/pubring.kbx && gpg2 --export > ~/.gnupg/pubring.gpg 
mocks:
	go install github.com/golang/mock/mockgen@v1.5.0
	mockgen -source=internal/storage/storage.go -destination=internal/storage/storage_mock.go -package=storage
	mockgen -source=internal/location/location.go -destination=internal/location/location_mock.go -package=location

.PHONY: check
check:
	go test ./...
	goimports -w ./
