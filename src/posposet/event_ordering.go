package posposet

import (
	"bytes"
	"sort"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
	"github.com/Fantom-foundation/go-lachesis/src/inter/idx"
)

func (p *Poset) fareOrdering(frame idx.Frame, fiWitness hash.Event, unordered inter.Events) inter.Events {
	// sort by lamport timestamp & hash
	sort.Slice(unordered, func(i, j int) bool {
		a, b := unordered[i], unordered[j]

		if a.Lamport != b.Lamport {
			return a.Lamport < b.Lamport
		}

		return bytes.Compare(a.Hash().Bytes(), b.Hash().Bytes()) < 0
	})
	ordered := unordered

	// calculate difference between highest and lowest period
	highestLamport := ordered[len(ordered)-1].Lamport
	lowestLamport := ordered[0].Lamport
	frameLamportPeriod := idx.MaxLamport(highestLamport-lowestLamport, 1)

	// calculate difference between fiWitness's median time and previous fiWitness's consensus time (almost the same as previous median time)
	frameTimePeriod := inter.MaxTimestamp(p.GetEvent(fiWitness).MedianTime-p.LastConsensusTime, 1)

	// Calculate time ratio & time offset
	timeRatio := inter.MaxTimestamp(frameTimePeriod/inter.Timestamp(frameLamportPeriod), 1)

	lowestConsensusTime := p.LastConsensusTime + timeRatio
	timeOffset := lowestConsensusTime - inter.Timestamp(lowestLamport)*timeRatio

	// Calculate consensus timestamp of an event with highestLamport (it's always fiWitness)
	p.LastConsensusTime = inter.Timestamp(highestLamport)*timeRatio + timeOffset

	// Save new timeRatio & timeOffset to frame
	p.frames[frame].SetTimes(timeOffset, timeRatio)

	return ordered
}
