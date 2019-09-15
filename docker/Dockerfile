FROM golang:latest 
RUN mkdir -p /app/data /app/static
ADD subway /app/ 
COPY data/ /app/data
COPY static/ /app/static
WORKDIR /app 
CMD ["/app/subway"]
