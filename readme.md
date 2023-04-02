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