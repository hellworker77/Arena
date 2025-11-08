package runtime

type Inventory struct {
	Slots []InventorySlot
}

type InventorySlot struct {
	ItemID string
	Count  int
}
