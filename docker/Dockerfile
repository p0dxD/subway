FROM golang:latest 
RUN mkdir -p /app/data
ADD subway /app/ 
COPY data/ /app/data
WORKDIR /app 
CMD ["/app/subway"]
