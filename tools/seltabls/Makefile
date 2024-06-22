.PHONY: all
all: install

.PHONY: act
act:
	# run github actions locally
	@act

.PHONY: install
install:
	@sh ./scripts/makefile/install.sh

.PHONY: test
test:
	@sh ./scripts/makefile/test.sh

.PHONY: coverage
coverage:
	@sh ./scripts/makefile/coverage.sh

.PHONY: lint
lint:
	@sh ./scripts/makefile/lint.sh

.PHONY: clean
clean:
	@sh ./scripts/makefile/clean.sh

.PHONY: tidy
tidy:
	@sh ./scripts/makefile/tidy.sh

.PHONY: database
database:
	@sh ./scripts/makefile/database.sh

.PHONY: docs
docs:
	@sh ./scripts/makefile/docs.sh

.PHONY: fmt
fmt:
	@sh ./scripts/makefile/fmt.sh
