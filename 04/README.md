![Data-flow graph](https://github.com/LucasGdosR/advent_of_code_23/blob/main/04/04.jpg)

Each line can be parsed in parallel, and in fact each number can be parsed in parallel. As winner numbers are parsed, they be incrementally added to a set, marked by the dashed pipeline. As scratchcard numbers are parsed, they can be incrementally tested for membership in the winner set, but they require the winner set to be complete, so this is a synchronization point. As the number of winners per scratchcard are discovered, the points they score can be summed in a pipeline for part 1. For part 2, the winnings of a card depend on the winnings of the following cards. With this dependency, no parallelism is possible. We must count the winnings of the last scratchcard, then the previous, and so on. This is a dynamic programming problem that does not lend itself well to parallel approaches.

**How to apply concurrency to this problem**

1. Producing

- **Solution**: the file has a perfectly regular structure. Every line has exactly the same ammount of bytes, and every section of the information is at the same byte range. This leads itself to a very neat solution using a memory map. The producer only needs to pass the start of a line as a pointer to each worker.
- **Challenges**: the second part of the problem leads itself to dynamic programming. However, dynamic programming imposes a dependency chain on the computed values. Multithreading does not mix well with dependencies, so this part of the problem could not be entirely multi-threaded.

2. Consuming
- **Solution**: accumulating the result to a thread-local variable leads to a lock-free implementation that does not share variables. At the end of the work, return the whole batch through a shared channel (once per batch instead of once per line). `matches [][]int` were cached in an array for later single-threaded processing via dynamic programming.
- **Challenges**: as there are roughly 53 lines per worker, a one dimensional array would lead to false-sharing between threads, especially if it was implemented as a `[]byte`, so I opted for a two dimensional array instead.

3. Merging

The main thread must reduce all partial results into a single final result. This isn't as trivial for the second part, as the two dimensional array led to some edge cases that needed handling, which I solved with `mutI, mutJ`.

4. Benchmark

The multithreaded solution using 4 threads led to an overall 1,33x speedup. The input is small, and not all work could be multi-threaded. However, this is still a win compared to day 3, which had a slowdown.