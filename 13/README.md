**How to apply concurrency to this problem**

Each pattern in the file can be parsed and solved independently. This is a little bit tricky, as each pattern has a variable length, and its end is flagged by two line breaks in succession. Producing the file indices for the patterns is the main part to do this concurrently. The next part is parsing the patterns with a less friendly API, manipulating a one-dimensional byte array directly. Once that's done, it's just a matter of solving each pattern and accumulating the results.

1. Producing

The main thread sends a range of indices to be processed by each worker.

2. Consuming

Each worker adjusts the ranges to get whole patterns, solves each pattern, and accumulates the results. The accumulated result is sent back to the main thread via a channel.

3. Merging

Just sum all partial results.

4. Benchmark

The multithreaded version runs with 4 threads. It results in a 2x speedup on average.