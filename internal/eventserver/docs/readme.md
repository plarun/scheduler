# Event Server
* App server to listen on services `client` and `worker`
* Routes the task action definition request to `validator` for syntax and conflict checks
* Except `allocator` all other services are depend on `event server`