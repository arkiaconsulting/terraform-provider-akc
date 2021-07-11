# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.187.0/containers/ubuntu/.devcontainer/base.Dockerfile

FROM mcr.microsoft.com/vscode/devcontainers/base:buster

COPY library-scripts/*.sh /tmp/library-scripts/

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# Install needed packages and setup non-root user. Use a separate RUN statement to add your own dependencies.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG INSTALL_ZSH="true"
ARG INSTALL_AZURE_CLI="true"
ARG UPGRADE_PACKAGES="true"
ARG GO_VERSION="1.16.5"

RUN wget -c https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" "true" "true" \
    && apt-get -y install build-essential \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

RUN bash /tmp/library-scripts/azcli-debian.sh \
    && echo '#!/bin/bash\n"$@"' > /usr/local/share/docker-init.sh && chmod +x /usr/local/share/docker-init.sh \
    && rm -rf /var/lib/apt/lists/*

ARG TERRAFORM_VERSION=1.0.2
ARG TFLINT_VERSION=0.8.2

RUN bash /tmp/library-scripts/terraform-debian.sh "${TERRAFORM_VERSION}" "${TFLINT_VERSION}" \
    && rm -rf /tmp/library-scripts

ADD https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash /home/vscode/.git-completion.bash
RUN chown vscode /home/vscode/.git-completion.bash \
    && echo 'if [ -f ~/.git-completion.bash ]; then\n' \
    '. ~/.git-completion.bash\n' \
    'fi\n' \
    >> /home/vscode/.bashrc

RUN echo '. ~/.ipc' >> home/vscode/.bashrc

ENTRYPOINT [ "/usr/local/share/docker-init.sh" ]
CMD [ "sleep", "infinity" ]