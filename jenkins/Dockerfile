FROM jenkins/jenkins:lts

USER root

# Ensure the jenkins user has sudo privileges
RUN echo "jenkins ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Switch to Jenkins user
USER jenkins


# Install Docker CLI and dependencies
RUN apt-get -y update && apt-get install -y sudo git wget curl apt-transport-https ca-certificates gnupg2 software-properties-common

# Install Docker using the official installation script
RUN curl -fsSL https://get.docker.com | sh

# Install Docker Compose
RUN curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
    chmod +x /usr/local/bin/docker-compose

# Set up Docker group
# RUN usermod -aG docker jenkins

# Ensure Docker socket permissions
RUN chmod 666 /var/run/docker.sock || true

# Enable Docker-in-Docker
ENV DOCKER_TLS_CERTDIR=/certs
RUN mkdir /certs /certs/client && chmod 1777 /certs /certs/client
ENV DOCKER_DRIVER=overlay2


# Set the working directory
WORKDIR /var/jenkins_home

# Start Docker daemon in the background
# ENTRYPOINT ["sh", "-c", "dockerd-entrypoint.sh & jenkins"]