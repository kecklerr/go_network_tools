FROM centos:7 
COPY template/ /template/
COPY go_network_tools /go_network_tools
EXPOSE 8080 
CMD /go_network_tools
