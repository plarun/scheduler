# Processes in Allocator
- Staging
- Queuing
- Splitting

# Tables involved
- `sched_task` - Tasks information
- `sched_stage` - Keeps track on scheduled task till run is completed
- `sched_queue` - Holds only callable tasks for condition check
- `sched_ready` - Holds the tasks which are ready for running
- `sched_wait` - Holds the waiting task which are not met to start condition

## Staging
### Steps
1. Lock tasks in `sched_task` which are to be scheduled on the current cycle
2. Put an entry for locked tasks into `sched_stage` for staging with flag=0 (locked for staging)
3. For tasks with flag=0 on `sched_stage`, update task status to `staged` on `sched_task`
4. For tasks with flag=0 on `sched_stage` and status=`staged` on `sched_task`, set the flag=1 (staged) on `sched_stage`

## Queuing
### Steps
1. For tasks with flag=1 and status=`staged`, set the flag=2 (locked for queuing) on `sched_stage`
2. For tasks with flag=2,
    a. If task is callable, put an entry into `sched_queue`
    b. If task is bundle, put all of its callable tasks into `sched_queue`
3. For tasks in `sched_queue` with status!=`queued`, update task status to `queued`
4. For tasks with flag=2 and status=`queued`, set the flag=3 (queued) on `sched_stage`

## Splitting
### Steps
1. Lock tasks in `sched_queue` for which flag=3 and status=`queued`
2. Check start condition for locked task
    a. If satisfied, set the status=`ready` on `sched_task`
    b. If not satisfied, set the status=`waiting` on `sched_task`
3. Put an entry for locked task with status=`ready` into `sched_ready`
4. Put an entry for locked task with status=`waiting` into `sched_wait`
5. Remove a task from `sched_queue` when status is either `ready` or `waiting`