//kage:unit pixels
package main

var Resolution vec2
var Time float

// 3D Gradient noise from: https://www.shadertoy.com/view/Xsl3Dl
func hash(p vec3) vec3 { // replace this by something better
	p = vec3(dot(p, vec3(127.1, 311.7, 74.7)),
		dot(p, vec3(269.5, 183.3, 246.1)),
		dot(p, vec3(113.5, 271.9, 124.6)))

	return -1.0 + 2.0*fract(sin(p)*43758.5453123)
}

func noise(p vec3) float {
	i := floor(p)
	f := fract(p)

	u := f * f * (3.0 - 2.0*f)

	return mix(mix(mix(dot(hash(i+vec3(0.0, 0.0, 0.0)), f-vec3(0.0, 0.0, 0.0)),
		dot(hash(i+vec3(1.0, 0.0, 0.0)), f-vec3(1.0, 0.0, 0.0)), u.x),
		mix(dot(hash(i+vec3(0.0, 1.0, 0.0)), f-vec3(0.0, 1.0, 0.0)),
			dot(hash(i+vec3(1.0, 1.0, 0.0)), f-vec3(1.0, 1.0, 0.0)), u.x), u.y),
		mix(mix(dot(hash(i+vec3(0.0, 0.0, 1.0)), f-vec3(0.0, 0.0, 1.0)),
			dot(hash(i+vec3(1.0, 0.0, 1.0)), f-vec3(1.0, 0.0, 1.0)), u.x),
			mix(dot(hash(i+vec3(0.0, 1.0, 1.0)), f-vec3(0.0, 1.0, 1.0)),
				dot(hash(i+vec3(1.0, 1.0, 1.0)), f-vec3(1.0, 1.0, 1.0)), u.x), u.y), u.z)
}

// from Unity's black body Shader Graph node
func UnityBlackbodyFloat(Temperature float) vec3 {
	color := vec3(255.0, 255.0, 255.0)
	color.x = 56100000.*pow(Temperature, (-3.0/2.0)) + 148.0
	color.y = 100.04*log(Temperature) - 623.6
	if Temperature > 6500.0 {
		color.y = 35200000.0*pow(Temperature, (-3.0/2.0)) + 184.0
	}
	color.z = 194.18*log(Temperature) - 1448.6
	color = clamp(color, 0.0, 255.0) / 255.0
	if Temperature < 1000.0 {
		color *= Temperature / 1000.0
	}
	return color
}

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	// Normalized pixel coordinates (from 0 to 1)
	// uv := (2.0 * dstPos.xy - Resolution.xy)/Resolution.y
	uv := dstPos.xy / Resolution.xy

	// Stars computation:
	stars_direction := normalize(vec3(uv*2.0-1.0, 1.0)) // could be view vector for example
	stars_threshold := 8.0                              // modifies the number of stars that are visible
	stars_exposure := 200.0                             // modifies the overall strength of the stars
	stars := pow(clamp(noise(stars_direction*200.0), 0.0, 1.0), stars_threshold) * stars_exposure
	stars *= mix(0.4, 1.4, noise(stars_direction*100.0+vec3(Time))) // time based flickering

	// star color by randomized temperature
	stars_temperature := noise(stars_direction*150.0)*0.5 + 0.5
	stars_color := UnityBlackbodyFloat(mix(1500.0, 65000.0, pow(stars_temperature, 4.0)))

	return vec4(stars_color*stars, 1.0)
}
