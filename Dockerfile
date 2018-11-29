FROM golang:latest
MAINTAINER Dimitri Sokolyuk <demon@dim13.org>

RUN apt-get update && apt-get install -y \
	texlive-latex-extra \
	texlive-fonts-extra \
	texlive-lang-all \
	texlive-xetex \
	latexmk

COPY . src/github.com/dim13/gotex
RUN go install -v github.com/dim13/gotex

WORKDIR src/github.com/dim13/gotex
CMD exec gotex

EXPOSE 8080
