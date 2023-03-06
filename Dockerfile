FROM debian:stable-slim
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y
RUN apt-get install apt-transport-https -y
RUN apt-get install apt-utils -y
RUN apt-get install gcc -y
RUN apt-get install g++ -y
RUN apt-get install nano -y
RUN apt-get install tar -y
RUN apt-get install bash -y
RUN apt-get install sudo -y
RUN apt-get install openssl -y
RUN apt-get install git -y
RUN apt-get install make -y
RUN apt-get install wget -y
RUN apt-get install curl -y
RUN apt-get install net-tools -y
RUN apt-get install iproute2 -y
RUN apt-get install ffmpeg -y
RUN apt-get install bc -y
RUN apt-get install imagemagick -y
RUN apt-get install libmagickwand-dev -y

ENV TZ="US/Eastern"
ARG USERNAME="morphs"
ARG PASSWORD="asdf"
RUN useradd -m $USERNAME -p $PASSWORD -s "/bin/bash"
RUN mkdir -p /home/$USERNAME
RUN chown -R $USERNAME:$USERNAME /home/$USERNAME
RUN usermod -aG sudo $USERNAME
RUN echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
USER $USERNAME
WORKDIR /home/$USERNAME

# install go with specific version and progress
COPY ./go_install.sh /home/$USERNAME/go_install.sh
RUN sudo chmod +x /home/$USERNAME/go_install.sh
RUN sudo chown $USERNAME:$USERNAME /home/$USERNAME/go_install.sh
RUN /home/$USERNAME/go_install.sh
RUN sudo tar --checkpoint=100 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo 
ZXhwb3J0IEdPX1RBUl9LSUxPQllURVM9JChwcmludGYgIiUuM2ZcbiIgJChlY2hvICIkKHN0YXQgLS1mb3JtYXQ9IiVzIiAvaG9tZS9tb3JwaHMvZ28udGFyLmd6KSAvIDEwMDAiIHwgYmMgLWwpKSAmJiBlY2hvIEV4dHJhY3RpbmcgWyRUQVJfQ0hFQ0tQT0lOVF0gb2YgJEdPX1RBUl9LSUxPQllURVMga2lsb2J5dGVzIC91c3IvbG9jYWwvZ28= 
| base64 -d ; echo); eval $cmd"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
RUN echo "PATH=$PATH:/usr/local/go/bin" | tee -a /home/$USERNAME/.bashrc

# Imagemagick Fix
RUN echo "PKG_CONFIG_PATH=/usr/local/lib/pkgconfig" | tee -a /home/$USERNAME/.bashrc
RUN echo "CGO_CFLAGS_ALLOW='-Xpreprocessor'" | tee -a /home/$USERNAME/.bashrc

ARG GO_ARCH=amd64
RUN git clone "https://github.com/0187773933/ImageUploadServerGo"
WORKDIR ImageUploadServerGo
RUN /usr/local/go/bin/go mod tidy
# force image-magick-v2
RUN /usr/local/go/bin/go get github.com/gographics/imagick@none
RUN /usr/local/go/bin/go get gopkg.in/gographics/imagick.v2/imagick
RUN PKG_CONFIG_PATH=/usr/local/lib/pkgconfig CGO_CFLAGS_ALLOW='-Xpreprocessor' GOOS=linux GOARCH=$GO_ARCH /usr/local/go/bin/go build -o /home/morphs/ImageUploadServer
#ENTRYPOINT [ "PKG_CONFIG_PATH=/usr/local/lib/pkgconfig CGO_CFLAGS_ALLOW='-Xpreprocessor' /home/morphs/ImageUploadServer" ]
ENTRYPOINT [ "/home/morphs/ImageUploadServer" ]
