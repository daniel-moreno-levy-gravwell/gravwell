/*************************************************************************
 * Copyright 2025 Gravwell, Inc. All rights reserved.
 * Contact: <legal@gravwell.io>
 *
 * This software may be modified and distributed under the terms of the
 * BSD 2-clause license. See the LICENSE file for details.
 **************************************************************************/

// Package debugout use to print out simple debug messages
package debugout

type DebugOutFunc func(string, ...interface{})

var DebugOut DebugOutFunc = noop

func noop(format string, args ...interface{}) {}
