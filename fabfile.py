import sys

sys.path.insert(0, r'./libs')

from fabric.api import run,put,env,cd,local,lcd,sudo,settings, hide
from fabric.contrib.files import exists
from fabric.state import connections

#env.hosts = ['ubuntu@staging-agent-1.wixcore3.com']
#env.key_filename=['~/.ssh/core3-ci.pem']

env.hosts = ['vagrant@192.168.1.79']
env.user = 'vagrant'
env.password = 'vagrant'


def buildAndTransfer():
    local('CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=\"-s -w -X main.Version=1.0\" -o claptrap')
    local('upx claptrap')
    put('claptrap','/home/vagrant/claptrap')
    sudo('chmod 755 /home/vagrant/claptrap')

