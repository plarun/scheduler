# Task Scheduler

- Schedules the task to run at the specified batch run.
- Task can be created/updated/deleted by Job definition syntax.
- Client can view existing job's definition.
- Client can manually action on the job.
- Client can view job run status and run history.

Subcommand syntax | Desc
--- | ---
`event <job_name> <event_type = start,abort,freeze,reset,green>` | send event to job where `start` runs the job, `abort` stops the running job, `freeze` scheduler will ignore the task, `reset` sets the job to default mode, `green` change the job status to SUCCESS
`submit <file_path>` | submits the Job definition to create/update/delete job.
`status <job_name>` | view the latest run status or job
`job <job_name>` | view the job definition
`history <job_name>` | view the job run history

***

## Job definition to create new job
![alt text](https://github.com/plarun/scheduler/blob/master/img/sched0.png)

***

## Usage
![alt text](https://github.com/plarun/scheduler/blob/master/img/sched1.png)

***

## Job definition
![alt text](https://github.com/plarun/scheduler/blob/master/img/sched2.png)

***

## Event
![alt text](https://github.com/plarun/scheduler/blob/master/img/sched3.png)

***

## Job run status
![alt text](https://github.com/plarun/scheduler/blob/master/img/sched4.png)

***

## Job run history
![alt text](https://github.com/plarun/scheduler/blob/master/img/sched5.png)