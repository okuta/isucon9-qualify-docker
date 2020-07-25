FROM ruby:2.7.1
WORKDIR /app

RUN apt-get update \
 && apt-get install -y \
    mariadb-client \
    less

COPY ./Gemfile ./
COPY ./Gemfile.lock ./
RUN bundle install

CMD /bin/sh -c "while sleep 1000; do :; done"
