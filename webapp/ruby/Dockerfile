FROM ruby:2.7.1
WORKDIR /app

COPY ./ ./
RUN bundle install

CMD /bin/sh -c "while sleep 1000; do :; done"
