.PHONY: tf-lint
tf-lint: tf-init
	cd infrastructure/aws && ./lint

.PHONY: tf-init 
tf-init: 
	cd infrastructure/aws && terraform init