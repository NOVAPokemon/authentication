FROM brunoanjos/nova-server-base:latest

ENV executable="executable"
COPY $executable .

EXPOSE 8001

CMD ["sh", "-c", "./$executable -a -l"]