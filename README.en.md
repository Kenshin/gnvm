![logo](http://i.imgur.com/uPmYWlEs.png) GNVM - Node.js version manager on Windows by Go  
================================  
[![Travis][travis-badge]][travis-link]
[![Version][version-badge]][version-link]
[![Gitter][gitter-badge]][gitter-link]
[![Slack][slack-badge]][slack-link]
[![Jianliao][jianliao-badge]][jianliao-link]  
#### `GNVM` is simple multiple Node.js version manager by Go, like `nvm` `nvmw` `nodist`.  
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

Characteristic
---
* Single file, not dependent on any environment.
* Direct use, no configuration.
* Color stdout.
* Support multiple download.
* Built-in [TAOBAO] (http://npm.taobao.org/mirrors/node), convenient switching, also support custom.
* Support `NPM` download / install.

Website
---
[![Website][www-badge]][www-link]

Document
---
[简体中文](https://github.com/kenshin/gnvm/blob/master/README.md) | [繁體中文](https://github.com/kenshin/gnvm/blob/master/README.tw.md)

Download
---
* [32-bit](https://app.box.com/gnvm/1/2014967291) | [64-bit](https://app.box.com/gnvm/1/2014967689) Host by Box.com
* [32-bit](https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe?raw=true) | [64-bit](https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe?raw=true) Host by Github.com

* For go user, please usage  
  `go get github.com/Kenshin/gnvm`

* For git user, please usage  
  `git clone git@github.com:Kenshin/gnvm-bin.git`

* For curl user, please usage  
  `curl -L https://github.com/Kenshin/gnvm-bin/blob/master/32-bit/gnvm.exe?raw=true -o gnvm.exe`  
  `curl -L https://github.com/Kenshin/gnvm-bin/blob/master/64-bit/gnvm.exe?raw=true -o gnvm.exe`

Installation
---
* Not exist Node.js Environment
  > Download and unzip `gnvm.exe` save to any local folder and add this folder to Environment `Path`.

* Exist Node.js Environment
  > Download and unzip `gnvm.exe` save to the same `Node.js` folder.

Validation
---
* Run `cmd`(administrator permissions) and input `gnvm version`, if output print `gnvm version` configuration is successful.

Feature
---
```
config       Setter and getter .gnvmrc file
use          Use any the local already exists of Node.js version
ls           Show all [local] [remote] Node.js version
install      Install any Node.js version
uninstall    Uninstall local Node.js version and npm
update       Update Node.js latest version
npm          NPM version management
session      Set any local Node.js version to session Node.js version
search       Search and Print Node.js version detail usage wildcard mode or regexp mode
node-version Show [global] [latest] Node.js version
reg          Add config property [noderoot] to Environment variable [NODE_HOME]
version      Print GNVM version number
```
![Feature](http://i.imgur.com/E7MvvQv.png)

Definitions
---
* `global`   current `Node.js` version.
* `latest`   latest `Node.js` version.
* `session`  current `cmd` Environment.( Temporary environment )
* `.gnvmrc`  `gnvm`configure file, can be auto created and it saved local/remote Node.js version information.
    - `registry` `node.exe` download URL, default is [DEFAULT](http://nodejs.org/dist/), can be choose [TAOBAO](http://nodejs.org/dist/), and support custom `url`.
    - `noderoot` save global Node.js path.

Getting Started
---
> `gnvm.exe` is a single exe file, don't need to configure, direct usage.

**.gnvmrc**

```
globalversion: 5.0.1
latestversion: 5.10.1
noderoot: /Users/kenshin/Work/28-GO/01-work/src/gnvm
registry: http://npm.taobao.org/mirrors/node/
```

**Change fast registry**
  > `gnvm.exe` built-in [DEFAULT](http://nodejs.org/dist/) and [TAOBAO](http://nodejs.org/dist/) two registry.

```
gnvm config registry TAOBAO
```

**Install multiple Node.js**
  > Install any Node.js version include: automatic recognition of `latest` version, `io.js` version and specified arch, e.g. `x.xx.xx-x64`.

```
gnvm install latest 1.0.0-x86 1.0.0-x64 5.0.0
```

**Uninstall local Node.js version**
```
gnvm uninstall latest 1.0.0-x86 1.0.0-x64 5.0.0
```

**Usage any local Node.js version**
```
gnvm use 5.10.1
```

**List all local Node.js versions**
```
c:\> gnvm ls
5.1.1 -- latest
1.0.0
1.0.0 -- x86
5.0.0 -- global
```

**Update local Node.js latest version**
```
gnvm update latest
```

**Install npm**
  > `gnvm` support install `npm`, download npm latest version, usage `gnvm npm latest`.

```
gnvm npm latest
```

**Search Node.js version from .gnvmrc registry**
  > you can usage `*` or `/regxp/`, e.g. `gnvm search 5.*.*` or `gnvm search /.10./` .

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

Example
---
**1. Not exist Node.js Environment and download Node.js latest version and usage it.**
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

**2. Update local Node.js latest version.**
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

**3. See Node.js global and latest version.**
```
c:\> gnvm node-version
Node.js latest version is 5.10.1.
Node.js global version is 5.10.1.
```

**4. Verify config registry.**
```
c:\> gnvm config registry test
Notice: gnvm config registry http://npm.taobao.org/mirrors/node/ valid ................... ok.
Notice: gnvm config registry http://npm.taobao.org/mirrors/node/index.json valid ......... ok.
```

**5. Local not exist npm and install local Node.js version matching npm version.**
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

**6. Install latest npm version.**
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

Dependency
---
* <https://github.com/Kenshin/curl>
* <https://github.com/Kenshin/cprint>
* <https://github.com/Kenshin/regedit>

Other package
---
* <https://github.com/spf13/cobra>
* <https://github.com/tsuru/config>
* <https://github.com/pierrre/archivefile>
* <https://github.com/daviddengcn/go-colortext>
* <https://github.com/bitly/go-simplejson>

To-Do
---
- [ ] Add `HTTP_PROXY` .
- [ ] Auto `Upgrade`, usage `gnvm upgrad` .
- [ ] Add `gnvm.exe` to `Chocolatey`.
- [ ] Multiple system. ( `MAC`, `Linux` )

About
---
* [CHANGELOG](https://github.com/kenshin/gnvm/blob/master/CHANGELOG.md)
* [Contact](http://kenshin.wang/) | [Email](kenshin@ksria.com) | [Twitter](https://twitter.com/wanglei001)
* [Feedback](https://github.com/kenshin/gnvm/issues)

Thanks
---
* Icon <http://www.easyicon.net> .
* Theme reference [You-Get](https://you-get.org/) .

Licenses
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
