FROM node:12.22-alpine3.15

ENV EDITOR_PORT 8080
ENV SWAGGER_FILE /spec/swagger.yaml

RUN npm install -g swagger-editor-live@2.1.8

CMD swagger-editor-live $SWAGGER_FILE --host=0.0.0.0 --port=$EDITOR_PORT
