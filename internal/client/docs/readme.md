# Scheduler Client

- Client allows the user to interact with event server.
- It provides multiple cli commands for the user to interact with tasks.

---

## CLI commands

### 1. schd_event
```txt
Usage:
    schd_event EVENT TASK

Send an event to task.

TASK:
    -j, --task=string   task name
EVENT:
    -e, --event=EVENT   event name should be one of following
                        start - starts the task
                        abort - stops the task
                        freeze - change the status of task to FROZEN
                        reset - change the status of task to IDLE
                        green - change the status of task to SUCCESS
                        red - change the status of task to FAILURE
```

Example:
```txt
# start a task test_1_task_1
$ schd_event -j test_1_task_1 -e start
```

### 2. schd_task
```txt
Usage:
    schd_task TASK

Print task definition.

TASK:
    -j, --task=string   task name
```

Example:
```txt
# view definition of task test_1_task_1
$ schd_task -j test_1_task_1

Task Name: test_1_task_1
Type: callable
Command: /opt/work/scripts/test.sh 5
OutLogFile: /opt/work/logs/test_1_task_1.out
ErrLogFile: /opt/work/logs/test_1_task_1.err
Label: callable batch task
RunDays: mo,we,th,fr
StartTimes: 08:13
Priority: 0

```

### 3. schd_status
```txt
Usage:
    schd_status TASK_NAME

Display current run and status of the task.

TASK:
    -j, --task=string   task name
```

Example
```txt
# view latest run status of callable task test_1_task_1
$ schd_status -j test_1_task_1

Task Name                                                         Start Time        End Time          Status    
_________________________________________________________________ _________________ _________________ __________
test_1_task_1                                                     20230331 15:14:59 20230331 15:33:05 success   

# view latest run status of bundle task test_5_box_1
$ schd_status -j test_5_box_1

Task Name                                                         Start Time        End Time          Status    
_________________________________________________________________ _________________ _________________ __________
test_5_box_1                                                      20230325 16:06:24 20230325 16:06:53 success   
 test_5_task_1                                                    20230325 16:06:29 20230325 16:06:30 success   
 test_5_task_2                                                    20230325 16:06:29 20230325 16:06:31 success   
 test_5_task_3                                                    20230325 16:06:33 20230325 16:06:34 success   
 test_5_task_4                                                    20230325 16:06:38 20230325 16:06:40 success   
 test_5_task_5                                                    20230325 16:06:43 20230325 16:06:46 success   
 test_5_task_6                                                    20230325 16:06:52 20230325 16:06:53 success   

```

### 4. schd_runs
```txt
Usage:
    schd_runs [OPTION]... TASK

Display previous runs and status of the task.

OPTION:
    -c, --count=NUM     number of last runs
    -d, --date=strings  only runs of given date if any
TASK:
    -j, --task=string    task name
```

Example
```txt
# view last run for task task test_1_task_1
$ schd_runs -j test_1_task_1

Task Name                                                         Start Time        End Time          Status    
_________________________________________________________________ _________________ _________________ __________
test_1_task_1                                                     20230331 15:14:59 20230331 15:15:04 success   

# view last 3 runs for task test_1_task_1
$ schd_runs -j test_1_task_1 -c 3

Task Name                                                         Start Time        End Time          Status    
_________________________________________________________________ _________________ _________________ __________
test_1_task_1                                                     20230331 15:14:59 20230331 15:15:04 success   
test_1_task_1                                                     20230331 11:19:41 20230331 11:19:45 success   
test_1_task_1                                                     20230331 11:16:22 20230331 11:16:27 success   

# view last 2 runs on 31-Mar-2023 for task test_1_task_1
$ schd_runs -j test_1_task_1 -c 2 -d "31/03/23"

Task Name                                                         Start Time        End Time          Status    
_________________________________________________________________ _________________ _________________ __________
test_1_task_1                                                     20230331 15:14:59 20230331 15:15:04 success   
test_1_task_1                                                     20230331 11:19:41 20230331 11:19:45 success   

```

### 5. schd_def
```txt
Usage:
    schd_def [OPTION] FILE

Check and process the task actions in the file.

OPTION:
    -c, --only-check    dont process the file
FILE:
    -f, --file          input file containing task definition
                        and action for one or more tasks
```

Example
```txt
# validate task definition actions at file "internal/client/etc/test/def/def1.txt"
$ schd_def -c -f internal/client/etc/test/def/def1.txt

# validate and implement task definition actions at file "internal/client/etc/test/def/def1.txt"
$ schd_def -f internal/client/etc/test/def/def1.txt

```
