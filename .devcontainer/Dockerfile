ARG VARIANT="buster"
FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

ENV TERRAFORM_VERSION 0.13.2
ENV YQ_VERSION=3.3.4
ENV HELM_VERSION=3.3.4
ENV AWS_AIM_AUTHENTICATOR_VERSION=1.17.9

# Install Terraform, yq, make, helm, aws-iam-authenticator
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && terraform --version \

    && wget https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 -O /usr/bin/yq \
    && chmod +x /usr/bin/yq \

    && apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends make \ 

    && wget https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz \
    && tar -zxvf helm-v${HELM_VERSION}-linux-amd64.tar.gz \
    && mv linux-amd64/helm /usr/local/bin/helm \
    && rm -rf linux-amd64 \
    && chmod +x /usr/local/bin/helm \
    && helm repo add stable https://kubernetes-charts.storage.googleapis.com \

    && wget https://amazon-eks.s3.us-west-2.amazonaws.com/${AWS_AIM_AUTHENTICATOR_VERSION}/2020-08-04/bin/linux/amd64/aws-iam-authenticator \
        -O /usr/local/bin/aws-iam-authenticator \
    && chmod +x /usr/local/bin/aws-iam-authenticator

