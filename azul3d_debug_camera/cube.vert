#version 120

attribute vec3 Vertex;

uniform mat4 MVP;

varying vec3 position;

void main()
{
	position = Vertex;
	gl_Position = MVP * vec4(Vertex, 1.0);
}
