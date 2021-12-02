FROM centos:7

WORKDIR "/usr/local"

RUN yum install -y yum-utils \
&& yum-config-manager --add-repo http://mirrors.aliyun.com/repo/Centos-7.repo \
&& yum-config-manager --add-repo http://mirrors.aliyun.com/repo/epel.repo \
&& yum-config-manager --add-repo http://mirrors.aliyun.com/repo/epel-7.repo \
&& yum repolist \
&& yum install -y golang \
&& go build boot.go

CMD boot