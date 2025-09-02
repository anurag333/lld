package main

import vm "awesomeProject/vending_machine_rack"

func vending_machine_rack() {
	vendingMachine := vm.NewVendingMachine()

	// Add products to the inventory
	coke := vm.NewProduct(1, "Coke", 1.50)
	pepsi := vm.NewProduct(2, "Pepsi", 1.75)

	vendingMachine.AddRack(vm.NewRack(1, coke, 3))
	vendingMachine.AddRack(vm.NewRack(2, pepsi, 3))

	// Insert money into the vending machine
	vendingMachine.InsertMoney(2.00)

	// Select a product to purchase
	vendingMachine.SelectRack(1)

	// Dispense the selected product
	vendingMachine.DispenseProduct()

	// Return any remaining change
	vendingMachine.ReturnChange()

	// Insert more money to purchase another product with insufficient funds
	vendingMachine.InsertMoney(1.50)

	// Select a second product to purchase
	vendingMachine.SelectRack(2)

	// Dispense the selected product
	vendingMachine.DispenseProduct()

	// Return any remaining change
	vendingMachine.ReturnChange()

}
