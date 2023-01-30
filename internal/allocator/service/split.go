package service

// 1. set flag to 2 (pre queue)
// 2. if task is callable
// 2a. push it into queue
// 2b. change status to 'queued'
// 2c. set flag to 3 (post queue)
// 3. if task is bundle
// 3a. push its callable tasks into queue
// 3b. change status of callable tasks to 'queued'
// 3c. change status of bundle to 'queued'
// 3d. set flag to 3 (post queue) for bundle
// 3e. change status of bundle to 'running'
type TaskSplitter struct{}

// Split
func (t *TaskSplitter) Split() {

}

func (t *TaskSplitter) checkStartCondition() {
	// get status of distinct tasks in start condition
	// su(task1)|(su(task2)&su(task3))|(fa(task4)|nr(task5))
	/*
		(
			su(task1) | (
				su(task2) &
				su(task3)
			) | (
				fa(task4) |
				nr(task5)
			)
		)


	*/
}
