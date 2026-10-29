# Placeholder until the stack is chosen.

FROM busybox:1.37
RUN mkdir -p /www \
    && echo "agneta placeholder" > /www/index.html \
    && adduser -D -u 10001 app
USER app
EXPOSE 8000
CMD ["httpd", "-f", "-v", "-p", "8000", "-h", "/www"]