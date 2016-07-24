// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package cfgraph_test

import (
	"fmt"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/bblock"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/cfgraph"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/graph"
	"io/ioutil"
	"testing"
)

// verifyControlFlowGraphs is an helper method for tests to check if the generated expectedCfGraph
// is equal the actual graph.
func VerifyControlFlowGraphs(expectedCfGraph *cfgraph.ControlFlowGraph, correctCfGraph *graph.Graph) error {
	if expectedCfGraph.GetNumberOfNodes() != correctCfGraph.GetNumberOfNodes() {
		return fmt.Errorf("Number of nodes in graph should be %d, but are %d!", correctCfGraph.GetNumberOfNodes(),
			expectedCfGraph.GetNumberOfNodes())
	}
	if expectedCfGraph.GetNumberOfEdges() != correctCfGraph.GetNumberOfEdges() {
		return fmt.Errorf("Number of edges in graph should be %d, but are %d!", correctCfGraph.GetNumberOfEdges(),
			expectedCfGraph.GetNumberOfEdges())
	}

	for key, correctNode := range correctCfGraph.Nodes {
		if expectedNode, ok := expectedCfGraph.Nodes[key]; ok {
			//Node exist in graph, now check its edges.
			for index, correctEdge := range correctNode.GetOutNodes() {
				if correctNode.GetOutDegree() != expectedNode.GetOutDegree() {
					return fmt.Errorf("Node %s should have %d out-edges!\n", correctNode, correctNode.GetOutDegree())
				}
				if correctEdge.UID() != expectedNode.GetOutNodes()[index].UID() {
					return fmt.Errorf("Edge ( %s -> %s ) should not be present!", correctNode, correctEdge.Value)
				}
			}

		} else {
			return fmt.Errorf("Node %s is not found in the graph!", correctNode)
		}
	}
	return nil
}

// VerifyBasicBlocks checks the list of expected basic-blocks with the list of actual basic-blocks.
func VerifyBasicBlocks(expectedBasicBlocks []*bblock.BasicBlock, correctBasicBlocks []*bblock.BasicBlock) error {
	if len(expectedBasicBlocks) != len(correctBasicBlocks) {
		return fmt.Errorf("Number of basic-blocks should be %d, but are %d!\n", len(correctBasicBlocks), len(expectedBasicBlocks))
	}

	//Loop through all generated basic-blocks and check if they are similar to the correct once.
	for index := range expectedBasicBlocks {
		if expectedBasicBlocks[index].Type != correctBasicBlocks[index].Type {
			//Check that basic-block type is correct.
			return fmt.Errorf("Basic block nr. %d should be of type %s, but are of type %s!\n",
				index, correctBasicBlocks[index].Type.String(), expectedBasicBlocks[index].Type.String())
		}

		//Check that length of generate basic-blocks successors are equal correct number of successor blocks.
		if len(expectedBasicBlocks[index].GetSuccessorBlocks()) != len(correctBasicBlocks[index].GetSuccessorBlocks()) {
			return fmt.Errorf("Number of successors in basic-block nr. %d should be %d, and not %d!\n",
				expectedBasicBlocks[index].Number, len(correctBasicBlocks[index].GetSuccessorBlocks()),
				len(expectedBasicBlocks[index].GetSuccessorBlocks()))
		}

		//Check that basic block starts at right line.
		if expectedBasicBlocks[index].EndLine != correctBasicBlocks[index].EndLine {
			return fmt.Errorf("Basic block nr. %d should end at line number %d, and not %d!\n", expectedBasicBlocks[index].Number,
				correctBasicBlocks[index].EndLine, expectedBasicBlocks[index].EndLine)
		}

		//Check that that basic-block has correct successor blocks, and their order.
		for i, successorBlock := range expectedBasicBlocks[index].GetSuccessorBlocks() {
			if successorBlock.Number != correctBasicBlocks[index].GetSuccessorBlocks()[i].Number {
				return fmt.Errorf("Basic block nr. %d's successor block nr. %d should be nr. %d, and not %d!\n",
					index, i, correctBasicBlocks[index].GetSuccessorBlocks()[i].Number, successorBlock.Number)
			}
		}

	}
	return nil
}

func TestSimpleControlFlowGraph(t *testing.T) {
	sourceFile, err := ioutil.ReadFile("./testcode/_simple.go")
	if err != nil {
		t.Fatal(err)
	}
	basicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedGraph := cfgraph.GetControlFlowGraph(basicBlocks)
	correctGraph := graph.NewGraph()

	START := bblock.NewBasicBlock(-1, bblock.START, 0)
	EXIT := bblock.NewBasicBlock(-1, bblock.EXIT, 0)

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 10)

	BB0.AddSuccessorBlock(BB1)

	correctBasicBlocks := []*bblock.BasicBlock{BB0, BB1}

	//Test basic-blocks.
	if err := VerifyBasicBlocks(basicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}

	correctGraph.InsertEdge(&graph.Node{Value: START}, &graph.Node{Value: BB0})
	correctGraph.InsertEdge(&graph.Node{Value: BB0}, &graph.Node{Value: BB1})
	correctGraph.InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: EXIT})
	correctGraph.InsertEdge(&graph.Node{Value: EXIT}, &graph.Node{Value: START})

	if err := VerifyControlFlowGraphs(expectedGraph[0], correctGraph); err != nil {
		t.Fatal(err)
	}
}

func TestIfElseControlFlowGraph(t *testing.T) {
	sourceFile, err := ioutil.ReadFile("./testcode/_ifelse.go")
	if err != nil {
		t.Fatal(err)
	}
	basicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedGraph := cfgraph.GetControlFlowGraph(basicBlocks)
	correctGraph := graph.NewGraph()

	START := bblock.NewBasicBlock(-1, bblock.START, 0)
	EXIT := bblock.NewBasicBlock(-1, bblock.EXIT, 0)

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 6)
	BB1 := bblock.NewBasicBlock(1, bblock.IF_CONDITION, 11)
	BB2 := bblock.NewBasicBlock(2, bblock.ELSE_CONDITION, 14)
	BB3 := bblock.NewBasicBlock(3, bblock.ELSE_BODY, 18)
	BB4 := bblock.NewBasicBlock(4, bblock.RETURN_STMT, 19)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3)
	BB2.AddSuccessorBlock(BB4)
	BB3.AddSuccessorBlock(BB4)

	correctBasicBlocks := []*bblock.BasicBlock{BB0, BB1, BB2, BB3, BB4}

	//Test basic-blocks.
	if err := VerifyBasicBlocks(basicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}

	correctGraph.InsertEdge(&graph.Node{Value: START}, &graph.Node{Value: BB0})
	correctGraph.InsertEdge(&graph.Node{Value: BB0}, &graph.Node{Value: BB1})
	correctGraph.InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: BB2})
	correctGraph.InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: BB3})
	correctGraph.InsertEdge(&graph.Node{Value: BB2}, &graph.Node{Value: BB4})
	correctGraph.InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: BB4})
	correctGraph.InsertEdge(&graph.Node{Value: BB4}, &graph.Node{Value: EXIT})
	correctGraph.InsertEdge(&graph.Node{Value: EXIT}, &graph.Node{Value: START})

	if err := VerifyControlFlowGraphs(expectedGraph[0], correctGraph); err != nil {
		t.Fatal(err)
	}
}

func TestForLoopControlFlowGraph(t *testing.T) {
	sourceFile, err := ioutil.ReadFile("./testcode/_looper.go")
	if err != nil {
		t.Fatal(err)
	}
	basicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedGraph := cfgraph.GetControlFlowGraph(basicBlocks)
	correctGraph := graph.NewGraph()

	START := bblock.NewBasicBlock(-1, bblock.START, 0)
	EXIT := bblock.NewBasicBlock(-1, bblock.EXIT, 0)

	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.FOR_STATEMENT, 11)
	BB2 := bblock.NewBasicBlock(2, bblock.FOR_BODY, 14)
	BB3 := bblock.NewBasicBlock(3, bblock.RETURN_STMT, 16)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3)
	BB2.AddSuccessorBlock(BB1)

	correctBasicBlocks := []*bblock.BasicBlock{BB0, BB1, BB2, BB3}

	//Test basic-blocks.
	if err := VerifyBasicBlocks(basicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}

	correctGraph.InsertEdge(&graph.Node{Value: START}, &graph.Node{Value: BB0})
	correctGraph.InsertEdge(&graph.Node{Value: BB0}, &graph.Node{Value: BB1})
	correctGraph.InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: BB2})
	correctGraph.InsertEdge(&graph.Node{Value: BB2}, &graph.Node{Value: BB1})
	correctGraph.InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: BB3})
	correctGraph.InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: EXIT})
	correctGraph.InsertEdge(&graph.Node{Value: EXIT}, &graph.Node{Value: START})

	//Test control-flow-graph.
	if err := VerifyControlFlowGraphs(expectedGraph[0], correctGraph); err != nil {
		t.Fatal(err)
	}
}

func TestSwitchControlFlowGraph(t *testing.T) {
	sourceFile, err := ioutil.ReadFile("./testcode/_switcher.go")
	if err != nil {
		t.Fatal(err)
	}
	basicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedGraphs := cfgraph.GetControlFlowGraph(basicBlocks)
	correctGraph := []*graph.Graph{
		graph.NewGraph(), // func 'main'
		graph.NewGraph(), // func 'integerToString'
	}

	START := bblock.NewBasicBlock(-1, bblock.START, 0)
	EXIT := bblock.NewBasicBlock(-1, bblock.EXIT, 0)

	// Function 'main'
	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.RETURN_STMT, 11)

	// Function 'integerToString'
	BB2 := bblock.NewBasicBlock(2, bblock.FUNCTION_ENTRY, 13)
	BB3 := bblock.NewBasicBlock(3, bblock.SWITCH_STATEMENT, 14)
	BB4 := bblock.NewBasicBlock(4, bblock.RETURN_STMT, 16)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 18)
	BB6 := bblock.NewBasicBlock(6, bblock.RETURN_STMT, 20)
	BB7 := bblock.NewBasicBlock(7, bblock.RETURN_STMT, 22)
	BB8 := bblock.NewBasicBlock(8, bblock.RETURN_STMT, 24)

	BB0.AddSuccessorBlock(BB1)
	BB2.AddSuccessorBlock(BB3)
	BB3.AddSuccessorBlock(BB4, BB5, BB6, BB7, BB8)

	correctBasicBlocks := []*bblock.BasicBlock{BB0, BB1, BB2, BB3, BB4, BB5, BB6, BB7, BB8}

	//Test basic-blocks.
	if err := VerifyBasicBlocks(basicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}

	// Control flow graph for function 'main'.
	correctGraph[0].InsertEdge(&graph.Node{Value: START}, &graph.Node{Value: BB0})
	correctGraph[0].InsertEdge(&graph.Node{Value: BB0}, &graph.Node{Value: BB1})
	correctGraph[0].InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: EXIT})
	correctGraph[0].InsertEdge(&graph.Node{Value: EXIT}, &graph.Node{Value: START})

	// Control flow graph for function 'integerToString'.
	correctGraph[1].InsertEdge(&graph.Node{Value: START}, &graph.Node{Value: BB2})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB2}, &graph.Node{Value: BB3})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: BB4})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: BB5})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: BB6})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: BB7})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: BB8})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB8}, &graph.Node{Value: EXIT})
	correctGraph[1].InsertEdge(&graph.Node{Value: EXIT}, &graph.Node{Value: START})

	if err := VerifyControlFlowGraphs(expectedGraphs[0], correctGraph[0]); err != nil {
		t.Fatal(err)
	}
	if err := VerifyControlFlowGraphs(expectedGraphs[1], correctGraph[1]); err != nil {
		t.Fatal(err)
	}
}

func TestGreatestCommonDivisorControlFlowGraph(t *testing.T) {
	sourceFile, err := ioutil.ReadFile("./testcode/_gcd.go")
	if err != nil {
		t.Fatal(err)
	}
	basicBlocks, err := bblock.GetBasicBlocksFromSourceCode(sourceFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedGraphs := cfgraph.GetControlFlowGraph(basicBlocks)
	correctGraph := []*graph.Graph{
		graph.NewGraph(), // func 'gcd'
		graph.NewGraph(), // func 'main'
	}

	// Function 'gcd'
	START0 := bblock.NewBasicBlock(-1, bblock.START, 0)
	EXIT0 := bblock.NewBasicBlock(-1, bblock.EXIT, 0)
	BB0 := bblock.NewBasicBlock(0, bblock.FUNCTION_ENTRY, 8)
	BB1 := bblock.NewBasicBlock(1, bblock.FOR_STATEMENT, 9)
	BB2 := bblock.NewBasicBlock(2, bblock.FOR_BODY, 11)
	BB3 := bblock.NewBasicBlock(3, bblock.RETURN_STMT, 12)

	// Function 'main'
	START1 := bblock.NewBasicBlock(-1, bblock.START, 0)
	EXIT1 := bblock.NewBasicBlock(-1, bblock.EXIT, 0)
	BB4 := bblock.NewBasicBlock(4, bblock.FUNCTION_ENTRY, 15)
	BB5 := bblock.NewBasicBlock(5, bblock.RETURN_STMT, 18)

	BB0.AddSuccessorBlock(BB1)
	BB1.AddSuccessorBlock(BB2, BB3)
	BB2.AddSuccessorBlock(BB1)
	BB4.AddSuccessorBlock(BB5)

	correctBasicBlocks := []*bblock.BasicBlock{BB0, BB1, BB2, BB3, BB4, BB5}

	// Function 'gcd'
	correctGraph[0].InsertEdge(&graph.Node{Value: START0}, &graph.Node{Value: BB0})
	correctGraph[0].InsertEdge(&graph.Node{Value: BB0}, &graph.Node{Value: BB1})
	correctGraph[0].InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: BB2})
	correctGraph[0].InsertEdge(&graph.Node{Value: BB2}, &graph.Node{Value: BB1})
	correctGraph[0].InsertEdge(&graph.Node{Value: BB1}, &graph.Node{Value: BB3})
	correctGraph[0].InsertEdge(&graph.Node{Value: BB3}, &graph.Node{Value: EXIT0})
	correctGraph[0].InsertEdge(&graph.Node{Value: EXIT0}, &graph.Node{Value: START0})

	// Function 'main'
	correctGraph[1].InsertEdge(&graph.Node{Value: START1}, &graph.Node{Value: BB4})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB4}, &graph.Node{Value: BB5})
	correctGraph[1].InsertEdge(&graph.Node{Value: BB5}, &graph.Node{Value: EXIT1})
	correctGraph[1].InsertEdge(&graph.Node{Value: EXIT1}, &graph.Node{Value: START1})

	// Test basic-blocks.
	if err := VerifyBasicBlocks(basicBlocks, correctBasicBlocks); err != nil {
		t.Fatal(err)
	}
	// Test control-flow graph.
	if err := VerifyControlFlowGraphs(expectedGraphs[0], correctGraph[0]); err != nil {
		t.Fatal(err)
	}
	if err := VerifyControlFlowGraphs(expectedGraphs[1], correctGraph[1]); err != nil {
		t.Fatal(err)
	}
}
