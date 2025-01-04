**How to apply concurrency to this problem**

Refer to [day 01](https://github.com/LucasGdosR/advent_of_code_23/tree/main/01/README.md). It is exactly the same.

Benchmark:

The multithreaded solution was about 20-25% faster with 4 threads. This file is only 100 lines long, so the overhead of thread creation really makes a difference. It makes more sense to optimize the string processing in this case. It could be greatly sped up, I think.