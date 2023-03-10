FROM debian:stable-20230208-slim

# create service dir
RUN mkdir -p /service/log

# set user env for glog
ENV USER=root

# copy service files
WORKDIR /service
COPY .env .
COPY main .

ENTRYPOINT ["/service/main", "-log_dir=./log"]