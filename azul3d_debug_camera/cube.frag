#version 120

varying vec3 position;

void main()
{
	gl_FragColor = vec4(abs(position)/5, 1.0);
}
