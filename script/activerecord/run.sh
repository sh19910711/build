#!/bin/sh

cd $(dirname "${BASH_SOURCE[0]}")
bundle install -j4
bundle exec ruby seed.rb

cd app
docker build -t build .
