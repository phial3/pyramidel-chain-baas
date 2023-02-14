#!/bin/bash
[[ -n $DEBUG ]] && set -x
set -o errtrace # Make sure any error trap is inherited
set -o nounset  # Disallow expansion of unset variables
set -o pipefail # Use last non-zero exit code in a pipeline

TMP_DIR="$(rm -rf /tmp/baas* && mktemp -d -t baas.XXXXXXXXXX)"
LOG_FILE="${TMP_DIR}/baas.log"

trap trap::info 1 2 3 15 EXIT

######################################################################################################
# function
######################################################################################################

function trap::info() {
  # 信号处理

  [[ ${#ERROR_INFO} -gt 37 ]] && echo -e "$ERROR_INFO"
  [[ ${#ACCESS_INFO} -gt 38 ]] && echo -e "$ACCESS_INFO"
  [ -f "$LOG_FILE" ] && echo -e "\n\n  See detailed log >>> $LOG_FILE \n\n"
  trap '' EXIT
  exit
}

function log::info() {
  # 基础日志
  printf "[%s]: \033[32mINFO:    \033[0m%s\n" "$(date +'%Y-%m-%dT%H:%M:%S.%N%z')" "$*" | tee -a "$LOG_FILE"
}

function log::warning() {
  # 警告日志

  printf "[%s]: \033[33mWARNING: \033[0m%s\n" "$(date +'%Y-%m-%dT%H:%M:%S.%N%z')" "$*" | tee -a "$LOG_FILE"
}

function log::error() {
  # 错误日志

  local item
  item="[$(date +'%Y-%m-%dT%H:%M:%S.%N%z')]: \033[31mERROR:   \033[0m$*"
  ERROR_INFO="${ERROR_INFO}${item}\n  "
  echo -e "${item}" | tee -a "$LOG_FILE"
}

function log::exec() {
  # 执行日志

  printf "[%s]: \033[34mEXEC:    \033[0m%s\n" "$(date +'%Y-%m-%dT%H:%M:%S.%N%z')" "$*" >>"$LOG_FILE"
}

function check::command_exists() {
  # 检查命令是否存在

  local cmd=${1}
  local package=${2}
  if command -V "$cmd" >/dev/null 2>&1; then
    log::info "[check]" "$cmd command exists."
  else
    log::warning "[check]" "I require $cmd but it's not installed."
    log::warning "[check]" "install $package package."
    case $cmd in
    docker)
      command::exec "install::docker"
      check::exit_code "$?" "check" "$package install" "exit"
      ;;
    juicefs)
      command::exec "curl -sSL https://d.juicefs.com/install | sh -"
      check::exit_code "$?" "check" "$package install" "exit"
      ;;
    *)
      command::exec "yum install -y ${package}"
      check::exit_code "$?" "check" "$package install" "exit"
      ;;
    esac
  fi
}

function check::command() {
  # 检查用到的命令
  check::command_exists curl curl
  check::command_exists ssh openssh-clients
  check::command_exists sshpass sshpass
  check::command_exists wget wget
  check::command_exists juicefs juicefs
  check::command_exists tar tar
  check::command_exists docker docker
}

function check::exit_code() {
  # 检查返回码
  local code=${1:-}
  local app=${2:-}
  local desc=${3:-}
  local exit_script=${4:-}
  if [[ "${code}" == "0" ]]; then
    log::info "[${app}]" "${desc} succeeded."
  else
    log::error "[${app}]" "${desc} failed."
    [[ "$exit_script" == "exit" ]] && exit "$code"
  fi
}

function command::exec() {
  local command="$*"
  log::exec "[command]" "bash -c $(printf "%s" "${command}")"
  # shellcheck disable=SC2094
  COMMAND_OUTPUT=$(eval "${command}" 2>>"$LOG_FILE" | tee -a "$LOG_FILE")
  local status=$?
  return $status
}

function install::docker() {
  local version="-20.10.17"
  yum clean all
  cat <<EOF >/etc/yum.repos.d/docker-ce.repo
[docker-ce-stable]
name=Docker CE Stable - \$basearch
baseurl=https://mirrors.aliyun.com/docker-ce/linux/centos/$(rpm --eval '%{centos_ver}')/\$basearch/stable
enabled=1
gpgcheck=1
gpgkey=https://mirrors.aliyun.com/docker-ce/linux/centos/gpg
EOF

  [ -f "$(which docker)" ] && yum remove -y docker-ce docker-ce-cli containerd.io
  yum install -y "docker-ce${version}" "docker-ce-cli${version}" containerd.io bash-completion

  [ -f /usr/share/bash-completion/completions/docker ] &&
    cp -f /usr/share/bash-completion/completions/docker /etc/bash_completion.d/

  [ ! -d /etc/docker ] && mkdir /etc/docker

  cat >/etc/docker/daemon.json <<EOF
{
    "data-root": "/var/lib/docker",
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "200m",
        "max-file": "5"
    },
    "default-ulimits": {
        "nofile": {
            "Name": "nofile",
            "Hard": 655360,
            "Soft": 655360
        },
        "nproc": {
            "Name": "nproc",
            "Hard": 655360,
            "Soft": 655360
        }
    },
    "live-restore": true,
    "oom-score-adjust": -1000,
    "max-concurrent-downloads": 10,
    "max-concurrent-uploads": 10,
    "storage-driver": "overlay2",
    "storage-opts": [
        "overlay2.override_kernel_check=true"
    ],
    "exec-opts": [
        "native.cgroupdriver=systemd"
    ],
    "registry-mirrors": [
        "https://yssx4sxy.mirror.aliyuncs.com/"
    ],
    "insecure-registries": [
        "harbor.sxtxhy.com"
    ]
}
EOF

  sed -i "13c ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock" /usr/lib/systemd/system/docker.service
  systemctl daemon-reload
  systemctl enable docker
  systemctl restart docker.service
}

function init_host() {
  # 设置时区并同步时间
  timedatectl set-timezone Asia/Shanghai
  if
    ! crontab -l | grep ntpdate &
    </dev/null
  then
    (
      echo "* 1 * * * ntpdate ntp.aliyun.com >/dev/null 2>&1"
      crontab -l
    ) | crontab
  fi

  # 禁用selinux
  sed -i 's#SELINUX=enforcing#SELINUX=disabled#g' /etc/sysconfig/selinux
  sed -i 's#SELINUX=enforcing#SELINUX=disabled#g' /etc/selinux/config
  #systemctl disable --now NetworkManager

  # 关闭防火墙
  if egrep "7.[0-9]" /etc/redhat-release &>/dev/null; then
    systemctl stop firewalld
    systemctl disable firewalld
  elif egrep "6.[0-9]" /etc/redhat-release &>/dev/null; then
    service iptables stop
    chkconfig iptables off
  fi

  # 历史命令显示操作时间
  if ! grep HISTTIMEFORMAT /etc/bashrc; then
    echo 'export HISTTIMEFORMAT="%F %T `whoami` "' >>/etc/bashrc
  fi

  #关闭邮件服务
  systemctl stop postfix && systemctl disable postfix

  # 禁止定时任务发送邮件
  sed -i 's/^MAILTO=root/MAILTO=""/' /etc/crontab

  # 设置最大打开文件数
  if ! grep "* soft nofile 65535" /etc/security/limits.conf &>/dev/null; then
    cat >>/etc/security/limits.conf <<EOF
## PYC baas start
root soft nofile 655360
root hard nofile 655360
root soft nproc 655360
root hard nproc 655360
root soft core unlimited
root hard core unlimited


* soft nofile 655360   #可打开的文件描述符的最大数
* hard nofile 655360   #可打开的文件描述符的最大数
* soft nproc 655360   #单个用户可用的最大进程数量
* hard nproc 655360   #单个用户可用的最大进程数量
* soft memlock unlimited   #系统硬限制无上限
* hard memlock unlimited   #系统软限制无上限
EOF
  fi

  #设置连接数最大值
  cat <<EOF >>/etc/rc.local
ulimit -SHn 65535
EOF

  # 系统内核优化
  cat <<EOF >/etc/sysctl.d/s.conf
net.ipv4.ip_forward = 1        #打开路由转发功能
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
fs.may_detach_mounts = 1
vm.overcommit_memory=1
net.ipv4.conf.all.route_localnet = 1
vm.panic_on_oom=0
fs.inotify.max_user_watches=89100
fs.file-max=52706963        #系统级打开最大文件句柄的数量
fs.nr_open=52706963
net.netfilter.nf_conntrack_max=2310720
net.ipv4.tcp_keepalive_time = 600        #间隔多久发送1次keepalive探测包
net.ipv4.tcp_keepalive_probes = 3        #探测失败后，最多尝试探测几次
net.ipv4.tcp_keepalive_intvl =15        #探测失败后，间隔几秒后重新探测
net.ipv4.tcp_max_tw_buckets = 36000        #针对TIME-WAIT，配置其上限
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_max_orphans = 327680
net.ipv4.tcp_orphan_retries = 3
net.ipv4.tcp_syncookies = 1
net.ipv4.ip_conntrack_max = 65536
net.ipv4.tcp_max_syn_backlog = 16384        #增大SYN队列的长度，容纳更多连接
net.ipv4.tcp_timestamps = 0        #TCP时间戳
net.core.somaxconn = 16384       #已经成功建立连接的套接字将要进入队列的长度
EOF
  sysctl --system

  #配置下载源
  yum -y install wget curl
  curl -o /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-7.repo
  wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo
  sed -i -e '/mirrors.cloud.aliyuncs.com/d' -e '/mirrors.aliyuncs.com/d' /etc/yum.repos.d/CentOS-Base.repo

  # 安装系统性能分析工具及其他
  yum update -y --exclude=kernel*
  yum -y install gcc make autoconf vim net-tools ntpdate sysstat iftop iotop lrzsz glances htop wget jq psmisc yum-utils device-mapper-persistent-data lvm2 git bash-completion tree unzip bzip2 gdisk telnet lsof dstat cmake gcc-c++ zlib-devel openssl-devel pcre pcre-devel curl fontconfig ipvsadm ipset sysstat conntrack libseccomp

  # 更新系统及内核
  wget 8.142.32.48/pack/kernel-ml-devel-4.19.12-1.el7.elrepo.x86_64.rpm
  wget 8.142.32.48/pack/kernel-ml-4.19.12-1.el7.elrepo.x86_64.rpm
  yum localinstall -y kernel-ml* && rm -rf kernel-ml*
  grub2-set-default 0 && grub2-mkconfig -o /etc/grub2.cfg
  grubby --args="user_namespace.enable=1" --update-kernel="$(grubby --default-kernel)"

  # 配置ipvs模块
  cat <<EOF >>/etc/modules-load.d/ipvs.conf
ip_vs
ip_vs_lc
ip_vs_wlc
ip_vs_rr
ip_vs_wrr
ip_vs_lblc
ip_vs_lblcr
ip_vs_dh
ip_vs_sh
ip_vs_fo
ip_vs_nq
ip_vs_sed
ip_vs_ftp
ip_vs_sh
nf_conntrack
ip_tables
ip_set
xt_set
ipt_set
ipt_rpfilter
ipt_REJECT
ipip
EOF

  #   减少SWAP使用
  swapoff -a && sysctl -w vm.swappiness=0
  sed -ri '/^[^#]*swap/s@^@#@' /etc/fstab
}

function juicefs::mount() {
  juicefs mount --background --cache-size 512000 redis://:Txhy2020@39.100.224.84:7000/1 /root/txhyjuicefs

  cp /usr/local/bin/juicefs /sbin/mount.juicefs

  cat >>/etc/fstab <<EOF
redis://:Txhy2020@39.100.224.84:7000/1    ~/txhyjuicefs       juicefs     _netdev,max-uploads=50,writeback,cache-size=512000     0  0
EOF
}

function psutil::up() {
  if [ -f /root/txhyjuicefs/psutil/linux/psutil ]; then
    nohup /root/txhyjuicefs/psutil/linux/psutil -port=8082 >nohub.out 2>&1 &
  else
    echo "不存在,处理"
  fi
  netstat -nultp | grep 8082
}

init_host
[ -f "$(which docker)" ] && yum remove -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin.x86_64
check::command
juicefs::mount
psutil::up
