FROM debian:stable-slim
COPY digitalshelf /bin/digitalshelf

CMD ["/bin/digitalshelf"]