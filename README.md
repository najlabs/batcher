# batcher
build comprehensive batch application

## create a new Batch
```go
type Batcher interface {
    Start(ctx JobContext)
    Schedule(ctx JobContext)
}
type Job interface {
    Start(ctx JobContext)
}
type ItemReader interface {
    Read() chan interface{}
}
// batcher.newBatch()
b := batcher.newJob()
b.withName("checker")
b.withItemReader(csvreader.Read("path to file"))
b.withItemReader(jsonreader.Read("path to file"))
b.withReader(itemreader.newCsvReader("path to file"))
b.withProcessor()
b.withItemWriter()

batch := batcher.newBatch() // newScheduledBatch
batch.StartsWith(job1)
batch.Then(job1)
batch.Then(job1)
batch.EndsWith(job1)
// Hooks
batch.BeforeAll()
batch.AfterAll()
batch.BeforeEach()
batch.AfterEach()
// batch.Schedule("")
batch.Start()
```