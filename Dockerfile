FROM golang

#RUN apt-get update && apt-get install -y texlive-full latexmk
RUN apt-get update
RUN apt-get install -y texlive-latex-extra
RUN apt-get install -y texlive-fonts-extra
RUN apt-get install -y texlive-lang-all
RUN apt-get install -y latexmk

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["go-wrapper", "run"]

COPY . /go/src/app
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 8080
