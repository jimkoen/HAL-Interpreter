package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type halInstruction struct {
	name    string
	operand float64
}

type halState struct {
	name            string
	accumulator     float64
	programCounter  int
	registers       [16]float64
	memory          []halInstruction
	inputRegisters  map[int]chan float64
	outputRegisters map[int]chan float64
}

type opFunction func(*halState)

var opcodes = map[string]opFunction{
	"START":    opStart,
	"STOP":     opStop,
	"OUT":      opOut,
	"IN":       opIn,
	"LOAD":     opLoad,
	"LOADNUM":  opLoadNum,
	"STORE":    opStore,
	"JUMPNEG":  opJumpNeg,
	"JUMPPOS":  opJumpPos,
	"JUMPNULL": opJumpNull,
	"JUMP":     opJump,
	"ADD":      opAdd,
	"ADDNUM":   opAddNum,
	"SUB":      opSub,
	"SUBNUM":   opSubNum,
	"MUL":      opMul,
	"MULNUM":   opMulNum,
	"DIV":      opDiv,
	"DIVNUM":   opDivNum,
}

func serialize(pathname string) (instructions []halInstruction, err error) {
	fmt.Println("Serializing input...")
	file, err := os.Open(pathname)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return instructions, nil
			} else {
				return nil, err
			}
		}
		line = strings.Trim(line, "\n")
		tokens := strings.Split(line, " ")
		instruction := strings.ToUpper(tokens[1])
		var operand float64
		if len(tokens) <= 2 {
			operand = 0
		} else {
			operand, err = strconv.ParseFloat(tokens[2], 64)
		}
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, halInstruction{instruction, operand})

		if err != nil {
			return instructions, err
		}
	}

	return
}

func serializeConcurrent(pathname string) (processors []halState, err error) {
	fmt.Println("Serializing input...")
	file, err := os.Open(pathname)
	var paths []string

	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	line = strings.Trim(line, "\n")
	if line == "--processors--" {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return nil, err
				} else {
					return nil, err
				}
			}
			line = strings.Trim(line, "\n")
			if line == "--connections--" {
				break
			}
			line = strings.SplitAfter(line, " ")[1]
			paths = append(paths, line)
		}

	}
	for i, path := range paths {
		fmt.Println(i)
		instructions, err := serialize(path)
		if err != nil {
			fmt.Println("Error in serializing instructions while appening halStates (in serializeConcurrenet")
		}
		processors = append(processors, halState{
			name:            path,
			accumulator:     0,
			programCounter:  0,
			registers:       [16]float64{},
			memory:          instructions,
			inputRegisters:  map[int]chan float64{},
			outputRegisters: map[int]chan float64{},
		})
	}

	return processors, nil
}

func connect(processorA, portA, processorB, portB int, states *[]halState) {
	(*states)[processorA].outputRegisters[portA] = make(chan float64, 2)
	(*states)[processorB].inputRegisters[portB] = (*states)[processorA].outputRegisters[portA]
}

func serializeAndConnect(pathname string, processors []halState, err error) ([]halState, error) {

	file, err := os.Open(pathname)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		line = strings.Trim(line, "\n")

		if line == "--connections--" {
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						return processors, err
					} else {
						return nil, err
					}
				}
				line = strings.Trim(line, "\n")
				if line == "\n" || err == io.EOF {
					break
				}
				var connectionCombination []string = strings.Split(line, " > ")
				var connectionsA []string = strings.Split(connectionCombination[0], ":")
				var connectionsB []string = strings.Split(connectionCombination[1], ":")
				cpuIndexA, err := strconv.Atoi(connectionsA[0])
				cpuPortA, err := strconv.Atoi(connectionsA[1])
				cpuIndexB, err := strconv.Atoi(connectionsB[0])
				cpuPortB, err := strconv.Atoi(connectionsB[1])

				connect(cpuIndexA, cpuPortA, cpuIndexB, cpuPortB, &processors)

			}
			return processors, nil
			break
		}

	}
	return processors, nil
}

func halProcessor(state halState, debug bool, wg *sync.WaitGroup) {
	fmt.Println(state.name, "- routine starting...\n")
	start := time.Now()
	if debug {
		for {
			fmt.Println(state.name, "- routine ", state.programCounter, " ", state.memory[state.programCounter])
			fmt.Println("State before: ", state.registers, " ", "ACC: ", state.accumulator)
			if state.memory[state.programCounter].name == "STOP" {
				break
			}
			opcodes[state.memory[state.programCounter].name](&state)
			fmt.Println("State after:  ", state.registers, " ", "ACC: ", state.accumulator)
		}

	} else {
		for {
			if state.memory[state.programCounter].name == "STOP" {
				break
			}
			opcodes[state.memory[state.programCounter].name](&state)
		}
	}
	elapsed := time.Since(start)

	fmt.Println("HAL Program executed in", elapsed)
	wg.Done()
}

func readInput() (input chan float64) {
	input = make(chan float64, 1)
	var u float64
	_, err := fmt.Scanf("%g\n", &u)
	if err != nil {
		panic(err)
	}
	input <- u
	return input
}

func outputresult(state halState, output chan float64, wg *sync.WaitGroup) {

	for {
		select {
		case msg := <-output:
			fmt.Println("STD OUTPUT RESULT for ", state.name, ":\n", msg)
			return
		default:

		}

	}
}

func main() {

	//mem, error := serialize(os.Args[1])
	processors, error := serializeConcurrent(os.Args[1])
	if error != nil {
		fmt.Println(error)

	}
	processors, error = serializeAndConnect(os.Args[1], processors, error)
	fmt.Println(processors)
	processors[0].inputRegisters[0] = readInput()

	var waitGroupOuput sync.WaitGroup

	for i := 0; i < len(processors); i++ {
		processors[i].outputRegisters[1] = make(chan float64)
		go outputresult(processors[i], processors[i].outputRegisters[1], &waitGroupOuput)
	}

	var waitGroupHalProcessors sync.WaitGroup

	for i := 0; i < len(processors); i++ {
		waitGroupHalProcessors.Add(1)
		go halProcessor(processors[i], true, &waitGroupHalProcessors)
	}

	waitGroupHalProcessors.Wait()

	fmt.Println("main is about to exit")

	//var chFloat = make(chan float64)

	/*var mainState = halState{
		accumulator:    0,
		programCounter: 0,
		registers:      [16]float64{},
		memory:         mem,
		inputRegisters: map[int]chan float64{
			1: chFloat,
		},
	}*/

	/*if(os.Args[2] != nil){

	}*/
	/*if(strings.ToUpper(os.Args[2]) == "DEBUG"){
	fmt.Println(halProcessor(mainState, true))
	}else{
		fmt.Println(halProcessor(mainState, false))
	}*/

}
