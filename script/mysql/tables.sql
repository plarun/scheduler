-- drop database
Drop Database If Exists sched_test;

-- create database
Create Database If Not Exists sched_test;

-- connect
Use sched_test;

-- sched_job contains tasks
Create Table `sched_job` (
    `job_id`          int Not Null AUTO_INCREMENT,
    `parent_id`       int Default Null,
    `job_name`        varchar(64) Not Null,
    `job_type`        varchar(16) Not Null,
    `owner_id`        int Default Null,
    `machine`         int Default Null,
    `run_flag`        varchar(16) Default Null,
    `start_condition` varchar(2048) Default Null,
    `command`         varchar(256) Default Null,
    `std_out_log`     varchar(256) Default Null,
    `std_err_log`     varchar(256) Default Null,
    `label`           varchar(256) Default Null,
    `job_profile`     varchar(256) Default Null,
    `run_days_bit`    int Default Null,
    `start_window`    time Default Null,
    `end_window`      time Default Null,
    `next_run`        time Default Null,
    `priority`        int Default '0',
    `created_on`      datetime Default CURRENT_TIMESTAMP,
    `created_by`      varchar(32) Default Null,
    `last_updated_on` datetime Default Null,
    `last_updated_by` varchar(32) Default Null,
    `last_start_time` datetime Default Null,
    `last_end_time`   datetime Default Null,
    `current_status`  varchar(16) Default 'IDLE',

    Primary Key (`job_id`),
    Key         `job_name` (`job_name`)
);

-- batch run details
Create Table `sched_batch_run` (
    `job_id`     int Not Null,
    `start_time` time Default Null
);

 -- window run details
Create Table `sched_window_run` (
    `job_id`    int Not Null,
    `start_min` int Default Null
);

-- start condition relation between tasks
Create Table `sched_job_relation` (
    `job_id`      int Not Null,
    `cond_job_id` int Not Null
);


Create Table `sched_abort` (
    `job_id`         int Not Null,
    `sys_entry_date` datetime Default CURRENT_TIMESTAMP,
    `lock_flag`      int Default Null
);

Create Table `sched_exec` (
    `job_id`        int Not Null,
    `last_beat`     datetime Default Null,
    `remote_server` varchar(256) Default Null,
    `remote_fid`    varchar(256) Default Null,
    `lock_flag`     int Default Null
);

Create Table `sched_queue` (
    `job_id`         int Not Null,
    `sys_entry_date` datetime Default CURRENT_TIMESTAMP,
    `priority`       int Default Null,
    `lock_flag`      int Default Null
)

Create Table `sched_ready` (
   `job_id`         int Not Null,
   `sys_entry_date` datetime Default CURRENT_TIMESTAMP,
   `priority`       int Default Null,
   `lock_flag`      int Default Null
);

Create Table `sched_stage` (
    `job_id`         int Not Null,
    `sys_entry_date` datetime DEFAULT CURRENT_TIMESTAMP,
    `priority`       int Not Null DEFAULT '0',
    `flag`           int Not Null,
    `is_bundle`      int Default Null,
    `lock_flag`      int Default Null
);

Create Table `sched_wait` (
    `job_id`         int Not Null,
    `sys_entry_date` datetime Default CURRENT_TIMESTAMP,
    `priority`       int Default Null
);