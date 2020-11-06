ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

AWSKS := epiphanyplatform/awsks:0.0.1

#makes it easier to replace the value
M_NAME            ?= awsks-example-tests
M_SHARED_EXTERNAL ?= $(ROOT_DIR)/shared

VMS_RSA_FILENAME ?= vms_rsa

-include ./vars.mk

export

.PHONY: all apply destroy clean setup

all: apply

apply: init-awsks plan-awsks apply-awsks

destroy: plan-destroy-awsks destroy-awsks

clean:
	@rm -rf $(M_SHARED_EXTERNAL)

.PHONY: init-awsks plan-awsks apply-awsks output-awsks plan-destroy-awsks destroy-awsks

init-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		init \
		M_NAME=$(M_NAME) \
		M_REGION=$(M_REGION) \
		M_VPC_ID=$(M_VPC_ID) \
		M_SUBNET_IDS=$(M_SUBNET_IDS)

plan-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		plan \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

apply-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		apply \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY \
		kubeconfig

output-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		output \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

plan-destroy-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		plan-destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

destroy-awsks: setup
	@docker run --rm \
		-v $(M_SHARED_EXTERNAL):/shared \
		-t $(AWSKS) \
		destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

setup: $(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME)

$(M_SHARED_EXTERNAL)/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
