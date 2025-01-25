**How to apply concurrency to this problem**

Each line in this file can be parsed and processed in parallel. Processing a line is done via a recursive function that splits its execution whenever a `'?'` shows up in the string. This is a tree structure. As such, it cannot be known in advance. Still, each branch on a split could be processed in parallel, theoretically. That's an awful idea in practice, as we'll see.

This problem requires memoization to run in a reasonable ammount of time. We can either hold a shared cache across threads, which might be reasonable when the cache is pretty much mostly read, or hold an exclusive cache per thread, which does not benefit from sharing previous work across threads, but does not require synchronization. In Go, this means using a `map` for each thread, or a `sync.Map` across all threads. A middle of the road approach could be achieved with partitions. I tried both setups and benchmarked their results.

1. Producing

The main thread sends a range of lines to be processed by each worker.

2. Consuming

Each worker solves each line and accumulates the results just as the single threaded solution. The accumulated result is sent back to the main thread via a channel.

3. Merging

Just sum all partial results.

4. Benchmark

The multithreaded version runs with 4 threads.

For the shared cache with `sync.Map`, the multithreaded version is 5x slower on average. There's a ton of synchronization work, as the cache ends up with 260,417 elements, meaning it was locked a bunch of times.

When each thread has its own local cache, each ends up with 64k to 70k elements, meaning there would be very little sharing between threads. Effectively, the memoization is only truly useful for solving each line, instead of having a line help in solving future lines. Having a unified cache has synchronization overheads, while giving no computation benefits on this particular input for this particular problem.

Having a local cache for each thread led to a 3,2x speedup over the single-threaded version, and a whopping 16x speedup vs the `sync.Map` version.