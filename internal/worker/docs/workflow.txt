# Worker

## Roles of worker service
1. Periodically check on sched_ready for tasks to execute
2. Pull ready tasks and put an entry on sched_exec then remove it from sched_ready
3. Execute the tasks in the sched_exec with status 'READY'
4. Once execution is started, update the task status to 'RUNNING'
5. Periodically check the running task's process whether is it alive or terminated by OS
6. Once task's command execution completed, get its exit code
7. If exit code is 0, mark the task status as 'SUCCESS' on sched_exec
8. If exit code is 1, mark the task status as 'FAILURE' on sched_exec
9. Periodically check sched_exec for status in ('SUCCESS', 'FAILURE')

# Locks on sched_ready
- 0 - newly inserted ready tasks
- 1 - locked for readyTaskPull
- 2 - pulled for execution