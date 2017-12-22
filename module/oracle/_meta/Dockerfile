FROM wnameless/oracle-xe-11g

COPY healthcheck.sh /usr/sbin/healthcheck.sh

HEALTHCHECK --interval=30s --start-period=20s CMD [ "/usr/sbin/healthcheck.sh" ]
