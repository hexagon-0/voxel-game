# Voxel Game

This project is a demo of basic Minecraft-style voxel mechanics. It is written
in Go using OpenGL bindings from [go-gl](https://github.com/go-gl). Currently,
it is able to render a single chunk mesh using a culling method (hidden faces
aren't produced). The main reference for this was [0fps.net](https://0fps.net/2012/06/30/meshing-in-a-minecraft-game/)
and the [corresponding repo](https://github.com/mikolalysenko/mikolalysenko.github.com/blob/master/MinecraftMeshes2/js/culled.js).

I will continue to work in this project in the future, but most likely ported
to a different language.

## Building

To build this project, you will need the Go toolchain as well as a C compiler
(I've used GCC):

```
go run ./cmd/client/main.go
```

Or use the `build` command to generate an executable. Regardless of the method,
you'll need to run this from the same directory where the `assets` folder is
located.
