docker run -e SONAR_HOST_URL=https://sonarcloud.io/ -e SONAR_TOKEN=$SONAR_TOKEN --user="0:0" -it -v "$(pwd):/usr/src" sonarsource/sonar-scanner-cli /usr/bin/entrypoint.sh -X