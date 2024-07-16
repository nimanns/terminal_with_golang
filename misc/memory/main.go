package main

import (
	"fmt"
	"math/rand"
)

const (
	page_size         = 4096
	virtual_mem_size  = 1 << 20
	physical_mem_size = 1 << 19
)

type page struct {
	data [page_size]byte
}

type memory_manager struct {
	virtual_mem  map[uint32]*page
	physical_mem [physical_mem_size / page_size]*page
	page_table   map[uint32]int
}

func new_memory_manager() *memory_manager {
	return &memory_manager{
		virtual_mem: make(map[uint32]*page),
		page_table:  make(map[uint32]int),
	}
}

func (mm *memory_manager) allocate(virtual_addr uint32) {
	page_num := virtual_addr / page_size
	if _, exists := mm.virtual_mem[page_num]; !exists {
		new_page := &page{}
		for i := range new_page.data {
			new_page.data[i] = byte(rand.Intn(256))
		}
		mm.virtual_mem[page_num] = new_page
	}
}

func (mm *memory_manager) access_memory(virtual_addr uint32) (byte, error) {
	page_num := virtual_addr / page_size
	offset := virtual_addr % page_size

	if _, exists := mm.virtual_mem[page_num]; exists {
		if frame_num, in_physical := mm.page_table[page_num]; in_physical {
			return mm.physical_mem[frame_num].data[offset], nil
		} else {
			frame_num := mm.load_into_physical_memory(page_num)
			return mm.physical_mem[frame_num].data[offset], nil
		}
	}

	return 0, fmt.Errorf("segmentation fault: accessing unallocated memory at page %d", page_num)
}


func (mm *memory_manager) load_into_physical_memory(page_num uint32) int {
	frame_num := mm.find_free_frame()
	
	mm.physical_mem[frame_num] = mm.virtual_mem[page_num]
	mm.page_table[page_num] = frame_num
	
	return frame_num
}

func (mm *memory_manager) find_free_frame() int {
	for i, frame := range mm.physical_mem {
		if frame == nil {
			return i
		}
	}
	
	evict_frame := rand.Intn(len(mm.physical_mem))
	for vpn, fpn := range mm.page_table {
		if fpn == evict_frame {
			delete(mm.page_table, vpn)
			break
		}
	}
	return evict_frame
}

func main() {
	mm := new_memory_manager()

	for i := uint32(0); i < 10; i++ {
		mm.allocate(i * page_size)
	}

	fmt.Println("Allocated pages:")
	for page := range mm.virtual_mem {
		fmt.Printf("Page %d\n", page)
	}

	for i := 0; i < 20; i++ {
		addr := uint32(rand.Intn(10 * int(page_size)))
		value, err := mm.access_memory(addr)
		if err != nil {
			fmt.Printf("Error accessing address 0x%x: %v\n", addr, err)
		} else {
			fmt.Printf("Value at address 0x%x (page %d, offset %d): %d\n", 
				addr, addr/page_size, addr%page_size, value)
		}
	}
}
