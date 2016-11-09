# fritzctl

[![Build Status](https://travis-ci.org/bpicode/fritzctl.svg)](https://travis-ci.org/bpicode/fritzctl) [![CircleCI](https://circleci.com/gh/bpicode/fritzctl/tree/master.svg?style=shield)](https://circleci.com/gh/bpicode/fritzctl) [![build status](https://gitlab.com/bpicode/fritzctl/badges/master/build.svg)](https://gitlab.com/bpicode/fritzctl/commits/master)

[![Go Report Card](https://goreportcard.com/badge/github.com/bpicode/fritzctl)](https://goreportcard.com/report/github.com/bpicode/fritzctl) [![codecov](https://codecov.io/gh/bpicode/fritzctl/branch/master/graph/badge.svg)](https://codecov.io/gh/bpicode/fritzctl) [![Issue Count](https://codeclimate.com/github/bpicode/fritzctl/badges/issue_count.svg)](https://codeclimate.com/github/bpicode/fritzctl) [![Code Climate](https://codeclimate.com/github/bpicode/fritzctl/badges/gpa.svg)](https://codeclimate.com/github/bpicode/fritzctl)

[![Download](https://api.bintray.com/packages/bpicode/fritzctl_deb/fritzctl/images/download.svg)](https://bintray.com/bpicode/fritzctl_deb/fritzctl/_latestVersion)

[![Get automatic notifications about new fritzctl versions](https://www.bintray.com/docs/images/bintray_badge_color.png)](https://bintray.com/bpicode/fritzctl_deb/fritzctl?source=watch)

fritzctl is a command line client for the AVM FRITZ!Box primarily focused on the
[AVM Home Automation HTTP Interface](https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf).

## Install (debian/ubuntu/...)
Add the repository
```bash
echo "deb https://dl.bintray.com/bpicode/fritzctl_deb jessie main" | sudo tee -a /etc/apt/sources.list
```
and its signing key
```bash
wget -qO - https://api.bintray.com/users/bpicode/keys/gpg/public.key | sudo apt-key add -
```
The fingerprint of the repository key `3072D/35E71039` is `93AC 2A3D 418B 9C93 2986  6463 15FC CFC9 35E7 1039`.
Update your local repository data and install
```bash
sudo apt update
sudo apt install fritzctl
```

## Direct downloads

There are several locations from where one can download the packages, e.g.
*   the [most recent build from CI](https://gitlab.com/bpicode/fritzctl/builds/artifacts/master/download?job=build),
    where one also finds [older builds](https://gitlab.com/bpicode/fritzctl/pipelines),
*   directly from the [debian repository](https://dl.bintray.com/bpicode/fritzctl_deb/fritzctl)
*   directly from the [rpm repository](https://bintray.com/bpicode/fritzctl_rpm/fritzctl)


## Usage

![Demo usage](/images/fritzctl_demo.gif?raw=true "Demo usage")
