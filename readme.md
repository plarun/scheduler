# Task scheduler
- Tasks can be scheduled with dependencies with other tasks.
- Allows manual actions on tasks using cli.
- Multiple tasks can be wrapped under a bundle task.
- Task can have dependencies, by setting a prerequisite (start condition) - task will only be executed if the condition satisfied.
- Tasks can be inserted/updated/deleted via specific syntax [task action definition](internal/client/etc/test/def/).
- Task can be scheduled with batch run or window run or manual run.

---

## Scheduler Architecture

![scheduler arch](/img/sched_arch.png)

---

## Services
### 1. [Client](internal/client/docs/readme.md)
* CLI tool to interact with the `event server`

### 2. Event Server
* App server to listen on services (`client`, `worker`)
* Handles all the requests from the `client` and `worker` services
* Routes the task action definition request to `validator` for syntax and conflict checks

### 3. Validator
* Creating new task or updating attributes of existing tasks has to be done via a specific task action syntax
* Validates the syntax and checks for conflicts as a result of implementing the actions mentioned in task action

### 4. Allocator
* Schedules the tasks per in task definition for execution and make the task ready for execution

### 5. Worker
* Tasks which are ready for execution (allocated by `allocator`) will be picked and executed
* Tasks' run status will be updated to `event server`
