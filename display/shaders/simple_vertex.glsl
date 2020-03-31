#version 330
//uniform mat4 projection;
//uniform mat4 camera;
//uniform mat4 model;

in vec3 position;

void main() {
//    gl_Position = projection * camera * model * vec4(vert, 1);
    gl_Position = vec4(position, 1);
}
