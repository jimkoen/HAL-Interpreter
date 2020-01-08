package main

import "fmt"

func opStart(state *halState) {
	state.programCounter++
}

func opStop(state *halState) {
	fmt.Println("Program has stopped.")
	state.programCounter++

}

func opOut(state *halState) {
	fmt.Println("Sending from ", state.name)
	state.outputRegisters[int(state.memory[state.programCounter].operand)] <- state.accumulator
	state.programCounter++
}

func opIn(state *halState) {
	fmt.Println(state.name, " is receiving a value")
	state.accumulator = <-state.inputRegisters[int(state.memory[state.programCounter].operand)]
	state.programCounter++
}

func opLoad(state *halState) {
	state.accumulator = state.registers[int(state.memory[state.programCounter].operand)]
	state.programCounter++
}

func opLoadNum(state *halState) {
	state.accumulator = state.memory[state.programCounter].operand
	state.programCounter++
}

func opStore(state *halState) {
	state.registers[int(state.memory[state.programCounter].operand)] = state.accumulator
	state.programCounter++
}

func opJumpNeg(state *halState) {
	if state.accumulator < 0 {
		state.programCounter = int(state.memory[state.programCounter].operand)
	} else {
		state.programCounter++
	}

}

func opJumpPos(state *halState) {
	if state.accumulator > 0 {
		state.programCounter = int(state.memory[state.programCounter].operand)
	} else {
		state.programCounter++
	}

}

func opJumpNull(state *halState) {
	if state.accumulator == 0 {
		state.programCounter = int(state.memory[state.programCounter].operand)
	} else {
		state.programCounter++
	}

}

func opJump(state *halState) {
	state.programCounter = int(state.memory[state.programCounter].operand)
}

func opAdd(state *halState) {
	state.accumulator = state.accumulator + state.registers[int(state.memory[state.programCounter].operand)]
	state.programCounter++
}

func opAddNum(state *halState) {
	state.accumulator = state.accumulator + state.memory[state.programCounter].operand
	state.programCounter++

}

func opSub(state *halState) {
	state.accumulator = state.accumulator - state.registers[int(state.memory[state.programCounter].operand)]
	state.programCounter++
}

func opSubNum(state *halState) {
	state.accumulator = state.accumulator - state.memory[state.programCounter].operand
	state.programCounter++
}

func opMul(state *halState) {
	state.accumulator = state.accumulator * state.registers[int(state.memory[state.programCounter].operand)]
	state.programCounter++
}

func opMulNum(state *halState) {
	state.accumulator = state.accumulator * state.memory[state.programCounter].operand
	state.programCounter++
}

func opDiv(state *halState) {
	state.accumulator = state.accumulator / state.registers[int(state.memory[state.programCounter].operand)]
	state.programCounter++
}

func opDivNum(state *halState) {
	state.accumulator = state.accumulator / state.memory[state.programCounter].operand
	state.programCounter++
}
