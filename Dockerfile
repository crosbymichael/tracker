FROM scratch

COPY bttracker/bttracker /usr/bin/

EXPOSE 80
ENTRYPOINT ["bttracker"]
CMD ["-redis-addr", "redis"]
