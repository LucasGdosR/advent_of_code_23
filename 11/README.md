**How to apply concurrency to this problem**

There are a bunch of tasks in this problem, so we should use a thread pool and pass them tasks.

*Task 1 (or maybe 1+2)*

The input is a single grid. We must:
- Find all `'#'` and store their coordinates
- Find all rows filled with `'.'`
- Find all columns filled with `'.'`

This can all be done in parallel using idempotent operations. If we prefill an array with `true`, calling all rows and columns empty, and then overwrite them with `false` whenever a `'#'` is found, each byte in the input can be processed independently, as long as we translate the byte position to an `i, j` coordinate. This's one approach.

When thinking about contention, the previous approach may be suboptimal. We can either split the work among threads by rows or columns. Rows lead to better cache locality, so that's the better way. However, all threads now look for `'#'` for all columns, and they all write to all indices of the empty columns array. This leads to contention when both want to write to the same index, and also leads to false-sharing when indices are close.

This led me to try a different approach. Have a job for scanning all bytes to find all `'#'`. Divide it by rows for cache locality. Each thread is the only one accessing its rows, so there's minimal contention (only potential false-sharing near edges, but this is negligible in this case). Have a separate job for scanning for empty cols. Scanning cols should be slower than rows because of cache locality, but we can fight that tendency by early breaks whenever a `'#'` is found, instead of scanning the entire column. We can also overlap appending all `'#'` coordinates with reading columns. As a side note, having a linked-list of slices would lead to a significantly better implementation than a single dynamic array that appends all elements, but this challenge is about multi-threading, not about data structures.

*Task 2*

Each `'#'` coordinate must be corrected given the empty rows and cols. We must synchronize before we start. This can be done independently for each coordinate, so each thread works on a slice of the array. This can lead to false-sharing too, but it's also negligible. When I say it is negligible, it's because the only way two threads can share the same cache line is if one of them is at the end of the work, and the other is at the beginning. That's not a likely event.

*Task 3*

We need to sum distances, and that requires all coordinates to be corrected, so we must synchronize before starting this task. We must calculate the distance between each pair of `'#'`. This is a little bit tricky, as the first `'#'` takes a bunch of work, and the last takes no work (as there's no pair left). This means evenly sharing work among threads requires us to solve an equation, which I did before coding. This is its solution:

`end = GALAXY_COUNT - int(math.Sqrt(float64(GALAXY_COUNT*GALAXY_COUNT-2*GALAXY_COUNT*(start+1)+start*start+2*start-2*step+1)))`

With this, each thread receives a range of all `'#'` that takes roughly the same work to process. The main thread them sums all partial results.

1. Producing

The main thread initializes shared variables with the workers, synchronizes when it's needed, and shares work evenly by sending equivalent slices of work to each worker.

2. Consuming

Each worker identifies which task it received and executes it. The accumulated results are sent back to the main thread either via a channel or via shared arrays.

3. Merging

This might require: 1. appending coordinate arrays together 2. doing nothing 3. summing partial sums.

4. Benchmark

The multithreaded version runs with 4 threads. Using the approach that has one job for columns and another for rows led to a 1,5x speedup. The approach that reads everything in a single pass led to a 1,6x speedup. In this case, fewer memory accesses trumped over lower contention.