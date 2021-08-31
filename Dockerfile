FROM golang:1.5
EXPOSE 8080
WORKDIR /app
COPY . /app
RUN chmod a+x ./build ./run ./test
RUN ["./build"]
CMD ./run
