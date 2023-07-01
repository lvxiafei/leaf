安装
--

macOS（支持 Sierra 以上版本）可以直接通过 Homebrew 安装：

```bash
$ brew cask install multipass
```

上手例子
--

先来创建一个容器：

```
$ multipass launch --name ubuntu
Launched: ubuntu
```

初次创建时需要下载镜像，网络畅通的情况下，稍等片刻即可。

容器创建后 multipass 会马上启动它，这样创建好容器后我们就可以直接使用了：

```
$ multipass exec ubuntu -- lsb_release -d
Description:	Ubuntu 22.04.3 LTS
```

`lsb_release` 会打印 Linux 发行版的信息。 之前我们创建容器的时候并没有指定使用什么样的镜像，上面命令的输出表明，multipass 默认会使用当前 LTS 版本的 Ubuntu。

除了直接在容器上运行（`exec`）命令外，还可以通过 `shell` 命令「进入」容器：

```
$ multipass shell ubuntu
```

定制
--

```
$ multipass launch --name ubuntu --disk 40G --cloud-init config.yaml 22.04
```

容器创建成功后，通过 `multipass info` 可以查看容器的基本信息， 至于 `config.yaml` 则是容器的初始化配置文件

手动挂载
--

```
$ multipass mount /some/local/path:/some/instance/path


# 挂载也可以在launch命令中作为选项之一：
$ multipass launch --mount /some/local/path:/some/instance/path
```

更多
--

运行 `multipass list` 可以列出所有的容器

结语
--

Multipass 使用起来十分简洁直观。 它是由 Canonical （Ubuntu 背后的公司）推出的，因此使用的镜像由 Canonical 负责更新，包含最近的安全更新，以及专门为各个平台的虚拟化方案（Windows 的 Hyper-V、macOS 的 HyperKit、Linux 的 KVM）优化的内核。 不过也因为同样的原因，目前支持的镜像也只限于 Ubuntu。
