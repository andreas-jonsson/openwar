workspace "openwar"
    configurations { "release", "debug" }
    platforms { "native" }
    language "C"
    cdialect "C11"
    pic "On"

    filter "configurations:debug"
        symbols "On"

    filter "configurations:release"
        defines "NDEBUG"
        optimize "On"

    filter "system:windows"
        defines "_CRT_SECURE_NO_WARNINGS"

    filter "toolset:clang or gcc"
        buildoptions { "-Wall", "-Wextra" }

    project "lua"
        kind "StaticLib"
        files { "lua/src/*.h", "lua/src/*.c" }
        removefiles { "lua/src/lua.c", "lua/src/luac.c" }

    project "core"
        kind "SharedLib"
        targetname "openwar_libretro"
        targetprefix ""
        targetdir ""
        pic "On"

        includedirs { "lua/src", "libretro" }
        links "lua"

        files { "src/*.h", "src/*.c" }
