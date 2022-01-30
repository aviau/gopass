//    Copyright (C) 2022 Alexandre Viau <alexandre@alexandreviau.net>
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

// Package alfred defines structs for Alfred script filters.
package alfred

type Output struct {
	Items []*Item `json:"items"`
}

type Item struct {
	UID      string `json:"uid"`
	Title    string `json:"title"`
	Arg      string `json:"arg"`
	Subtitle string `json:"subtitle,omitempty"`
}
