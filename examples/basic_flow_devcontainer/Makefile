ROOT_DIR := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

#makes it easier to replace the value
M_NAME            ?= ks-basic-flow

M_SHARED_EXTERNAL ?= $(ROOT_DIR)/shared

VMS_RSA_FILENAME ?= vms_rsa

-include ./vars.mk

export

.PHONY: all apply destroy clean setup

all: apply

apply: init plan apply

destroy: plan-destroy-awsbi destroy-awsbi

clean:
	@rm -rf $(M_SHARED_EXTERNAL)

.PHONY: init plan apply output plan-destroy destroy

init: setup
	@cd $(M_WORKDIR) && $(MAKE) init \
		M_NAME=$(M_NAME)

plan: setup
	@cd $(M_WORKDIR) && $(MAKE) plan \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

apply: setup
	@cd $(M_WORKDIR) && $(MAKE) apply \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY \
		kubeconfig

output: setup
	@cd $(M_WORKDIR) && $(MAKE) output \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

plan-destroy: setup
	@cd $(M_WORKDIR) && $(MAKE) plan-destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

destroy: setup
	@cd $(M_WORKDIR) && $(MAKE) destroy \
		M_AWS_ACCESS_KEY=$$AWS_ACCESS_KEY \
		M_AWS_SECRET_KEY=$$AWS_SECRET_KEY

setup: $(M_SHARED)/$(VMS_RSA_FILENAME)

$(M_SHARED)/$(VMS_RSA_FILENAME):
	@mkdir -p $(dir $@)
	@ssh-keygen -t rsa -b 4096 -f $@ -N '' 2>&1 >/dev/null
