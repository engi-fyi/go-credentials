sonar:
	go test -coverprofile=reports/coverage.out ./tests
	scripts/sonar.sh

test:
	go test -v ./...

cover:
	go test -coverprofile=reports/coverage.out ./tests >/dev/null 2>&1
	go get -u github.com/mcubik/goverreport >/dev/null 2>&1
	~/go/bin/goverreport -coverprofile=reports/coverage.out