//    Copyright (C) 2021 Alexandre Viau <alexandre@alexandreviau.net>
//
//    This file is part of gopass.
//
//    gopass is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    gopass is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with gopass.  If not, see <http://www.gnu.org/licenses/>.

package clitest

import "time"

// runOptions contain *optitonal* parameters for Run().
type runOptions struct {
	editFunc func(string) (string, error)
	nowFunc  func() time.Time
}

type RunOption func(*runOptions)

func WithEditFunc(editFunc func(string) (string, error)) RunOption {
	return func(opts *runOptions) {
		opts.editFunc = editFunc
	}
}

func WithNowFunc(nowFunc func() time.Time) RunOption {
	return func(opts *runOptions) {
		opts.nowFunc = nowFunc
	}
}

func WithFixedTime(t time.Time) RunOption {
	return WithNowFunc(
		func() time.Time {
			return t
		},
	)
}
