default: install

generate:
	go generate ./...

install:
	go install .

test:
	go test -count=1 -parallel=4 ./...

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

templates-get:
	bash ./elestio-templates/scripts/get_templates_list.sh

templates-generate:
	bash ./elestio-templates/scripts/generate_templates_examples.sh

templates-delete:
	bash ./elestio-templates/scripts/delete_generated_examples.sh