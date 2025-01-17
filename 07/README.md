**How to apply concurrency to this problem**

Each line in the file is an entry. This means we can parse lines independently. Not all lines have the same lenght, so in order to parse them in parallel using offsets into `Mmap`, some backtracking might be needed to find line breaks.

After parsing, the next step is sorting all hands. A hand's rank depends on all other hands, so sorting cannot be done fully in parallel. However, it can be done partially in parallel. Each worker sorts the hands it parsed, and so we have some slices which are already sorted. Merging sorted slices is cheaper than sorting a single unsorted array. The final merge cannot happen in parallel, but all others can.

In order to calculate the winnings, it is necessary to have the final rank of all hands. This means we must sync after sorting, just before counting winnings. Each hand's winnings can be calculated in parallel knowing the hand's final rank and bid.

The second part of the problem could be done in parallel with the first, but it would involve making a copy of all hands and bids, instead of mutating the same array. If you'd like to explore that path, tell me if it was worth it. The majority of the implementation is the same, but the sorting function is different, and there's no need to parse.

With this many parallel operations (parsing and sorting, merging sorted slices, sorting unsorted slices, counting winnings), spawning a goroutine for each would lead to quite a bit of overhead. Instead, I implemented a "worker pool". Workers are spawned in the beginning of the program and they enter a loop while a tasks channel is open. The same workers execute each part of the tasks (parsing, sorting, counting). This is implemented via a tagged union for tasks, which has an enum for what kind of task it is, and the parameters necessary for executing them. The workers also receive multiple channels when spawned in order to send back their results. Sending the channels once is better than sending them inside the task union.

1. Producing
- The main thread produces tasks for the worker threads. We want to send as little data as possible.
- **Solution**: have the task consist of a tag indicating which task to perform, a range indicated by a start and an end index, and which sorting function to use (regular or joker). This could be further reduced by creating different tags for each sorting function, and by sending int16 for the ranges.
- **Challenges**: sending only ranges in the tasks requires us to pass the arrays during worker instantiation. Both the `file` and `handsAndBids` arrays are passed. As `handsAndBids` is not initialized yet, a pointer to it must be passed.

2. Consuming
- **Solution**: consuming tasks requires the worker to identify which task to do by the tag byte, and then perform the task.
- **Challenges**: sending the tasks' output requires us to share a channel during worker instantiation. The channels could also be passed inside the task struct, but it is less overhead to pass all channels only once in the beginning.

3. Merging
- **Solution**: merging winnings is easy, it's just a sum. Merging sorted arrays efficiently requires writing a custom merge step, like merge-sort.
- **Challenges**: naively applying the default quicksort on the two concatenated slices leads to no speedup from sorting each slice. Making this step efficiently multi-threaded requires an algorithmic change.

4. Benchmark

The multithreaded version resulted in a 2,0-2,4x speedup. This problem is much larger than the previous ones, so there's enough work to get the overhead of spawning new threads be worth it. The worker pool implementation was about 20% faster than the implementation that spawned a new goroutine for each task.