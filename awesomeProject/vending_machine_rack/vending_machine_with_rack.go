package vending_machine_rack

type VendingMachine struct {
	Racks                 []*Rack
	rackMap               map[int]*Rack
	currentState          State
	moneyInsertedState    *MoneyInsertedState
	productSelectedState  *ProductSelectedState
	productDispensedState *ProductDispensedState
	selectedRack          *Rack
	insertedMoney         float64
	returnedChange        float64
}

var vendingMachineInstance *VendingMachine

func NewVendingMachine() *VendingMachine {
	if vendingMachineInstance == nil {
		vendingMachineInstance = &VendingMachine{Racks: make([]*Rack, 0), rackMap: make(map[int]*Rack)}
		vendingMachineInstance.moneyInsertedState = &MoneyInsertedState{vendingMachineInstance}
		vendingMachineInstance.productSelectedState = &ProductSelectedState{vendingMachineInstance}
		vendingMachineInstance.productDispensedState = &ProductDispensedState{vendingMachineInstance}
		vendingMachineInstance.currentState = &MoneyInsertedState{vendingMachine: vendingMachineInstance}
	}
	return vendingMachineInstance
}

func (v *VendingMachine) AddRack(rack *Rack) {
	v.Racks = append(v.Racks, rack)
	v.rackMap[rack.ID] = rack
}

func (v *VendingMachine) InsertMoney(amount float64) {
	v.currentState.InsertMoney(amount)
}

func (v *VendingMachine) SelectRack(rackId int) {
	v.currentState.SelectRack(rackId)
}

func (v *VendingMachine) ReturnChange() {
	v.currentState.ReturnChange()
}

func (v *VendingMachine) DispenseProduct() {
	v.currentState.DispenseProduct()
}

func (v *VendingMachine) UpdateState(newState State) {
	v.currentState = newState
}
