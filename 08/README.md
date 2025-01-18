**How to apply concurrency to this problem**

The first line of the file needs to be parsed by a single-thread, as we need to find the first line break. The rest of the file could be parsed in parallel using mmap easily, as all lines have the same length. However, the bulk of the work is not parsing, so this would be negligible.

Part 1 of the problem cannot be run in parallel. It is strictly sequential. I could not think of any possible trick. Part 2, however, fits right in, as it's just solving part 1 six times, and they can be done in parallel.

1. Producing
- **Solution**: the main thread sends a starting point for each worker.
- **Challenges**: instead of having a number of workers based on the number of available cores, this problem should spawn one worker per starting point, regardless if there are cores available. For instance, spawning 4 workers would solve 4 starting points in parallel (taking T time), and then 2 starting points in parallel (taking T time), resulting in 2T. With 6 workers, even with only 4 cores available, each worker is running roughly 2/3 of the time, so it takes 3T/2 for a solution, but they're done after 1,5T.

2. Consuming

Each worker finds how many steps it takes to get to the exit. If it's the starting point for part 1, it records that count directly. If it's for others, it sends the result for the main thread to merge.

3. Merging

Just multiply all partial results.

4. Benchmark

The multithreaded version resulted in a 2-3x speedup for 4 threads. It would still scale up to 6 threads, but no further, as that's how many starting points there are.