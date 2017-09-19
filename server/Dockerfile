FROM debian:stretch-slim

RUN mkdir -p /log

COPY rlog /

VOLUME ["/log"]

EXPOSE 1338

ENTRYPOINT ["/rlog", "-log", "/log/logregator.log", "-port", "1338", "-maddr", "224.0.0.1:9999"]
