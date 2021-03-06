/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	var rp *retryProcess

	mockChannel := newChannel("foo.channel", defaultPartition)
	flag := false

	noErrorFn := func() error {
		flag = true
		return nil
	}

	errorFn := func() error { return fmt.Errorf("foo") }

	t.Run("Exit", func(t *testing.T) {
		exitChan := make(chan struct{})
		close(exitChan)
		rp = newRetryProcess(mockRetryOptions, exitChan, mockChannel, "foo", noErrorFn)
		assert.Error(t, rp.retry(), "Expected retry to return an error")
		assert.Equal(t, false, flag, "Expected flag to remain set to false")
	})

	t.Run("Proper", func(t *testing.T) {
		exitChan := make(chan struct{})
		rp = newRetryProcess(mockRetryOptions, exitChan, mockChannel, "foo", noErrorFn)
		assert.NoError(t, rp.retry(), "Expected retry to return no errors")
		assert.Equal(t, true, flag, "Expected flag to be set to true")
	})

	t.Run("WithError", func(t *testing.T) {
		exitChan := make(chan struct{})
		rp = newRetryProcess(mockRetryOptions, exitChan, mockChannel, "foo", errorFn)
		assert.Error(t, rp.retry(), "Expected retry to return an error")
	})
}
