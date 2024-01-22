#!/bin/bash
CONFIG_FOLDER=test-config


# Function to generate random string
generate_random_string() {
    docker run --rm authelia/authelia:latest authelia crypto rand --length 128 --charset alphanumeric | awk '{print $3}'
}

# Function to generate random password
generate_random_password() {
    local password=$(docker run authelia/authelia:latest authelia crypto rand --length 72 --charset rfc3986 | awk '{print $NF}') 
    echo "        # Password for digest: $password"
    local digest=$(docker run authelia/authelia:latest authelia crypto hash generate pbkdf2 --variant sha512 --password "$password" | awk '{print $NF}')
    echo "        secret: $digest"
}

# Function to generate Argon2 password hash
generate_password_hash() {
    local password=$1
    local digest=$(docker run authelia/authelia:latest authelia crypto hash generate argon2 --password "$password")
    echo "$digest"
}

# Function to generate RSA key pair and save private key to private.pem
generate_rsa_key_pair() {
    docker run -u "$(id -u):$(id -g)" -v "$(pwd)":/keys authelia/authelia:latest authelia crypto pair rsa generate --bits 4096 --directory /keys
}


# Generate .env file
generate_env_file() {
    local domain=$1
    echo "# mongo"
    echo "MONGO_URI=mongodb://db:27017/tracking-detector"
    echo "USER_COLLECTION=users"
    echo "REQUEST_COLLECTION=requests"
    echo "TRAINING_RUNS_COLLECTION=training-runs"
    echo "MODELS_COLLECTION=models"
    echo ""
    echo "# minio"
    echo "MINIO_URI=minio:9000"
    echo "MINIO_ACCESS_KEY=adminadmin"
    echo "MINIO_PRIVATE_KEY=password123"
    echo "EXPORT_BUCKET_NAME=exports"
    echo "MODEL_BUCKET_NAME=models"
    echo ""
    echo "# admin"
    echo "ADMIN_API_KEY=$(generate_random_string)"
    echo ""
    echo "DOMAIN=$domain"
    echo "CONFIG_FOLDER=$CONFIG_FOLDER"
}

# Generate Authelia config
generate_authelia_config() {
    local domain=$1
    echo "---"
    echo "###############################################################"
    echo "#                   Authelia configuration                    #"
    echo "###############################################################"
    echo ""
    echo "jwt_secret: $(generate_random_string)"
    echo "default_redirection_url: https://auth.$domain"
    echo ""
    echo "log:"
    echo "  level: debug"
    echo ""
    echo "totp:"
    echo "  issuer: auth.$domain"
    echo ""
    echo "authentication_backend:"
    echo "  file:"
    echo "    path: /config/users_database.yml"
    echo ""
    echo "access_control:"
    echo "  default_policy: deny"
    echo "  rules:"
    echo "    - domain: '$domain'"
    echo "      policy: bypass"
    echo "    - domain: '*.$domain'"
    echo "      policy: two_factor"
    echo ""
    echo "session:"
    echo "  name: authelia_session"
    echo "  domain: $domain"
    echo "  secret: $(generate_random_string)"
    echo "  expiration: 1h"
    echo "  inactivity: 5m"
    echo "  remember_me_duration: 1M"
    echo ""
    echo "regulation:"
    echo "  max_retries: 3"
    echo "  find_time: 120"
    echo "  ban_time: 300"
    echo ""
    echo "storage:"
    echo "  encryption_key: $(generate_random_string)"
    echo "  local:"
    echo "    path: /config/db.sqlite3"
    echo ""
    echo "notifier:"
    echo "  filesystem:"
    echo "    filename: /config/notification.txt"
    echo ""
    echo "identity_providers:"
    echo "  oidc:"
    echo "    hmac_secret: $(generate_random_string)"
    echo "    issuer_private_key: |"
    sed 's/^/      /' private.pem
    echo "     clients:"
    echo "      - id: portainer"
    echo "        description: Portainer"
    generate_random_password 
    echo "        public: false"
    echo "        authorization_policy: two_factor"
    echo "        redirect_uris:"
    echo "          - https://portainer.$domain"
    echo "        scopes:"
    echo "          - openid"
    echo "          - profile"
    echo "          - groups"
    echo "          - email"
    echo "        userinfo_signing_algorithm: none"
    echo "      - id: minio"
    echo "        description: MinIO"
    generate_random_password  
    echo "        public: false"
    echo "        authorization_policy: two_factor"
    echo "        redirect_uris:"
    echo "          - https://minio.$domain/apps/oidc_login/oidc"
    echo "        scopes:"
    echo "          - openid"
    echo "          - profile"
    echo "          - email"
    echo "          - groups"
    echo "        userinfo_signing_algorithm: none"
}

# Generate users_database.yml
generate_users_database() {
    echo "###############################################################"
    echo "#                         Users Database                      #"
    echo "###############################################################"
    echo ""
    echo "users:"
    while true; do
        read -p "Enter username (or 'exit' to finish): " username
        if [ "$username" == "exit" ]; then
            break
        fi
        read -p "Enter display name for $username: " displayname
        read -p "Enter email for $username: " email
        read -p "Enter password for $username: " password
        local password_hash=$(generate_password_hash "$password")
        echo "  $username:"
        echo "    disabled: false"
        echo "    displayname: \"$displayname\""
        echo "    password: $password_hash  # Password: $password"
        echo "    email: $email"
        echo "    groups:"
        echo "      - user"
    done
}

# Main function
main() {
    rm -f .test-env
    rm -rf ./$CONFIG_FOLDER
    read -p "Enter domain (e.g., tracking-detector.duckdns.org): " domain
    generate_env_file "$domain" > .test-env
    mkdir $CONFIG_FOLDER
    generate_rsa_key_pair > /dev/null
    generate_authelia_config "$domain" > ./$CONFIG_FOLDER/configuration.yml
    generate_users_database > ./$CONFIG_FOLDER/users_database.yml
    rm ./private.pem
    rm ./public.pem
    echo "Configuration files generated successfully."
}

# Run the wizard
main