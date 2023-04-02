# Validator
* Creating new task or updating fields of existing tasks has to be done via a specific task action syntax
* Validates the syntax and checks for conflicts as a result of implementing the actions mentioned in task action

---

## Task Definition fields
Below are the actions allowed on task

[Link to sample task definition files](internal/client/etc/test/def/)

| Key | Value | Desc |
|:--- |:--- |:--- |
| `insert_task` | task name | create a new task definition |
| `update_task` | task name | update fields of existing task definition |
| `delete_task` | task name | delete an existing task definition |

Below are the fields allowed on task action
### 1. `type` field
- Task `type` can be either one of below two

| Value | Desc |
|:--- |:--- |
| `callable` | callable type task is an executable task |
| `bundle` | bundle type task is a container to keep one of more callable tasks |

| Allowed action | Mandatory |
| --- | --- |
| `insert_task` | yes |

### 2. `command` field
- Task `command` contains the script to be executed or command to be executed.

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | yes |

### 3. `condition` field
- Task's start condition is a boolean expression of tasks wrapped with status.
- Start condition of task will be evaluated when scheduled to run
- If the evaluated value is true then condition is satisfied so task can run else it has to wait till the condition get satisfied
- condition is built with one or more clause
- clause is nothing but task name wrapped with paranthesis with prefix of task state (su/fa/nr). That is `su(task1)` or `fa(task1)` or `nr(task1)`
- `su` is success, `fa` is failure, `nr` is not running
- That is the clause `su(task1)` represents that task1 should be in success state
- Multiple clauses can be joined with either `|` or `&` operator forming a boolean expression

| Example | Desc |
| --- | --- |
| `su(task1)` | task1 should be in success state |
| `fa(task1)` | task1 should be in failure state |
| `nr(task1)` | task1 should not be in running state |
| `su(task1) & su(task2)` | task1 and task2 should be in success state |
| `su(task1) \| fa(task2)` | either task1 in success state or task2 in failure state |
| `su(task1) & nr(task2)` | task1 in success state and task2 should not be in running state |
| `su(task1) \| (((su(task2) & su(task3) \| fa(task4)) \| nr(task5)) \| nr(task6)) & su(task7)` | nested boolean logic |

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | no |

### 4. `err_log_file` field
- File to write error log of executed command

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | yes |

### 5. `label` field
- Description about the task

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable`, `bundle` | no |

### 6. `machine` field
- Run machine (feature not yet implemented)
- If null then run machine is the machine where `allocator` is running

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable`, `bundle` | no |

### 7. `out_log_file` field
- File to write output log of executed command

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | yes |

### 8. `parent` field
- If task is inside a bundle task, then name of that bundle task is `parent`

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | no |

### 9. `priority` field
- priority of the task

| Value | Desc |
| --- | --- |
| `0` or `low` | Low priority is the lowest level |
| `1` or `normal` | Normal priority |
| `2` or `important` | Important priority |
| `3` or `critical` | Critical priority is the highest level |

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable`, `bundle` | no |

### 10. `profile` field
- `profile` is an executable file with exported variables for env setup to initate before executing the `command`

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable` | no |

### 11. `run_days` field
- Days of week in comma separated value at when task should run
- Following are the decodes for weeks `su, mo, tu, we, th, fr, sa`
- If value is `all` or null then all days of week will be set
- Repeated values are not allowed

| Example | Desc |
| --- | --- |
| `all` or null value | all days of week |
| `su` | sunday |
| `su,mo,we` | sunday, monday, wednesday |

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable`, `bundle` | no |

### 12. `run_window` field
- Used to represent run window where task will be scheduled on mentioned `start_mins`
- It should be used along with field `start_mins`
- It is represented as two time (hh24:mm) values separated by hypen (-)

Examples: `00:00-03:00`, `04:00-16:00`, `12:59-23:59`

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable`, `bundle` | yes if `start_mins` used |

### 13. `start_mins` field
- Used to represent start minutes where task will be scheduled on every mentioned minute in `run_window`
- It should be used along with field `run_window`
- It is represented as comma separated value of minutes (00 to 59)

Examples: `00`, `30`, `00,20`, `15,30,45`

Task with below settings will be scheduled to run for every 15th, 30th and 45th minute between 4 AM and 8 AM
```txt
run_window: 04:00-08:00
start_mins: 15,30,45
```

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable`, `bundle` | yes if `run_window` used |

### 14. `start_times` field
- Used to represent start times where task will be scheduled
- It should not be used with `start_mins` or `run_window`
- It is represented as comma separated value of time of format (hh24:mm)

Examples: `00:00`, `23:00`, `00:20,00:40`, `04:00,16:00`

| Allowed action | Allowed task type | Mandatory |
| --- | --- | --- |
| `insert_task`, `update_task` | `callable`, `bundle` | no |
