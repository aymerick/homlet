package homlet

import "sort"

// SortPackets sorts packets by device id
func SortPackets(packets []*Packet) {
	order := func(a, b *Packet) bool {
		return a.DeviceID < b.DeviceID
	}
	PacketsBy(order).Sort(packets)
}

// PacketsBy is a function that sorts a list of Packet
type PacketsBy func(a, b *Packet) bool

// Sort sorts Packets
func (by PacketsBy) Sort(l []*Packet) {
	sort.Sort(&PacketSorter{
		packets: l,
		by:      by,
	})
}

// PacketSorter implements sort.Interface
type PacketSorter struct {
	packets []*Packet
	by      func(a, b *Packet) bool
}

// Len implements sort.Interface
func (s *PacketSorter) Len() int {
	return len(s.packets)
}

// Swap implements sort.Interface
func (s *PacketSorter) Swap(i, j int) {
	s.packets[i], s.packets[j] = s.packets[j], s.packets[i]
}

// Less implements sort.Interface
func (s *PacketSorter) Less(i, j int) bool {
	return s.by(s.packets[i], s.packets[j])
}
