# Tracking Detector Backend Infrastructure

# Table of Contents

# Getting Started
This project aims to offer a comprehensive backend solution tailored for researchers in the privacy and web tracking domains. It provides a fully functional infrastructure to facilitate the collection of datasets, the training of TensorFlow-based models, and the comparative analysis of these models. Additionally, the backend empowers users to effortlessly create custom data exports. The overarching objective is to enable the development and validation of machine learning models proficient in detecting web trackers.

While we will host this service, the setup process is designed to be exceptionally straightforward. The project is open source, allowing users the freedom to deploy it on their servers. Feel free to utilize our backend if you choose to embark on your deployment journey.

# Prerequisites
To initiate the deployment of our service, it's essential to be the owner of a public domain. This requirement is crucial for the creation of valid SSL certificates.

# Deployment
```sh
# Clone Repository
git clone git@github.com:Tracking-Detector/td_backend_infra.git
# Give permission to execute setup wizard
chmod +x wizard.sh
# Run setup wizard to generate users passwords and env variables
./wizard.sh
# Build Docker images
docker-compose build
# Start application
docker-compose up -d
```


## Authors
@HenrySchwerdt @philip-raschke