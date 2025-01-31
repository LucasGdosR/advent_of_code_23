![Data-flow graph](https://github.com/LucasGdosR/advent_of_code_23/blob/main/08/08.jpg)

Parsing the first line and the other lines can be done in parallel as long as you know in advance how long the first line is. As the problem uses a hash table, there are certain criteria that must be met in order to parallelize the structure's construction. I considered it should be done single-threaded in the diagram. However, if would be possible to build it in parallel with one of the following conditions:
- A perfect hash function with no collisions and no resizing;
- A concurrent hash table (syncMap, map with a lock);
- Buckets, with each bucket requiring locking, but threads rarelly accessing the same bucket.

We cannot guarantee the first case without knowing the input in advance. The second should be really slow. The third is the most reasonable to implement, but it's pointless for this particular case.

Once we have the "Left, Right" commands from the first line and the graph from the other lines (the hash table is a graph), we can count steps to the exit given an entry point. There's a dependendy between each move in the graph, so this must be done sequentially for each entry. However, each entry can be explored in parallel, so that's where the concurrency in this problem comes from. The nature of the input of the problem is such that we must do the product of each path lengths divided by length of the sequence of commands (their least common multiple). This multiplication can be accumulated in parallel.

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