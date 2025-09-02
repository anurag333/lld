package vending_machine

type VendingMachine struct {
	Inventory             *Inventory
	currentState          State
	moneyInsertedState    *MoneyInsertedState
	productSelectedState  *ProductSelectedState
	productDispensedState *ProductDispensedState
	selectedProduct       *Product
	insertedMoney         float64
	returnedChange        float64
}

var vendingMachineInstance *VendingMachine

func NewVendingMachine() *VendingMachine {
	if vendingMachineInstance == nil {
		vendingMachineInstance = &VendingMachine{Inventory: NewInventory()}
		vendingMachineInstance.moneyInsertedState = &MoneyInsertedState{vendingMachineInstance}
		vendingMachineInstance.productSelectedState = &ProductSelectedState{vendingMachineInstance}
		vendingMachineInstance.productDispensedState = &ProductDispensedState{vendingMachineInstance}
		vendingMachineInstance.currentState = &MoneyInsertedState{vendingMachine: vendingMachineInstance}
	}
	return vendingMachineInstance
}

func (v *VendingMachine) InsertMoney(amount float64) {
	v.currentState.InsertMoney(amount)
}

func (v *VendingMachine) SelectProduct(product *Product) {
	v.currentState.SelectProduct(product)
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
