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

### 2. [Event Server](internal/eventserver/docs/readme.md)
* App server to listen on services (`client`, `worker`)

### 3. [Validator](internal/validator/docs/readme.md)
* Validates the syntax, data and conflicts of implementing the task action definition

### 4. [Allocator](internal/allocator/docs/readme.md)
* Schedules the tasks for execution and make the task ready for execution

### 5. [Worker](internal/worker/docs/readme.md)
* Tasks are executed by `worker` and will request `event server` for state update
