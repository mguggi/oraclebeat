FROM wnameless/oracle-xe-11g

# overriding healtcheck
HEALTHCHECK --interval=5s --retries=60 --start-period=30s CMD test `netstat -tulpen | grep 0.0.0.0:1521 | wc -l` -eq 1
