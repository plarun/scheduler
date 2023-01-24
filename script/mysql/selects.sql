-- To view task's run details
Select 
	j.job_id, 
    j.job_name, 
    j.job_type,
    run_flag,
    j.run_days_bit, 
    j.start_window, 
    j.end_window,
    (select group_concat(start_time order by start_time asc separator ',') from sched_batch_run b where b.job_id=j.job_id) start_times,
    (select group_concat(start_min order by start_min asc separator ',') from sched_window_run w where w.job_id=j.job_id) start_mins
From sched_job j;