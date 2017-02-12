# fritzctl

[![Build Status TravisCI](https://travis-ci.org/bpicode/fritzctl.svg)](https://travis-ci.org/bpicode/fritzctl)
[![Build Status CircleCI](https://circleci.com/gh/bpicode/fritzctl/tree/master.svg?style=shield)](https://circleci.com/gh/bpicode/fritzctl)
[![Build Status GitlabCI](https://gitlab.com/bpicode/fritzctl/badges/master/build.svg)](https://gitlab.com/bpicode/fritzctl/commits/master)
[![Build Status SemaphoreCI](https://semaphoreci.com/api/v1/bpicode/fritzctl/branches/master/shields_badge.svg)](https://semaphoreci.com/bpicode/fritzctl)
[![Build status AppVeyor](https://ci.appveyor.com/api/projects/status/k7qqx91w6mja3u7h?svg=true)](https://ci.appveyor.com/project/bpicode/fritzctl)

[![License](https://img.shields.io/github/license/bpicode/fritzctl.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/bpicode/fritzctl?status.svg)](https://godoc.org/github.com/bpicode/fritzctl)
[![Repository size](https://reposs.herokuapp.com/?path=bpicode/fritzctl)](https://github.com/bpicode/fritzctl)

[![Go Report Card](https://goreportcard.com/badge/github.com/bpicode/fritzctl)](https://goreportcard.com/report/github.com/bpicode/fritzctl)
[![codecov](https://codecov.io/gh/bpicode/fritzctl/branch/master/graph/badge.svg)](https://codecov.io/gh/bpicode/fritzctl)
[![Issue Count](https://codeclimate.com/github/bpicode/fritzctl/badges/issue_count.svg)](https://codeclimate.com/github/bpicode/fritzctl)
[![Code Climate](https://codeclimate.com/github/bpicode/fritzctl/badges/gpa.svg)](https://codeclimate.com/github/bpicode/fritzctl)
[![codebeat badge](https://codebeat.co/badges/605cf539-21dd-4a60-a892-e0d6da3021fe)](https://codebeat.co/projects/github-com-bpicode-fritzctl)

deb [![Download .deb](https://api.bintray.com/packages/bpicode/fritzctl_deb/fritzctl/images/download.svg)](https://bintray.com/bpicode/fritzctl_deb/fritzctl/_latestVersion)

rpm [![Download .rpm](https://api.bintray.com/packages/bpicode/fritzctl_rpm/fritzctl/images/download.svg)](https://bintray.com/bpicode/fritzctl_rpm/fritzctl/_latestVersion)

win [![Download](https://api.bintray.com/packages/bpicode/fritzctl_win/fritzctl/images/download.svg)](https://bintray.com/bpicode/fritzctl_win/fritzctl/_latestVersion)

[![Get automatic notifications about new fritzctl versions (deb)](https://www.bintray.com/docs/images/bintray_badge_color.png)](https://bintray.com/bpicode/fritzctl_deb/fritzctl?source=watch)
[![Get automatic notifications about new fritzctl versions (rpm)](https://www.bintray.com/docs/images/bintray_badge_bw.png)](https://bintray.com/bpicode/fritzctl_rpm/fritzctl?source=watch)
[![Get automatic notifications about new fritzctl versions (rpm)](https://www.bintray.com/docs/images/bintray_badge_greyscale.png)](https://bintray.com/bpicode/fritzctl_win/fritzctl?source=watch)

fritzctl is a command line client for the AVM FRITZ!Box primarily focused on the
[AVM Home Automation HTTP Interface](https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf).

Software is tested with

*   FRITZ!Box Fon WLAN 7390 running FRITZ!OS 06.51

## Install (debian/ubuntu/...)

Add the repository

```bash
echo "deb https://dl.bintray.com/bpicode/fritzctl_deb jessie main" | sudo tee -a /etc/apt/sources.list
```

and its signing key

```bash
wget -qO - https://api.bintray.com/users/bpicode/keys/gpg/public.key | sudo apt-key add -
```

The fingerprint of the repository key `3072D/35E71039` is
`93AC 2A3D 418B 9C93 2986  6463 15FC CFC9 35E7 1039`.
Update your local repository data and install

```bash
sudo apt update
sudo apt install fritzctl
```

## Install (opensuse)

Add the repository

```bash
wget https://bintray.com/bpicode/fritzctl_rpm/rpm -O bintray-bpicode-fritzctl_rpm.repo && sudo zypper ar -f bintray-bpicode-fritzctl_rpm.repo && rm bintray-bpicode-fritzctl_rpm.repo
```

Update your local repository data and install

```bash
sudo zypper refresh
sudo zypper in fritzctl
```

## Install (windows)

Windows binaries can found in the [windows directory](https://dl.bintray.com/bpicode/fritzctl_win/).

## Direct downloads

There are several locations from where one can download the packages, e.g.

*   directly from the [debian repository](https://bintray.com/bpicode/fritzctl_deb/fritzctl)
    or the [directory index](https://dl.bintray.com/bpicode/fritzctl_deb/)
*   directly from the [rpm repository](https://bintray.com/bpicode/fritzctl_rpm/fritzctl)
    or the [directory index](https://dl.bintray.com/bpicode/fritzctl_rpm/),
*   directly from the [windows repository](https://bintray.com/bpicode/fritzctl_win/fritzctl)
    or the [directory index](https://dl.bintray.com/bpicode/fritzctl_win/).

## Usage

![Demo usage](/images/fritzctl_demo.gif?raw=true "Demo usage")
