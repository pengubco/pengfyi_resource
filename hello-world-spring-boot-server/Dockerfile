FROM public.ecr.aws/amazonlinux/amazonlinux:2

# This is a hello-world demo. So just install the -devel.
RUN yum install -y java-17-amazon-corretto-devel

COPY ./target/HelloWorld-0.0.1-SNAPSHOT.jar webserver.jar
