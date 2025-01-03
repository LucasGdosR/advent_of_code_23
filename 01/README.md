**How to apply concurrency to this problem**

1. Producing
- Passing every line of the file as a string through a channel would suck
- Passing indices of lines does not work, as we only know the line breaks by reading through the file
- **Solution**: map the file to memory and pass a pointer with offsets for the start and end of the worker's share
- **Challenges**: deal with boundary conditions: offsets won't coincide with line breaks, so they must be adjusted to proccess each line exactly once

2. Consuming
- Passing the output of each line would lead to a lot of contention on the channel, as each worker would try to access it frequently
- Accumulating the result in a shared variable would either require locks (single variable), or lead to false-sharing (an array of variables would share a 64 bytes cache line among cores)
- **Solution**: accumulating the result to a thread-local variable leads to a lock-free implementation that does not share variables. At the of the work, return the whole batch through a shared channel (once per batch instead of once per line)

3. Merging

The main thread must reduce all partial results into a single final result. This is trivial, since solving each line is commutative.

4. Benchmark

The multithreaded solution using 4 threads led to a 2-3x speedup for a 1000-line file. This is pretty good considering reading the file contributes a bit to the time it takes, and it's single-threaded in both implementations. I was surprised that thread creation overhead was negligible.