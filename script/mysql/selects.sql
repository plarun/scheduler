-- Query tasks and its schedule details
Select 
	t.id, 
    t.name, 
    t.type,
    run_flag,
    t.run_days_bit, 
    t.start_window, 
    t.end_window,
    (select group_concat(start_time order by start_time asc separator ',') from sched_batch_run b where b.task_id=t.id) start_times,
    (select group_concat(start_min order by start_min asc separator ',') from sched_window_run w where w.task_id=t.id) start_mins
From sched_task t;

