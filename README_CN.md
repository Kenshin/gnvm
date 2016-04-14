GNVM - 使用 Go 语言编写的 Node.js 多版本管理器  
[![Travis][travis-badge]][travis-link]
[![Version][version-badge]][version-link]
[![Gitter][gitter-badge]][gitter-link]
[![Slack][slack-badge]][slack-link]
[![Jianliao][jianliao-badge]][jianliao-link]
================================  
`GNVM` 是一个简单的Windows下 Node.js 多版本管理器，类似的 `nvm` `nvmw` `nodist`

主页
---
[![Website][www-badge]][www-link]

文档
---
[English](https://github.com/kenshin/gnvm/blob/master/README.md) | [繁體中文](https://github.com/kenshin/gnvm/blob/master/README_tw.md)

下载
---
* [32-bit](https://app.box.com/gnvm/1/2014967291) | [64-bit](https://app.box.com/gnvm/1/2014967689) 常用地址，Box，速度稍慢
* [32-bit](http://pan.baidu.com/s/1gdmVgen#dir/path=%2F%E6%88%91%E7%9A%84%E5%85%B1%E4%BA%AB%2Fgnvm%2F32-bit) | [64-bit](http://pan.baidu.com/s/1gdmVgen#dir/path=%2F%E6%88%91%E7%9A%84%E5%85%B1%E4%BA%AB%2Fgnvm%2F64-bit) 备用地址，百度网盘，速度快
* [32-bit](https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe?raw=true) | [64-bit](https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe?raw=true) 备用地址，Github

其他方式
---
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
在cmd下（确保获取管理员权限），输入：`gnvm version`，如有`Current version x.x.x`则说明配置成功。（注：`x.xx.xx`以下载的版本为准。）

术语
---
* `global` 当前使用的node.exe。
* `latest` 稳定版本的node.exe。

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

入门指南
---
> `gnvm.exe` 是一个单文件 exe，无需任何配置，直接使用。

**.gnvmrc**
> `.gnvmrc` 无需手动建立，其中保存了 本地 / 远程 Node.js 版本信息等。

```
globalversion: 5.0.1
latestversion: 5.10.1
noderoot: /Users/kenshin/Work/28-GO/01-work/src/gnvm
registry: http://npm.taobao.org/mirrors/node/
```

**更换更快的库 registry**
  > `gnvm.exe` 内建了 [DEFAULT](http://nodejs.org/dist/) and [TAOBAO](http://nodejs.org/dist/) 两个库。

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

**gnvm npm latest**
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

依赖
---
* <https://github.com/Kenshin/curl>
* <https://github.com/Kenshin/cprint>

第三方包
---
* <https://github.com/spf13/cobra>
* <https://github.com/tsuru/config>
* <https://github.com/pierrre/archivefile>
* <https://github.com/daviddengcn/go-colortext>


下一步
---
- [ ] 增加 `HTTP_PROXY` 。
- [ ] 自动升级，例如 `gnvm upgrad` .
- [ ] `gnvm.exe` 增加 `Chocolatey` 方案。

相关链接
---
* [更新日志](https://github.com/kenshin/gnvm/blob/master/CHANGELOG.md)
* [联系方式](http://kenshin.wang/) | [邮件](kenshin@ksria.com) | [微博](http://weibo.com/23784148)
* [常见问题](https://github.com/kenshin/gnvm/wiki/常见问题)
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
