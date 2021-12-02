FROM centos:7

WORKDIR "/usr/local"

RUN yum install -y yum-utils \
&& yum-config-manager --add-repo http://mirrors.aliyun.com/repo/Centos-7.repo \
&& yum-config-manager --add-repo http://mirrors.aliyun.com/repo/epel.repo \
&& yum-config-manager --add-repo http://mirrors.aliyun.com/repo/epel-7.repo \
&& yum repolist \
&& yum install -y iproute net-tools ca-certificates iptables strongswan tcpdump wget vim tcpdump less iproute-doc net-tools nmap-ncat.x86_64 lsof unzip openssh-server openssh-clients python2 python2-pip \
&& yum install -y golang

COPY boot.go /usr/local/


RUN cd /usr/local && ls

RUN cd /usr/local && go build boot.go

CMD /usr/local/boot

# docker build . -t storage
# docker run -d -p 9090:8080 -v /data/storage/:/data/ storage /usr/local/boot 8080 /data