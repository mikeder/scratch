//kage:unit pixels
package main

var Resolution vec2
var Time float

func Fragment(targetCoords vec4, _ vec2, _ vec4) vec4 {
	uv := (2.0*targetCoords.xy - Resolution.xy) / Resolution.y
	battery := 1.0

	sunUV := -uv
	col := vec3(1.0, 0.2, 1.0)
	sunVal := sun(sunUV, battery)

	col = mix(col, vec3(1.0, 0.4, 0.1), sunUV.y*2.0+0.2)
	// col = mix(col, vec3(1.0, 0.4, 0.1), sunUV.y * 2.0 + 0.2) // amplify colors
	col = mix(vec3(0.0, 0.0, 0.0), col, sunVal)

	return vec4(col, 0.85)
}

// adapted from https://www.shadertoy.com/view/Wt33Wf
func sun(uv vec2, battery float) float {
	val := smoothstep(0.3, 0.29, length(uv))
	bloom := smoothstep(0.9, 0.0, length(uv)*0.55)
	cut := 3.0*sin((uv.y+Time*0.2*(battery+0.02))*100.0) + clamp(uv.y*14.0+1.0, -6.0, 6.0)
	cut = clamp(cut, 0.0, 1.0)
	return clamp(val*cut, 0.0, 1.0) + bloom*0.6
}
