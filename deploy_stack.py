#!/usr/bin/env python
__metaclass__ = type

import yaml

from collections import defaultdict
from cStringIO import StringIO
from datetime import datetime, timedelta
import subprocess
import sys


class ErroredUnit(Exception):

    def __init__(self, unit_name, state):
        Exception.__init__('Unit %s is in state %s' % unit_name, state)


class Environment:

    def __init__(self, environment):
        self.environment = environment

    def _full_args(self, command, *args):
        return ('juju', command, '-e', self.environment) + args

    def juju(self, command, *args):
        args = self._full_args(command, *args)
        print ' '.join(args)
        sys.stdout.flush()
        return subprocess.check_call(args)

    @staticmethod
    def agent_states(status):
        states = defaultdict(list)
        for machine_name, machine in sorted(status['machines'].items()):
            states[machine.get('agent-state', 'no-agent')].append(machine_name)
        for service in sorted(status['services'].values()):
            for unit_name, unit in service['units'].items():
                states[unit['agent-state']].append(unit_name)
        return states

    def wait_for_started(self):
        now = datetime.now()
        args = self._full_args('status')
        while now - datetime.now() < timedelta(300):
            status = yaml.safe_load(StringIO(subprocess.check_output(args)))
            states = self.agent_states(status)
            pending = False
            state_listing = []
            for state, entries in states.items():
                if state == 'started':
                    continue
                if 'error' in state:
                    raise ErroredUnit(entries[0],  state)
                pending = True
                state_listing.append('%s: %s' % (state, ' '.join(entries)))
            print ' / '.join(state_listing)
            sys.stdout.flush()
            if not pending:
                return
        raise Exception('Timed out!')

        def pending_entry(entry_name, entry):
            if 'error' in entry['agent-state']:
                raise ErroredUnit(entry_name,  entry['agent-state'])
            return entry['agent-state'] != 'started'


def deploy_stack(environment):
    env = Environment(environment)
    env.juju('bootstrap', '--constraints', 'mem=2G')
    env.juju('deploy', 'wordpress')
    env.juju('deploy', 'mysql')
    env.juju('add-relation', 'mysql', 'wordpress')
    env.wait_for_started()


def main():
    try:
        deploy_stack(sys.argv[1])
    except Exception as e:
        print e
        sys.exit(1)


if __name__ == '__main__':
    main()
