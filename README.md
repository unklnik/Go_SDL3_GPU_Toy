# Go SDL3 GPU Toy
A small, unfinished GUI for working with SDL3 GPU in 3D and 2D using https://github.com/Zyko0/go-sdl3 SDL3 bindings for Go. The idea was to learn a bit of SDL3 for use with 2D as well as 3D. Note that you will need to also have Shadercross https://github.com/libsdl-org/SDL_shadercross installed. The files are included here so you should be able to compile and run on Windows, on Linux you will need to add Shadercross (Shader translation library for SDL's GPU API). Note this is not a complete project, just a small project to learn a bit more about SDL3 GPU and put it up in case it is useful to someone else.

**Features:**
- Load .glb models
- Draw 2D geometry (lines/rectangles/triangles/grids)
- Draw 2D text
- Draw 3D primitives (cubes/spheres/cones)
- Load 2D textures & create animations

**Keys**<br>
F1 key > Settings<br>
F2 key > Show/Hide UI<br>
F10 key > Debug overlay<br>
F11 key > Color Palette<br>
F12 key > Color Palette with Alpha<br>
ESC key > Exit<br>

*AI Disclosure: Setting up SDL GPU to run with Go is not that easy and I didn't write that code at all, that was made using Google AI Studio and took a while with a lot of errors along the way. However, the engine itself was mainly coded by myself.*

https://github.com/user-attachments/assets/de19581b-6a04-40de-8a5d-be561c1b4121
