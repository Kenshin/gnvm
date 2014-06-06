GNVM: Node.js version manager on Windows by GO [![Build Status](https://travis-ci.org/Kenshin/gnvm.svg?branch=master)](https://travis-ci.org/Kenshin/gnvm)
================================
`GNVM` Windows下的Node.js多版本管理器，类似 `nvm` `nvmw` `nodist`

文档
---
[English](https://github.com/kenshin/gnvm/blob/master/README.md)

网盘下载
---
* [32-bit](https://app.box.com/gnvm/1/2014967291) | [64-bit](https://app.box.com/gnvm/1/2014967689) **常用地址，Box，速度稍慢**
* [32-bit](http://pan.baidu.com/s/1gdmVgen#dir/path=%2F%E6%88%91%E7%9A%84%E5%85%B1%E4%BA%AB%2Fgnvm%2F32-bit) | [64-bit](http://pan.baidu.com/s/1gdmVgen#dir/path=%2F%E6%88%91%E7%9A%84%E5%85%B1%E4%BA%AB%2Fgnvm%2F64-bit) **备用地址，百度网盘，速度快**
* [32-bit](https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe) | [64-bit](https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe) **备用地址，Github**

其他方式
---
* 已经安装了go的用户，使用 **go get**

  `go get github.com/Kenshin/gnvm`
* 已经安装了git的用户，使用 **git clone**

  `git clone git@github.com:Kenshin/gnvm-bin.git`
* 已经安装了curl的用户，使用 **curl -O**

  `curl -O https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe`

  `curl -O https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe`

配置
---

#### 本机已有node.exe
* 方式1: 将下载的`gnvm.exe`放到`node.exe`目录下。（**推荐方式**）
* 方式2: 将下载的`gnvm.exe`放到任意文件夹下。（确保此文件夹在Path环境下，或者手动添加此文件夹到Path环境）

#### 本机没有node.exe
* 将下载的`gnvm.exe`放到任意文件夹下。（确保此文件夹在Path环境下，或者手动添加此文件夹到Path环境）

验证
---
在cmd下（确保获取管理员权限），输入：`gnvm version`，如有`Current version x.x.x`则说明配置成功。（注：`x.xx.xx`以下载的版本为准。）

![gnvm version](http://i.imgur.com/hEyXZnl.png)

术语
---
* `global` 当前使用的node.exe。
* `latest` 稳定版本的node.exe。

使用
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

最佳实践
---
* 在cmd中运行`gnvm`需要管理员权限。
* 虽然可以直接使用`gnvm`的各种命令，但第一次运行`gnvm`时，建议使用`gnvm config INIT`来初始化一些配置参数。
* 虽然`gnvm`支持任意文件夹，但建议将`node.exe`与`gnvm.exe`放在同一目录下。
* 使用`gnvm config registry xxx`更换库，默认库：<http://nodejs.org/dist/>，只要xxx的结构与默认库一致即可。
* `gnvm`的使用依赖与`.gnvmrc`，请不要手动修改此文件。
* `gnvm install npm`支持安装最新版的npm，但`0.1.0`版本只支持安装最新版本到`node.exe`所在文件夹，不可自定义`npm`的文件夹。（npm的最新版本取决于`gnvm config registry`对应的最新版本。）

使用场景之一（本机已有node.exe）
---
    gnvm config INIT （第一次使用时，推荐做法）
    gnvm config registry dist.u.qiniudn.com （更换库）
    gnvm update latest （如果本机的latest过低，可以使用此方式升级。或者使用gnvm install latest）
    gnvm node-version（查看本机global与latest node.exe版本）
    gnvm install 0.11.1 0.11.2 0.11.3 （下载任意版本的node.exe）
    gnvm use 0.11.1 （切换本机已安装的任意版本node.exe）
    gnvm ls （查看当前共有多少个node.exe）
    gnvm uninstall 0.11.1 （删除0.11.1）

使用场景之二（本机没有node.exe）
---
    gnvm config INIT （第一次使用时，推荐做法）
    gnvm config registry dist.u.qiniudn.com （更换库）
    gnvm install latest -g （下载最新版本的latest并设置为全局node.exe）
    gnvm node-version（查看本机global与latest node.exe版本）
    gnvm ls （查看当前共有多少个node.exe）
    gnvm install npm （安装最新版本的npm到node.exe所在目录）

使用的第三方LIB
---
* <https://github.com/spf13/cobra>
* <https://github.com/tsuru/config>
* <https://github.com/pierrre/archivefile>
* <https://github.com/daviddengcn/go-colortext>
* icon <http://www.easyicon.net/1143807-update_icon.html>

功能一览
---
![功能一览](https://trello-attachments.s3.amazonaws.com/535f6fd8cb08b7fd799c2051/53606254da7b8f8b2f6c9d87/981x580/f6e58f47691d3d352f0b97ba94263df8/gnvm_0.1.0.png)

FAQ
---

#### Q. 在安装了XXX卫士的某些Windows系统下，使用诸如`gnvm use x.xx.xx`的命令会弹出警告。
A. 建议将`gnvm.exe`加入白名单。

#### Q. `gnvm`与`nvmw` `nvm`有什么区别？
A. `gnvm`是单文件CLI，同时比`nvmw`多了一些实用功能，如`gnvm update`, `gnvm install npm`, `gnvm config registry xxxx`等，在功能上更贴近`nvm`。

Help
---
* Email <kenshin@ksria.com>
* Github issue

CHANGELOG
---
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
