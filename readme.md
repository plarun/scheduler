Task scheduler can schedule tasks with dependencies for execution. It also allows manual actions on tasks.

---

# Scheduler Architecture

![scheduler arch](/img/sched_arch.png)

---

## Services
### 1. Client
* CLI tool to interact with the `event server`
* Below are the available commands to action an event or view details about tasks

| Command | Usage | Example |
|:--- |:--- |:--- |
| `schd_def` | Validate task definition action file | schd_def -c -f \<filename\> |
| `schd_def` | Validate and implement task definition action file | schd_def -f \<filename\> |
| `schd_event` | Send an event on existing task | schd_event -j \<taskname\> -e \<eventname\> |
| `schd_task` | View definition of existing task | schd_task -j \<taskname\> |
| `schd_runs` | View run details of existing task | schd_runs -j \<taskname\> -c \<last_n_runs\> -d \<runs_on_specific_date DD\/MM\/YYYY\> |
| `schd_status` | View latest run status of existing task | schd_status -j \<taskname\> |

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

---

"command"
	FIELD_CONDITION    Field = "condition"
	FIELD_ERR_LOG_FILE Field = "err_log_file"
	FIELD_LABEL        Field = "label"
	FIELD_MACHINE      Field = "machine"
	FIELD_OUT_LOG_FILE Field = "out_log_file"
	FIELD_PARENT       Field = "parent"
	FIELD_PRIORITY     Field = "priority"
	FIELD_PROFILE      Field = "profile"
	FIELD_TYPE         Field = "type"
	FIELD_RUN_DAYS     Field = "run_days"
	FIELD_RUN_WINDOW   Field = "run_window"
	FIELD_START_MINS   Field = "start_mins"
	FIELD_START_TIMES  Field = "start_times"

## Task Definition
Below are the actions allowed on task

| Key | Value | Desc |
|:--- |:--- |:--- |
| `insert_task` | task name | create a new task definition |
| `update_task` | task name | update attributes of existing task definition |
| `delete_task` | task name | delete an existing task definition |

Below are the attributes allowed on task action
### 1. `type` attribute
- Task `type` can be either one of below two

| Value | Desc |
|:--- |:--- |
| `callable` | callable type task is an executable task |
| `bundle` | bundle type task is a container to keep one of more callable tasks |

| Allowed action | Mandatory |
| --- | --- |
| `insert_task` | yes |

### 2. `command` attribute
- Task `command` contains the script to be executed or command to be executed.

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | yes |

### 3. `condition` attribute
- Task's start condition is a boolean expression of tasks wrapped with status.
- Start condition of task will be evaluated when scheduled to run
- If the evaluated value is true then condition is satisfied so task can run else it has to wait till the condition get satisfied

| Example | Explaination |
| --- | --- |
| `su(task1)` | task1 should be in success state |
| `fa(task1)` | task1 should be in failure state |
| `nr(task1)` | task1 should not be in running state |
| `su(task1) & su(task2)` | task1 and task2 should be in success state |
| `su(task1) \| fa(task2)` | either task1 in success state or task2 in failure state |
| `su(task1) & nr(task2)` | task1 in success state and task2 should not be in running state |
| `su(task1) | (((su(task2) & su(task3) | fa(task4)) | nr(task5)) | nr(task6)) & su(task7)` | nested boolean logic |

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | no |