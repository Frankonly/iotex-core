// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package ioctlio

import (
	"errors"
	"fmt"
	"log"

	"github.com/iotexproject/iotex-core/ioctl/util"
)

var (
	NoNewInput  = errors.New("no new input in buffer")
	NoNewOutput = errors.New("no new output in buffer")
)

type IoController struct {
	isTest       bool
	inputBuffer  []string
	outputBuffer []string
	lastInput    int
	lastOutput   int
}

// IoCtrl is the controller of ioctl' io to switch between user's call and test
var IoCtrl IoController

// ReadSecret used to get secret input
func ReadSecret() (string, error) {
	if IoCtrl.isTest {
		return IoCtrl.readInput(), nil
	}

	return util.ReadSecretFromStdin()
}

// Scanf used to invoke fmt.Scanf() or fmt.Sscanf()
func Scanf(format string, a ...interface{}) (err error) {
	if IoCtrl.isTest {
		input := IoCtrl.readInput()
		_, err = fmt.Sscanf(input, format, a...)
	} else {
		_, err = fmt.Scanf(format, a...)
	}

	return
}

// SetStatus sets whether program is test status
func (ctrl *IoController) SetStatus(isTest bool) {
	ctrl.isTest = isTest
}

// IsTest returns whether program is test status
func (ctrl *IoController) IsTest() bool {
	return ctrl.isTest
}

// OutputBuffer returns IoController's output buffer
func (ctrl *IoController) OutputBuffer() []string {
	return ctrl.outputBuffer
}

// ReadOutput returns latest outputs from IoController's output buffer
// with a counter to check whether buffer has new output
func (ctrl *IoController) ReadOutput() ([]string, error) {
	if len(ctrl.outputBuffer) <= ctrl.lastOutput {
		log.Panic(NoNewOutput)
	}

	result := ctrl.outputBuffer[ctrl.lastOutput:]
	ctrl.lastOutput = len(ctrl.outputBuffer)

	return result, nil
}

// SetTestBuffer overwrites the buffer of IoController
func (ctrl *IoController) SetInputBuffer(inputBuffer []string) {
	ctrl.inputBuffer = inputBuffer
}

// AppendInputBuffer append buffer of IoController to the end of IoController's input buffer
func (ctrl *IoController) AppendInputBuffer(inputBuffer ...string) {
	ctrl.inputBuffer = append(ctrl.inputBuffer, inputBuffer...)
}

func (ctrl *IoController) readInput() string {
	if len(ctrl.inputBuffer) <= ctrl.lastInput {
		log.Panic(NoNewInput)
	}

	result := ctrl.inputBuffer[ctrl.lastInput]
	ctrl.lastInput++

	return result
}

func (ctrl *IoController) writeOutput(outputBuffer ...string) {
	ctrl.outputBuffer = append(ctrl.outputBuffer, outputBuffer...)
}

func init() {
	IoCtrl.SetStatus(false)
}
