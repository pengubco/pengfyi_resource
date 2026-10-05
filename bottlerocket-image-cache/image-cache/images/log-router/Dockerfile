# A sidecar that every workload gets: reads the workload's log directory and
# forwards it. Here it just echoes, which is enough to show a second cached
# image starting next to the runtime.
FROM public.ecr.aws/docker/library/alpine:3.20

COPY log-router.sh /usr/local/bin/log-router
RUN chmod 0755 /usr/local/bin/log-router

ENTRYPOINT ["/usr/local/bin/log-router"]
