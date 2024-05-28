GOX := $(shell which go)
BIN := localdns

localdns:
	$(GOX) build \
		-v \
		-x \
		./cmd/localdns

clean:
	@rm localdns

.PHONY: clean
