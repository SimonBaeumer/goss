FROM centos:7

RUN echo "test content" > /tmp/test-file.txt
RUN chmod 0644 /tmp/test-file.txt
RUN ln -s /tmp/test-file.txt /tmp/symlink.txt
RUN mkdir /tmp/directory