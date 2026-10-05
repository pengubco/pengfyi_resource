# A "customized Python runtime": the base image every Python workload in the
# fleet runs on. Pre-installs the HTTP client and the in-house bootstrap that
# every workload would otherwise install on every start.
FROM public.ecr.aws/docker/library/python:3.12-alpine

RUN pip install --no-cache-dir "requests==2.32.3"

COPY bootstrap.py /opt/runtime/bootstrap.py
ENV PYTHONPATH=/opt/runtime

ENTRYPOINT ["python", "/opt/runtime/bootstrap.py"]
