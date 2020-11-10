FROM alpine:3

RUN apk update && apk add ca-certificates tzdata && rm -rf /var/cache/apk/*

COPY _output/whim_docker /app/whim_docker
COPY wait_for /app/wait-for.sh

WORKDIR /app

EXPOSE 7171

RUN chmod +x wait-for.sh

CMD ["./wait-for.sh" , "mysql:3306" , "--timeout=300" , "--" , "./whim_docker"]