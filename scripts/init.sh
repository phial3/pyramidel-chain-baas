#!/bin/bash
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
    ]
}
EOF

  sed -i "13c ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock" /usr/lib/systemd/system/docker.service
  systemctl daemon-reload
  systemctl enable docker
  systemctl restart docker.service
}

check::command
