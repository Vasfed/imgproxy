package main

import(
  "runtime"
  "sync/atomic"
)

type processStats struct {
  Start, Success, Fail uint64
  Processing int32
}

type processStatsCounter struct {
  data processStats
}

func (s *processStatsCounter) Start(){
  atomic.AddUint64(&s.data.Start, 1)
  atomic.AddInt32(&s.data.Processing, 1)
}
func (s *processStatsCounter) Success(){
  atomic.AddUint64(&s.data.Success, 1)
  atomic.AddInt32(&s.data.Processing, -1)
}
func (s *processStatsCounter) Fail(){
  atomic.AddUint64(&s.data.Fail, 1)
  atomic.AddInt32(&s.data.Processing, -1)
}
func (s *processStatsCounter) Discard(){
  atomic.AddUint64(&s.data.Start, ^uint64(0)) // decrement by 1
  atomic.AddInt32(&s.data.Processing, -1)
}
func (s *processStatsCounter) Export() *processStats {
  return &processStats{
    Start: atomic.LoadUint64(&s.data.Start),
    Success: atomic.LoadUint64(&s.data.Success),
    Fail: atomic.LoadUint64(&s.data.Fail),
    Processing: atomic.LoadInt32(&s.data.Processing),
  }
}

type serverStats struct {
	// currentConnections uint32
	// totalAccepted      uint64 //?
	// totalHandled       uint64 //?
	request, download, process, result processStatsCounter
}

type serverStatsResult struct {
  Requests, Download, Process, Result *processStats
}

var stats = serverStats{}


func (s *serverStats) Read() (*serverStatsResult){
  return &serverStatsResult{
    Requests: s.request.Export(),
    Download: s.download.Export(),
    Process: s.process.Export(),
    Result: s.result.Export(),
  }
}

type memStatsResult struct {
  Objects uint64
  NumGoroutines int

  Alloc, TotalAlloc,
  Sys, HeapSys, HeapInuse uint64

  LastGC        uint64
  PauseTotalNs  uint64
  GCCPUFraction float64
}

func (s* serverStats) ReadMem() (*memStatsResult) {
  var m runtime.MemStats
	runtime.ReadMemStats(&m)

  return &memStatsResult{
    Objects: m.Mallocs - m.Frees,
    NumGoroutines: runtime.NumGoroutine(),
  	Alloc: m.Alloc,
  	TotalAlloc: m.TotalAlloc,
  	Sys: m.Sys,
  	HeapSys: m.HeapSys,
  	HeapInuse: m.HeapInuse,

  	LastGC: m.LastGC,
  	PauseTotalNs: m.PauseTotalNs,
  	GCCPUFraction: m.GCCPUFraction,
  }
}
