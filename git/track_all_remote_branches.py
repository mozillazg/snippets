#!/usr/bin/env python
# -*- coding: utf-8 -*-
from __future__ import absolute_import, print_function, unicode_literals

from collections import namedtuple
import re
import subprocess


def run_cmd(cmd):
    process = subprocess.Popen(cmd, stdout=subprocess.PIPE,
                               stderr=subprocess.PIPE, shell=True)
    result = namedtuple('Result', ['stdout', 'stderr'])
    return result(*process.communicate())


def get_remote_branches():
    cmd = 'git branch -r | grep -vE "HEAD|master"'
    result = run_cmd(cmd)
    if result.stderr:
        raise OSError(result.stderr)
    branches = map(lambda x: x.strip(), result.stdout.split())
    for branch in filter(None, branches):
        if not branch.startswith('origin/'):
            continue
        yield branch


def add_local_branch(remote_branches, replace=re.compile('^origin/')):
    remote_branches = filter(replace.match, remote_branches)
    local_branches = map(lambda x: replace.sub('', x), remote_branches)
    cmd = 'git branch --track {local} {remote}'
    for remote, local in zip(remote_branches, local_branches):
        result = run_cmd(cmd.format(local=local, remote=remote))
        print(result.stdout, result.stderr)


def main():
    remote_branches = get_remote_branches()
    add_local_branch(remote_branches, re.compile('^origin/branches/'))

if __name__ == '__main__':
    main()
    # update changes with this shell:
    # for remote in `git branch`; do git checkout $remote ; git pull; done
