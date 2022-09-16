FROM registry.suse.com/bci/bci-micro:15.3
RUN echo 'pluginoperator:x:1000:1000::/home/pluginoperator:/bin/bash' >> /etc/passwd && \
    echo 'pluginoperator:x:1000:' >> /etc/group && \
    mkdir /home/pluginoperator && \
    chown -R pluginoperator:pluginoperator /home/pluginoperator
COPY bin/plugin-operator /usr/bin/
USER pluginoperator
CMD ["plugin-operator"]
