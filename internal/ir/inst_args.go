package ir

//go:generate stringer -type=SlotKind -trimprefix=Slot
type SlotKind uint8

const (
	SlotInvalid SlotKind = iota
	SlotCallArg
	SlotParam
	SlotTemp
	SlotUniq
	SlotDiscard
)

type Slot struct {
	ID   uint8
	Kind SlotKind
}

func NewCallArgSlot(id uint8) Slot { return Slot{ID: id, Kind: SlotCallArg} }
func NewParamSlot(id uint8) Slot   { return Slot{ID: id, Kind: SlotParam} }
func NewTempSlot(id uint8) Slot    { return Slot{ID: id, Kind: SlotTemp} }
func NewUniqSlot(id uint8) Slot    { return Slot{ID: id, Kind: SlotUniq} }

func (s Slot) ToInstArg() InstArg {
	return InstArg((uint16(s.ID) << 8) | uint16(s.Kind))
}

func (s Slot) IsInvalid() bool { return s.Kind == SlotInvalid }
func (s Slot) IsCallArg() bool { return s.Kind == SlotCallArg }
func (s Slot) IsParam() bool   { return s.Kind == SlotParam }
func (s Slot) IsTemp() bool    { return s.Kind == SlotTemp }
func (s Slot) IsUniq() bool    { return s.Kind == SlotUniq }
