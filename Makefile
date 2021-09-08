
.PHONY: mocks
gpg:
	test ! -f ~/.gnupg/pubring.kbx && gpg2 --export > ~/.gnupg/pubring.gpg 
mocks:
	go install github.com/golang/mock/mockgen@v1.5.0
	mockgen -source=internal/storage/storage.go -destination=internal/storage/storage_mock.go -package=storage
	mockgen -source=internal/location/location.go -destination=internal/location/location_mock.go -package=location
	mockgen -source=internal/terraform/terraform.go -destination=internal/terraform/terraform_mock.go -package=terraform
	mockgen -source=internal/repo/repo.go -destination=internal/repo/repo_mock.go -package=repo

.PHONY: check
check:
	goimports -w ./
	go test ./...
