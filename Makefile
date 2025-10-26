DATE_TODAY := $(shell date +%Y-%m-%d)
DATE_7_DAYS_AGO := $(shell date -v -7d +%Y-%m-%d)
PROFILE ?= "profile"
DLDIR ?= "photos"

build:
	go build

updateDeps:
	go get -u .
	go mod tidy

fmt:
	go fmt .

login:
	go run . -v -profile $(PROFILE) -dldir $(DLDIR) -from $(DATE_TODAY)

test-year-month: killChrome
	go run . -v -profile $(PROFILE) -dldir $(DLDIR) -from $(DATE_7_DAYS_AGO) -headless -yearmonth

test: killChrome
	go run . -v -profile $(PROFILE) -dldir $(DLDIR) -from $(DATE_7_DAYS_AGO) -headless

help:
	go run . -h

killChrome:
	./killChrome.sh
