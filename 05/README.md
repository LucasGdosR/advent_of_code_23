![Data-flow graph](https://github.com/LucasGdosR/advent_of_code_23/blob/main/05/05.jpg)

Seeds, intervals, and trees can be parsed in parallel. Once a seed / interval is parsed, it can enter a pipeline to be mapped. This has two meanings: each seed / interval can be mapped in parallel (with the individual seed going through all trees), as is also possible to map all seeds through one the same tree in parallel. In this case, it is better to go through all seeds / intervals once a tree is ready. The i-th mapping requires the i-th tree to be built. Although all trees can be built in parallel, they can only be used for mapping once all previous trees have already been built. As seeds and intervals are mapped, the minimum value can be gathered via a pipeline that keeps only the minimum value seen.

**How to apply concurrency to this problem**

This problem has three parts: building interval trees, mapping seeds according to the trees, and mapping intervals according to the trees. The input file's lines have a variable length, so we must parse the file to get to the trees. This means they must be built sequentially. Fundamentally, though, the trees do not depend on each other, so they could be built in parallel. Every interval and every seed could also be mapped in parallel, though they must go through the trees in sequential order.

Measuring the cycles for the three parts yields the results that building trees takes longer than the other two combined. This means that a single goroutine can take care of all the mappings.

1. Producing
- Seeds and intervals may be processed in parallel, but they require built interval trees.
- **Solution**: the main thread does the parsing. It passes the seeds and intervals to the worker and builds the trees, which are passed through a channel. The worker can start working as soon as trees are available.
- **Challenges**: one could manually inspect the file and hardcode indices for building trees in parallel without parsing the whole file by using a memory map. If you're building a solver for a single problem, you could just hardcode the solution and not run anything, though.

2. Consuming
- **Solution**: the worker thread receives a tree and maps both seeds and intervals according to it.

3. Merging

There's no merging, as there's no branching. There's only a single worker. The work could be reduced either by the consumer or the producer. It's less channel overhead to reduce it in the consumer, but reducing both answers could be done in parallel if the producer did it. Doesn't make any practical difference, as it's blazingly fast to get the min out of the seeds array.

4. Benchmark

The whole mapping phase, including both seeds and intervals, takes only ~100k cycles. This is not worth multithreading at all. In fact, the multithreaded solution was about 1,33x slower than the single-threaded solution. If it was possible to do build all trees in parallel (which it fundamentally is, the only reason it isn't is because of the file format), it wouldn't make much of a difference, as the solution only takes about ~300k cycles, and we've seen on day 3 that that's not enough for surpassing the multithreading overhead.