![Data-flow graph](https://github.com/LucasGdosR/advent_of_code_23/blob/main/02/02.jpg)

Each line could be parsed in parallel. The approach I took was first tokenizing, then parsing the colors. In this way, each color could be parsed in parallel. Each color must take the maximum value, so either all reds're processed by the same node, or they could be processed in a map-reduce style (many nodes get the local maximum, and one reducer node gets the global maximum from all previous nodes. This could even extend to more hierarchy levels). Alternatively, parser combinators could be used to tokenize and parse in one go. This forfeits parallelism at this level, though. It's still fine, since each line could be parsed in parallel. Each color produces either the line's id or 0. A new node must reduce each id from each color, so this is a synchronization point. Summing all the ids and powers results requires every line, so they must sync too (although this sum could be pipelined ids and powers are produced, indicated by the dashed lines). When both sums are ready, the program ends.

**How to apply concurrency to this problem**

Refer to [day 01](https://github.com/LucasGdosR/advent_of_code_23/tree/main/01/README.md). It is exactly the same.

Benchmark:

The multithreaded solution was about 20-25% faster with 4 threads. This file is only 100 lines long, so the overhead of thread creation really makes a difference. It makes more sense to optimize the string processing in this case. It could be greatly sped up, I think.
