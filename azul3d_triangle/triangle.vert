#version 120

attribute vec3 Vertex;
attribute vec4 Color;

uniform mat4 MVP;

varying vec4 frontColor;

void main()
{
	frontColor = Color;
	gl_PointSize = 25.0;
	gl_Position = MVP * vec4(Vertex, 1.0);
}
