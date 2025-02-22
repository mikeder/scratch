// Copyright 2020 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ignore

//kage:unit pixels

package main

// var Resolution vec2
var Time float

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	uv := dstPos.xy - imageDstOrigin()

	border := imageDstSize().y*0.65 + 4*cos(Time*3+uv.y/10)
	if uv.y < border {
		return imageSrc0UnsafeAt(srcPos)
	}

	xoffset := 4 * cos(Time*3+uv.y/10)
	yoffset := 5 * (1 + cos(Time*3+uv.y/20))
	srcOrigin := imageSrc0Origin()
	clr := imageSrc0At(vec2(
		srcPos.x+xoffset,
		-(srcPos.y+yoffset-srcOrigin.y)+border*2+srcOrigin.y,
	)).rgb

	overlay := vec3(0, 0, 0.3)
	return vec4(mix(clr, overlay, 0.25), 1)
}
