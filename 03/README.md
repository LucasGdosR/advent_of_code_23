**How to apply concurrency to this problem**

1. Producing
- Passing every line of the file as a string through a channel would suck
- Passing indices of lines does not work, as we only know the line breaks by reading through the file
- **Solution**: since the file is a square grid, the work can be split up by lines. Each thread gets its fair share of the lines, leading to threads working on contiguous blocks of memory
- **Challenges?**: unlike memory mapped files, we know exactly where lines break, so there are no boundary edge cases for us to deal with

2. Consuming
- Passing the output of each line would lead to a lot of contention on the channel, as each worker would try to access it frequently
- Accumulating the result in a shared variable would either require locks (single variable), or lead to false-sharing (an array of variables would share a 64 bytes cache line among cores)
- **Solution**: accumulating the result to a thread-local variable leads to a lock-free implementation that does not share variables. At the end of the work, return the whole batch through a shared channel (once per batch instead of once per line)
- **Challenges**: keeping state local to each thread (buffer that accumulates "parts")

3. Merging

The main thread must reduce all partial results into a single final result. This is trivial, since solving each line is commutative.

4. Benchmark

The multithreaded solution using 4 threads was actually slower than the single threaded implementation. I blame it on the small input, a 140x140 grid. Solving the problem took on average 300k cycles, which seems to be too little to take advantage of multithreading.