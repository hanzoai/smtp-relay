all: run

deps:
	go get -u github.com/julienschmidt/httprouter
	go get -u github.com/jordan-wright/email

run:
	sudo go run app.go


.PHONY: all deps run
