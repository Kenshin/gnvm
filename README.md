GNVM: Node.js version manager on Windows by GO [![Build Status](https://travis-ci.org/Kenshin/gnvm.svg?branch=master)](https://travis-ci.org/Kenshin/gnvm)
================================
`GNVM` is simple multi node.js version manager on Windows by GO, like `nvm` `nvmw` `nodist`

Documentation
---
[中文版](https://github.com/kenshin/gnvm/blob/master/README_CN.md)

Download
---
* [32-bit](https://app.box.com/gnvm/1/2014967291) | [64-bit](https://app.box.com/gnvm/1/2014967689) **Host by Box.com**
* [32-bit](https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe?raw=true) | [64-bit](https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe?raw=true) **Host by Github.com**

Geek way
---
* For go user, usage **go get**

  `go get github.com/Kenshin/gnvm`
* For git user, usage **git clone**

  `git clone git@github.com:Kenshin/gnvm-bin.git`
* For curl user, usage **curl -O**

  `curl -L https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe?raw=true -o gnvm.exe`

  `curl -L https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe?raw=true -o gnvm.exe`


Installation
---

#### exist node.exe
* Download `gnvm.exe` in `node.exe` folder.(**recommended**)
* Download `gnvm.exe` in any folder, add this folder to `path` environment variable.

#### not exist node.exe
* Download `gnvm.exe` in any folder, add this folder to `path` environment variable.

Validation
---
Run `cmd`(administrator permissions) and input `gnvm version`, if output print `Current version x.x.x` configuration is successful.

![gnvm version](http://i.imgur.com/hEyXZnl.png)

Definitions
---
* `global` The current node.exe version.
* `latest` The stable node.exe version.

Usage
---

    Usage:
      gnvm
      gnvm [command]

    Available Commands:
      version                   :: Print the version number of gnvm.exe
      install                   :: Install any node.exe version
      uninstall                 :: Uninstall local node.exe version
      use                       :: Use any version of the local already exists
      update                    :: Update latest node.exe
      ls                        :: List show all <local> <remote> node.exe version
      node-version              :: Show <global> <latest> node.exe version
      config                    :: Setter and getter registry
      help [command]            :: Help about any command

Best practices
---
* Run `gnvm` need administrator permissions.
* The first run `gnvm` need use `gnvm config INIT`(**recommended**)
* `gnvm` can download to any folder, suggest `gnvm.exe` in `node.exe` folder.
* Use `gnvm config registry xxx` change reigistry, default registry is <http://nodejs.org/dist/>.
* `gnvm` depend on`.gnvmrc`, please don't modify manually.
* `gnvm install npm` support latest npm, but `0.1.0` version only support install npm to `node.exe` folder, can't custom npm path.

Usage scenarios( exist node.exe )
---
    gnvm config INIT
    gnvm config registry dist.u.qiniudn.com
    gnvm update latest
    gnvm node-version
    gnvm install 0.11.1 0.11.2 0.11.3
    gnvm use 0.11.1
    gnvm ls
    gnvm uninstall 0.11.1

Usage scenarios( not exist node.exe )
---
    gnvm config INIT
    gnvm config registry dist.u.qiniudn.com
    gnvm install latest -g
    gnvm node-version
    gnvm ls
    gnvm install npm

Dependency
---
* <https://github.com/Kenshin/curl>
* <https://github.com/Kenshin/cprint>

Other package
---
* <https://github.com/spf13/cobra>
* <https://github.com/tsuru/config>
* <https://github.com/pierrre/archivefile>
* <https://github.com/daviddengcn/go-colortext>

Icon
---
* <http://www.easyicon.net/1143807-update_icon.html>

Feature
---
![Feature](https://trello-attachments.s3.amazonaws.com/535f6fd8cb08b7fd799c2051/53606254da7b8f8b2f6c9d87/981x580/f6e58f47691d3d352f0b97ba94263df8/gnvm_0.1.0.png)

FAQ
---

#### Q. The difference between `gnvm` and `nvmw`, `nvm`?
A. `gnvm` is single cli file, more than `nvmw` feature, e.g. `gnvm update`, `gnvm install npm`, `gnvm config registry xxxx`, more like `nvm`.

Help
---
* Email <kenshin@ksria.com>
* Github issue

CHANGELOG
---
* **2014-07-15, Version `0.1.2`:**
    * Adapter go version 1.3.
    * Fix bug of usage `gnvm update latest -g` adapter go version error.

* **2014-06-06, Version `0.1.1`:**
    * Change `util/p/print.go` to `github.com/Kenshin/cprint`.
    * change `util/curl/curl.go` to `github.com/Kenshin/curl`.
    * Add this project to travis-ci.org.
    * Remove `nodehandle.cmd` method.
    * Optimize `nodehandle.copy` method logic.
    * Fix bug of When not global node.exe, need get gnvm.exe path.

* **2014-05-30, Version `0.1.0`:**
    * version
    * install
    * uninstall
    * use
    * update
    * ls
    * node-version
    * config

LICENSE
---
(The MIT License)

Copyright (c) 2014 Kenshin Wang <kenshin@ksria.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
