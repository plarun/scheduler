# Scheduler Client

- Client allows the user to interact with event server.
- It provides multiple cli commands for the user to interact with jobs.

---

## CLI commands

### 1. schd_event
- schd_event command allows the user to send an event on existing job for some actions
- it takes two arguments
- arg 1 is job name
- arg 2 is one of the below events
- a. start - it starts the job if its not already running. job shouldn't be "FROZEN"
- b. abort - it stops the job if its already running
- c. frost - it changes job status to "FROZEN"
- d. reset - it resets the job status to "IDLE" - default state
- e. chg_succ - it changes the status of job to "SUCCESS"
- f. chg_fail - it changes the status of job to "FAILED"

### 2. schd_job
- sched_job shows definition of job
- it takes one argument
- arg 1 is job name

### 3. schd_status
- sched_status shows the current status of job
- it takes one argument
- arg 1 is job name

### 4. schd_runs
- sched_runs shows the recent run details of job
- it takes one argument
- arg 1 is job name
- flag -n <number> : shows nth run of the job
- flag -d <date> : shows only runs of specified date
- by default it shows all available runs

### 5. schd_def
- sched_def checks the syntax and job definition in the file
- it takes one argument
- arg 1 is file name
- flag -c : only checks the syntax and validates the definition
- by default it checks the syntax and validates the definition then process the actions

### 6. schd_dep
- sched_dep shows the dependencies of the job
- it takes one argument
- arg 1 is job name
