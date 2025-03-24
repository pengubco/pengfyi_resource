FROM public.ecr.aws/amazonlinux/amazonlinux:2

# This image is used to reproduce page thrashing on EC2.

RUN amazon-linux-extras install epel -y && \ 
  yum install stress -y 

COPY run-parallel.sh .
RUN chmod +x ./run-parallel.sh
