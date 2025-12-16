package runtime

type ItemSlotType int
type Rarity int

// Item Slot
const (
	SlotType_Head ItemSlotType = iota
	SlotType_Body
	SlotType_Gloves
	SlotType_Belt
	SlotType_Boots
	SlotType_Amulet
	SlotType_Ring
	SlotType_Hand
)

const (
	Rarity_Common Rarity = iota
	Rarity_Magic
	Rarity_Rare
	Rarity_Epic
	Rarity_Legendary
	Rarity_Mythical
	Rarity_Set
	Rarity_Runeword
)

type ItemProperty struct {
	Name  string
	Value float64
}

type Item struct {
	ID         int
	Name       string
	SlotType   ItemSlotType
	Rarity     Rarity
	Properties []ItemProperty
}
