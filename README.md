# Go Job Queue
A Job Queue that will serially process Jobs in the order they were pushed into the queue.

## Usage
- First, create a new Queue with a buffer size using the make function. If you do not use a buffer, Adding Jobs will block until the Queue is empty, rendering it rather useless.
- You can Add jobs to it at any time using AddJob.
- Whenever you are ready to start executing Jobs, call StartWorking. You can still add jobs afterwards.
- When you're done with the Queue, you can call Close on it. This will close the underlying channel.

## Exported functions and types

- ```type Queue chan job```
Can be created like a channel, using make with a buffer size:
  	```q := make(Queue, 20) // q is a Queue that can hold a maximum of 20 pending jobs at a time.```
  	
- ```func (q Queue) AddJob(transact func() error, rollback func(error)) (err error) ```
Adds a Job to the Queue it is called on. Takes two arguments:
	```transact func() error // The function that should contain a Job's logic.```
	```rollback func(error) // The function will only be called if transact() returns an error, with this error as an argument. Can be used for rolling back changes, doing cleanups, error logging etc.```
Returns an error if the Queue is already closed (so unlike channels, it will not panic).

- ```func (q Queue) Close()```
Closes the underlying channel.

- ```func (q Queue) StartWorking()```
Starts executing the jobs already in the Queue, and any new Jobs you add to it.