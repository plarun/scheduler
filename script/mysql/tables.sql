-- drop database
Drop Database If Exists sched_test;

-- create database
Create Database If Not Exists sched_test;

-- connect
Use sched_test;

-- sched_job contains tasks
CREATE TABLE `sched_task` (
   `id` int NOT NULL AUTO_INCREMENT,
   `parent_id` int DEFAULT NULL,
   `name` varchar(64) NOT NULL,
   `type` varchar(16) NOT NULL,
   `owner_id` int DEFAULT NULL,
   `machine` int DEFAULT NULL,
   `run_flag` varchar(16) DEFAULT NULL,
   `start_condition` varchar(2048) DEFAULT NULL,
   `command` varchar(256) DEFAULT NULL,
   `std_out_log` varchar(256) DEFAULT NULL,
   `std_err_log` varchar(256) DEFAULT NULL,
   `label` varchar(256) DEFAULT NULL,
   `profile` varchar(256) DEFAULT NULL,
   `run_days_bit` int DEFAULT NULL,
   `start_window` time DEFAULT NULL,
   `end_window` time DEFAULT NULL,
   `next_run` time DEFAULT NULL,
   `priority` int DEFAULT '0',
   `created_on` datetime DEFAULT CURRENT_TIMESTAMP,
   `created_by` varchar(32) DEFAULT NULL,
   `last_updated_on` datetime DEFAULT NULL,
   `last_updated_by` varchar(32) DEFAULT NULL,
   `last_start_time` datetime DEFAULT NULL,
   `last_end_time` datetime DEFAULT NULL,
   `current_status` varchar(16) DEFAULT 'IDLE',
   `lock_flag` int DEFAULT '0',

   PRIMARY KEY (`id`),
   KEY `job_name` (`name`)
 ) ;

-- batch run details
CREATE TABLE `sched_batch_run` (
   `task_id` int NOT NULL,
   `start_time` time DEFAULT NULL,

   KEY `task_id` (`task_id`),
   CONSTRAINT `sched_batch_run_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `sched_task` (`id`)
);

 -- window run details
CREATE TABLE `sched_window_run` (
   `task_id` int NOT NULL,
   `start_min` int DEFAULT NULL,

   KEY `task_id` (`task_id`),
   CONSTRAINT `sched_window_run_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `sched_task` (`id`)
);

-- start condition relation between tasks
CREATE TABLE `sched_task_relation` (
   `task_id` int NOT NULL,
   `cond_task_id` int NOT NULL,

   KEY `task_id_idx` (`task_id`),
   KEY `cond_task_id_idx` (`cond_task_id`)
);

-- queued tasks
CREATE TABLE `sched_queue` (
   `task_id` int NOT NULL,
   `sys_entry_date` datetime DEFAULT CURRENT_TIMESTAMP,
   `priority` int DEFAULT NULL,
   `lock_flag` int NOT NULL DEFAULT '0',

   PRIMARY KEY (`task_id`)
);

-- tasks ready for execution
CREATE TABLE `sched_ready` (
   `task_id` int NOT NULL,
   `sys_entry_date` datetime DEFAULT CURRENT_TIMESTAMP,
   `priority` int DEFAULT NULL,
   `lock_flag` int DEFAULT '0',

   PRIMARY KEY (`task_id`)
);

-- tasks tracked through execution cycle
CREATE TABLE `sched_stage` (
   `task_id` int NOT NULL,
   `sys_entry_date` datetime DEFAULT CURRENT_TIMESTAMP,
   `priority` int NOT NULL DEFAULT '0',
   `flag` int NOT NULL,
   `is_bundle` int DEFAULT NULL,
   `lock_flag` int DEFAULT NULL,

   PRIMARY KEY (`task_id`)
);

-- tasks waiting for its starting condition to satisfy
CREATE TABLE `sched_wait` (
   `task_id` int NOT NULL,
   `sys_entry_date` datetime DEFAULT CURRENT_TIMESTAMP,
   `priority` int DEFAULT NULL,

   PRIMARY KEY (`task_id`)
);

-- tasks' run history
CREATE TABLE `sched_run_history` (
   `run_id` int NOT NULL AUTO_INCREMENT,
   `task_id` int NOT NULL,
   `seq_id` int NOT NULL,
   `start_time` datetime DEFAULT NULL,
   `end_time` datetime DEFAULT NULL,
   `status` varchar(16) DEFAULT NULL,
   
   PRIMARY KEY (`run_id`),
   KEY `task_id` (`task_id`,`seq_id`),
   CONSTRAINT `sched_run_history_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `sched_task` (`id`)
);