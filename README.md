GNVM - 使用 Go 语言编写的 Node.js 多版本管理器  
================================  
[![Travis][travis-badge]][travis-link]
[![Version][version-badge]][version-link]
[![Gitter][gitter-badge]][gitter-link]
[![Slack][slack-badge]][slack-link]
[![Jianliao][jianliao-badge]][jianliao-link]  
#### `GNVM` 是一个简单的 `Windows` 下 Node.js 多版本管理器，类似的 `nvm` `nvmw` `nodist` 。  
```
c:\> gnvm install latest 1.0.0-x86 1.0.0-x64 5.0.0
Start download Node.js versions [5.10.1, 1.0.0, 1.0.0-x86, 5.0.0].
5.10.1: 18% [=========>__________________________________________] 4s
 1.0.0: 80% [==========================================>_________] 40s
1.0...: 50% [==========================>_________________________] 30s
 5.0.1: 100% [==================================================>] 20s
End download.

c:\> gnvm ls
5.1.1 -- latest
1.0.0
1.0.0 -- x86
5.0.0 -- global

c:\> gnvm use latest
Set success, current Node.js version is 5.10.0.

c:\> gnvm update latest
Update success, current Node.js latest version is 5.10.0.
```

特色
---
* 单文件，不依赖于任何环境。
* 下载即用，无需配置。
* 彩色日志输出。
* 支持多线程下载。
* 内置 [TAOBAO](http://npm.taobao.org/mirrors/node)，方便切换，也支持自定义。
* 支持 `NPM` 下载/安装/配置。

主页
---
[![Website][www-badge]][www-link]

文档
---
[English](https://github.com/kenshin/gnvm/blob/master/README.en.md) | [繁體中文](https://github.com/kenshin/gnvm/blob/master/README.tw.md)

下载
---
* [32-bit](https://app.box.com/gnvm/1/2014967291) | [64-bit](https://app.box.com/gnvm/1/2014967689) 常用地址，Box，速度稍慢
* [32-bit](http://pan.baidu.com/s/1gdmVgen#dir/path=%2F%E6%88%91%E7%9A%84%E5%85%B1%E4%BA%AB%2Fgnvm%2F32-bit) | [64-bit](http://pan.baidu.com/s/1gdmVgen#dir/path=%2F%E6%88%91%E7%9A%84%E5%85%B1%E4%BA%AB%2Fgnvm%2F64-bit) 备用地址，百度网盘，速度快
* [32-bit](https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe?raw=true) | [64-bit](https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe?raw=true) 备用地址，Github

* 已经安装了go的用户，使用  
  `go get github.com/Kenshin/gnvm`

* 已经安装了git的用户，使用  
  `git clone git@github.com:Kenshin/gnvm-bin.git`

* 已经安装了curl的用户，使用  
  `curl -L https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe?raw=true -o gnvm.exe`  
  `curl -L https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe?raw=true -o gnvm.exe`

安装
---
* 不存在 Node.js 环境
  > 下载并解压缩 `gnvm.exe` 保存到任意文件夹，并将此文件夹加入到环境变量 `Path` 。

* 存在 Node.js 环境
  > 下载并解压缩 `gnvm.exe` 保存到 `Node.js` 所在的文件夹。

验证
---
* 在 `cmd` 下，输入 `gnvm version`，输出 `版本说明` 则配置成功。

功能
---
```
config       Setter and getter .gnvmrc file
use          Use any the local already exists of Node.js version
ls           Show all [local] [remote] Node.js version
install      Install any Node.js version
uninstall    Uninstall local Node.js version and npm
update       Update Node.js latest version
npm          NPM version management
session      Use any Node.js version of the local already exists version by current session
search       Search and Print Node.js version detail usage wildcard mode or regexp mode
node-version Show [global] [latest] Node.js version
reg          Add config property 'noderoot' to Environment variable 'NODE_HOME'
version      Print GNVM version number
```
![功能一览](http://i.imgur.com/GqkZcjZ.png)

术语
---
* `global`   当前使用的 `Node.js` 。
* `latest`   稳定版本的 `Node.js` 。
* `session`  当前 `cmd` 所对应的环境。（临时环境）
* `.gnvmrc`  `gnvm`配置文件，无需手动建立，其中保存了 `本地` / `远程` Node.js 版本信息等。
    - `registry` 下载 `node.exe` 所对应的库，默认为 [DEFAULT](http://nodejs.org/dist/)，可以更换为 [TAOBAO](http://npm.taobao.org/mirrors/node)，也支持自定义。（**自定义库的结构需要保持一致。**）
    - `noderoot` 保存了全局 `Node.js` 所在的目录。（也是 `gnvm.exe` 所在的目录。）

入门指南
---
> `gnvm.exe` 是一个单文件 exe，无需任何配置，直接使用。

**.gnvmrc**

```
globalversion: 5.0.1
latestversion: 5.10.1
noderoot: /Users/kenshin/Work/28-GO/01-work/src/gnvm
registry: http://npm.taobao.org/mirrors/node/
```

**更换更快的库 registry**
  > `gnvm.exe` 内建了 [DEFAULT](http://nodejs.org/dist/) 和 [TAOBAO](http://npm.taobao.org/mirrors/node) 两个库。

```
gnvm config registry TAOBAO
```

**安装 多个 Node.js**
  > 安装任意版本的 Node.js 包括： 自动匹配 `latest` / `io.js` version 以及 选择 32 / 64 位，例如 `x.xx.xx-x64` 。

```
gnvm install latest 1.0.0-x86 1.0.0-x64 5.0.0
```

**卸载本地任意 Node.js 版本**
```
gnvm uninstall latest 1.0.0-x86 1.0.0-x64 5.0.0
```

**切换本地存在的任意版本 Node.js**
```
gnvm use 5.10.1
```

**列出本地已存在的全部 Node.js 版本**
```
c:\> gnvm ls
5.1.1 -- latest
1.0.0
1.0.0 -- x86
5.0.0 -- global
```

**更新本地的 Node.js latest 版本**
```
gnvm update latest
```

**安装 NPM**
  > `gnvm` 支持安装 `npm`, 例如：下载最新版的 npm version ，使用 `gnvm npm latest` 。

```
gnvm npm latest
```

**查询 Node.js 版本**
  > 可以使用关键字 `*` 或者 正则表达式 `/regxp/`，例如： `gnvm search 5.*.*` 或者 `gnvm search /.10./` 。

```
c:\> gnvm search 5.*.*
Search Node.js version rules [5.x.x] from http://npm.taobao.org/mirrors/node/index.json, please wait.
+--------------------------------------------------+
| No.   date         node ver    exec      npm ver |
+--------------------------------------------------+
1     2016-04-05   5.10.1      x86 x64   3.8.3
2     2016-04-01   5.10.0      x86 x64   3.8.3
3     2016-03-22   5.9.1       x86 x64   3.7.3
4     2016-03-16   5.9.0       x86 x64   3.7.3
5     2016-03-09   5.8.0       x86 x64   3.7.3
6     2016-03-02   5.7.1       x86 x64   3.6.0
7     2016-02-23   5.7.0       x86 x64   3.6.0
+--------------------------------------------------+
```

例子
---
**1. 不存在 Node.js 环境时，下载 Node.js latest version 并设置为全局 Node.js 。**
```
c:\> gnvm config registry TAOBAO
Set success, registry new value is http://npm.taobao.org/mirrors/node/
c:\> gnvm install latest -g
Notice: local  latest version is unknown.
Notice: remote latest version is 5.10.1.
Start download Node.js versions [5.10.1].
5.10.1: 100% [==================================================>] 13s
End download.
Set success, latestversion new value is 5.10.1
Set success, global Node.js version is 5.10.1.
```

**2. 升级本地 Node.js latest 版本。**
```
c:\> gnvm config registry TAOBAO
Set success, registry new value is http://npm.taobao.org/mirrors/node/
c:\> gnvm update latest
Notice: local  Node.js latest version is 5.9.1.
Notice: remote Node.js latest version is 5.10.1 from http://npm.taobao.org/mirrors/node/.
Waring: remote latest version 5.10.1 > local latest version 5.9.1.
Waring: 5.10.1 folder exist.
Update success, Node.js latest version is 5.10.1.
```

**3. 查看本地 Node.js global and latest 版本。**
```
c:\> gnvm node-version
Node.js latest version is 5.10.1.
Node.js global version is 5.10.1.
```

**4. 验证 .gnvmrc registry 正确性。**
```
c:\> gnvm config registry test
Notice: gnvm config registry http://npm.taobao.org/mirrors/node/ valid ................... ok.
Notice: gnvm config registry http://npm.taobao.org/mirrors/node/index.json valid ......... ok.
```

**5. 本地不存在 NPM 时，安装当前 Node.js 版本对应的 NPM 版本。**
```
c:\ gnvm npm global
Waring: current path C:\xxx\xxx\nodejs\ not exist npm.
Notice: local    npm version is unknown
Notice: remote   npm version is 3.8.3
Notice: download 3.8.3 version [Y/n]? y
Start download new npm version v3.8.3.zip
v3.8.3.zip: 100% [==================================================>] 4s
Start unzip and install v3.8.3.zip zip file, please wait.
Set success, current npm version is 3.8.3.
c:\> npm -v
3.8.7
```

**6. 安装 NPM latest 版本。**
```
c:\ gnvm npm laltest
Notice: local    npm version is 3.7.3
Notice: remote   npm version is 3.8.7
Notice: download 3.8.7 version [Y/n]? y
Start download new npm version v3.8.7.zip
v3.8.7.zip: 100% [==================================================>] 3s
Start unzip and install v3.8.7.zip zip file, please wait.
Set success, current npm version is 3.8.7.
c:\> npm -v
3.8.7
```

依赖
---
* <https://github.com/Kenshin/curl>
* <https://github.com/Kenshin/cprint>
* <https://github.com/Kenshin/regedit>

第三方包
---
* <https://github.com/spf13/cobra>
* <https://github.com/tsuru/config>
* <https://github.com/pierrre/archivefile>
* <https://github.com/daviddengcn/go-colortext>
* <https://github.com/bitly/go-simplejson>

下一步
---
- [ ] 增加 `HTTP_PROXY` 。
- [ ] 自动升级，例如 `gnvm upgrad` .
- [ ] `gnvm.exe` 增加 `Chocolatey` 方案。

相关链接
---
* [更新日志](https://github.com/kenshin/gnvm/blob/master/CHANGELOG.md)
* [联系方式](http://kenshin.wang/) | [邮件](kenshin@ksria.com) | [微博](http://weibo.com/23784148)
* [反馈](https://github.com/kenshin/gnvm/issues)

感谢
---
* 图标来自 <http://www.easyicon.net> 。
* 页面设计参考 [You-Get](https://you-get.org/) 。

许可
---
[![license-badge]][license-link]

<!-- Link -->
[www-badge]:        https://img.shields.io/badge/website-gnvm.ksria.com-1DBA90.svg
[www-link]:         http://ksria.com/gnvm
[version-badge]:    https://img.shields.io/badge/lastest_version-0.2.0-blue.svg
[version-link]:     https://github.com/kenshin/gnvm/releases
[travis-badge]:     https://travis-ci.org/Kenshin/gnvm.svg?branch=master
[travis-link]:      https://travis-ci.org/Kenshin/gnvm
[gitter-badge]:     https://badges.gitter.im/kenshin/gnvm.svg
[gitter-link]:      https://gitter.im/kenshin/gnvm?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge
[slack-badge]:      https://img.shields.io/badge/chat-slack-orange.svg
[slack-link]:       https://gnvm.slack.com/
[jianliao-badge]:   https://img.shields.io/badge/chat-jianliao-yellowgreen.svg
[jianliao-link]:    https://guest.jianliao.com/rooms/76dce8b01v
[license-badge]:    https://img.shields.io/github/license/mashape/apistatus.svg
[license-link]:     https://opensource.org/licenses/MIT
