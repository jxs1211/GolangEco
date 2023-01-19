# 无法拉取 gcr.io 镜像？用魔法来打败魔法

Original 邹俊豪 [gopher云原生](javascript:void(0);) *2022-01-24 12:12*

收录于合集

\#Docker1个

\#GitHub Actions1个

目前常用的 Docker Registry 公开服务有：

- `docker.io` ：Docker Hub 官方镜像仓库，也是 Docker 默认的仓库
- `gcr.io`、`k8s.gcr.io` ：谷歌镜像仓库
- `quay.io` ：Red Hat 镜像仓库
- `ghcr.io` ：GitHub 镜像仓库

当使用 `docker pull 仓库地址/用户名/仓库名:标签` 时，会前往对应的仓库地址拉取镜像，标签无声明时默认为 `latest`， 仓库地址无声明时默认为 `docker.io` 。

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmf5ZOpdWS1AaobQXkJLjwiaibQvpEWNHkEhPSQlicAOHjTZyibKpq24XnY9DArURAsibicojkFRia6NNHKLw/640?wx_fmt=gif&wxfrom=13&tp=wxpic)

众所周知的原因，在国内访问这些服务异常的慢，甚至 `gcr.io` 和 `quay.io` 根本无法访问。

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmf5ZOpdWS1AaobQXkJLjwiaibHeNYUv3G9eibY6F0hXOf39DS3TibQcXib4K2QCzmYlkQoa1QE6cTkFGzw/640?wx_fmt=gif&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

## 解决方案：镜像加速器

针对 `Docker Hub` ，Docker 官方和国内各大云服务商均提供了 Docker 镜像加速服务。

你只需要简单配置一下（以 Linux 为例）：

```
sudo mkdir -p /etc/docker

sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["镜像加速器"]
}
EOF

sudo systemctl daemon-reload
sudo service docker restart
```

便可以通过访问国内镜像加速器来加速 `Docker Hub` 的镜像下载。

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmf5ZOpdWS1AaobQXkJLjwiaibgJa1k7NHGYZeeshhNbwKLJwicozyMP9iadBLQkHD1gOaibeMsNkwOaCWQ/640?wx_fmt=gif&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

不过这种办法也只能针对 `docker.io` ，其它的仓库地址并没有真正实际可用的加速器（至少我目前没找到）。

## 解决方案：用魔法打败魔法

既然无法治本，那治治标还是可以的吧。

若我们使用一台魔法机器从 `gcr.io` 或 `quay.io` 等仓库先把我们无法下载的镜像拉取下来，然后重新上传到 `docker.io` ，是不是就可以使用 `Docker Hub` 的镜像加速器来下载了。

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmf5ZOpdWS1AaobQXkJLjwiaibzEje2Q6fr8eDAibE4guvvhGrU4FNt2qcmgwGxSpU2Zy4gdNsoiar3F5w/640?wx_fmt=gif&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

镜像仓库迁移的功能，我这里采用了 Go Docker SDK ，整体实现也比较简单。

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmf5ZOpdWS1AaobQXkJLjwiaibu2FAw7gAsqH92htgEicW3aADk9znRG4Sq6cbKBF8vibahXIQAzL4SUOA/640?wx_fmt=gif&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

以需要转换的 `gcr.io/google-samples/microservices-demo/emailservice:v0.3.5` 为例，使用方式：

![Image](data:image/svg+xml,%3C%3Fxml version='1.0' encoding='UTF-8'%3F%3E%3Csvg width='1px' height='1px' viewBox='0 0 1 1' version='1.1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'%3E%3Ctitle%3E%3C/title%3E%3Cg stroke='none' stroke-width='1' fill='none' fill-rule='evenodd' fill-opacity='0'%3E%3Cg transform='translate(-249.000000, -126.000000)' fill='%23FFFFFF'%3E%3Crect x='249' y='126' width='1' height='1'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E)

功能实现了，剩下的就是找台带有魔法的机器了。

GitHub Actions 就是个好选择，我们可以利用提交 `issues` 来触发镜像仓库迁移的功能。

`workflow` 的实现如下：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmf5ZOpdWS1AaobQXkJLjwiaibQ3MztF0v7gqtfG17JTwuTmsucrsx5LuI5AxLWjaFW9icTeaMLjBSv3A/640?wx_fmt=gif&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

实际的使用效果：

![Image](https://mmbiz.qpic.cn/mmbiz_png/Ub8Xn54XTmf5ZOpdWS1AaobQXkJLjwiaib6oJlfPTmL9BGVcBxZnDmpqWXcMHp0Lia2Kh4MpHbIB8bib8TicR4RQSWQ/640?wx_fmt=gif&tp=wxpic&wxfrom=5&wx_lazy=1&wx_co=1)

只要执行最终输出的命令，就可以飞快的使用 Docker Hub 的加速器下载 `gcr.io` 或 `quay.io` 等镜像了。

## 最后

本篇的实现已放在 GitHub ：**`https://github.com/togettoyou/hub-mirror`**