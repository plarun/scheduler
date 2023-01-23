#!/bin/bash

# schd_def subcommand for job definition validation and processing

schd_home=/root/go/src/github.com/plarun/scheduler
schd_bin=${schd_home}/bin

schd_cli=${schd_bin}/client
sub_command=schd_def

# call client service
${schd_cli} ${sub_command} $0
