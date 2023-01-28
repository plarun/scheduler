# Allocator Service

## Actions
### Staging
- Periodically checks and poll the tasks on its scheduled time.
### Queueing
- Push the newly staged tasks into priority queue
### Splitting
- Check tasks from priority queue and split them to wait or ready based on its start condition
### Unstaging
- Check dependency tasks to be pushed to ready from wait
- Unstage the task and also its bundle if there is one