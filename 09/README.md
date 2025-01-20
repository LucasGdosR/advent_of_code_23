**How to apply concurrency to this problem**

Each line can be parsed in parallel and solved in parallel, so the input is divided up in chunks and each worker solves all the lines in their chunk. The main thread aggregates the results.

1. Producing

The main thread sends an offset into the memory mapped file for each worker. This is done via the `common` library at the repository's root with `SolveCommonCaseMmapLinesInt`.

2. Consuming

Each worker solves each line and accumulates the results just as the single threaded solution. The accumulated result is sent back to the main thread via a channel.

3. Merging

Just sum all partial results.

4. Benchmark

The multithreaded version with 4 threads is 10% faster on average. It can be up to 50% faster, but it also has a higher tail execution time.