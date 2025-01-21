**How to apply concurrency to this problem**

There's no parsing in this problem, just reading a file.

Then, performing a single traversal through a cycle cannot be done in parallel, as each step depends on the last step. Theoretically it might be possible to search the cycle from both ends at a time, but this would lead to contention when writing to a "visited" set. Getting the set size would at least make operations idempotent, eliminating race condition bugs.

Finally, when the cycle has been found, we can actually parallelize the solution. Each line can be processed independently when using a raycasting algorithm. If we were to use a flood fill algorithm, this would be trickier, leading to exactly the same considerations as finding the cycle from both ends simultaneously.

1. Producing

The main thread sends a range of lines to be processed by each worker.

2. Consuming

Each worker solves each line and accumulates the results just as the single threaded solution. The accumulated result is sent back to the main thread via a channel.

3. Merging

Just sum all partial results.

4. Benchmark

The multithreaded version with 4 threads is almost 20% slower on average. There's just not enough work for it to get over the overhead.