/*
Copyright (c) 2016-2023 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#include <assert.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>
#include <stdint.h>

#include <libretro.h>

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

#define AUDIO_FREQUENCY 44100

#define NL "\n"
#define LOG(...) log_cb(RETRO_LOG_INFO, __VA_ARGS__)

lua_State *L = NULL;

retro_log_printf_t log_cb = NULL;
retro_video_refresh_t video_cb = NULL;
retro_audio_sample_t audio_cb = NULL;
retro_input_poll_t input_poll_cb = NULL;
retro_input_state_t input_state_cb = NULL;
retro_environment_t environ_cb = NULL;

static void no_log(enum retro_log_level level, const char *fmt, ...) {
    va_list args;
    va_start(args, fmt);
    printf("[%d] ", (int)level);
    vprintf(fmt, args);
    va_end(args);
}

static void *lua_realloc(void *ud, void *ptr, size_t osize, size_t nsize) {
    (void)ud; (void)osize;
    if (!nsize) {
        if (ptr)
            free(ptr);
        return NULL;
    } else if (ptr) {
        return realloc(ptr, nsize);
    }
    return malloc(nsize);
}

static void audio_callback(void) {
    //audio_cb(sample, sample);
}

static void keyboard_event(bool down, unsigned keycode, uint32_t character, uint16_t key_modifiers) {
    (void)character; (void)key_modifiers; (void)down; (void)keycode;
}

static void mouse_event(void) {
    /*
    int x = input_state_cb(0, RETRO_DEVICE_MOUSE, 0, RETRO_DEVICE_ID_MOUSE_X);
    int y = input_state_cb(0, RETRO_DEVICE_MOUSE, 0, RETRO_DEVICE_ID_MOUSE_Y);
    int l = input_state_cb(0, RETRO_DEVICE_MOUSE, 0, RETRO_DEVICE_ID_MOUSE_LEFT);
    int r = input_state_cb(0, RETRO_DEVICE_MOUSE, 0, RETRO_DEVICE_ID_MOUSE_RIGHT);
    */
}

static void check_variables(void) {
}

void retro_init(void) {
    check_variables();

}

void retro_deinit(void) {

}

unsigned retro_api_version(void) {
    return RETRO_API_VERSION;
}

void retro_set_controller_port_device(unsigned port, unsigned device) {
    (void)port; (void)device;
}

void retro_get_system_info(struct retro_system_info *info) {
    memset(info, 0, sizeof(struct retro_system_info));
    info->library_name = "OpenWar";
    info->library_version = "0.0.1";
    info->valid_extensions = "war";
}

void retro_get_system_av_info(struct retro_system_av_info *info) {
    memset(info, 0, sizeof(struct retro_system_av_info));
    info->geometry.base_width = 320;
    info->geometry.base_height = 200;
    info->geometry.max_width = 320;
    info->geometry.max_height = 200;
    info->geometry.aspect_ratio = 4.0f / 3.0f;
    info->timing.sample_rate = (double)AUDIO_FREQUENCY;
    info->timing.fps = 60.0;
}

void retro_set_environment(retro_environment_t cb) {
    environ_cb = cb;
    struct retro_log_callback logging;
    log_cb = cb(RETRO_ENVIRONMENT_GET_LOG_INTERFACE, &logging) ? logging.log : no_log;
}

void retro_set_audio_sample(retro_audio_sample_t cb) {
    audio_cb = cb;
}

void retro_set_audio_sample_batch(retro_audio_sample_batch_t cb) {
    (void)cb;
}

void retro_set_input_poll(retro_input_poll_t cb) {
    input_poll_cb = cb;
}

void retro_set_input_state(retro_input_state_t cb) {
    input_state_cb = cb;
}

void retro_set_video_refresh(retro_video_refresh_t cb) {
    video_cb = cb;
}

void retro_reset(void) {
}

void retro_run(void) {
    bool updated = false;
    if (environ_cb(RETRO_ENVIRONMENT_GET_VARIABLE_UPDATE, &updated) && updated)
        check_variables();

    input_poll_cb();

    if (lua_getglobal(L, "update") == LUA_TFUNCTION)
        lua_pcall(L, 0, 0, -1);
    else
        lua_pop(L, 1);

    mouse_event();
}

bool retro_load_game(const struct retro_game_info *info) {
    check_variables();

    struct retro_keyboard_callback kbdesk = { &keyboard_event };
    environ_cb(RETRO_ENVIRONMENT_SET_KEYBOARD_CALLBACK, &kbdesk);

    enum retro_pixel_format fmt = RETRO_PIXEL_FORMAT_XRGB8888;
    if (!environ_cb(RETRO_ENVIRONMENT_SET_PIXEL_FORMAT, &fmt)) {
        log_cb(RETRO_LOG_ERROR, "XRGB8888 is not supported!\n");
        return false;
    }

    struct retro_audio_callback audio_cb = { &audio_callback, NULL };
    environ_cb(RETRO_ENVIRONMENT_SET_AUDIO_CALLBACK, &audio_cb);

    //if (!load_resources((unsigned char*)info->data, (int)info->size))
    //    return false;

    if (!(L = lua_newstate(&lua_realloc, NULL)))
        return false;

    luaL_openlibs(L);
    if (luaL_loadfile(L, "script/main.lua"))
        return false;

    return true;
}

void retro_unload_game(void) {
    if (lua_getglobal(L, "shutdown") == LUA_TFUNCTION)
        lua_pcall(L, 0, 0, -1);
    else
        lua_pop(L, 1);
    lua_close(L);
    L = NULL;
}

unsigned retro_get_region(void) {
    return RETRO_REGION_NTSC;
}

bool retro_load_game_special(unsigned type, const struct retro_game_info *info, size_t num) {
    (void)type; (void)info; (void)num;
    return false;
}

size_t retro_serialize_size(void) {
    return 0;
}

bool retro_serialize(void *data, size_t size) {
    (void)data; (void)size;
    return false;
}

bool retro_unserialize(const void *data, size_t size) {
    (void)data; (void)size;
    return false;
}

void *retro_get_memory_data(unsigned id) {
    (void)id;
    return NULL;
}

size_t retro_get_memory_size(unsigned id) {
    (void)id;
    return 0;
}

void retro_cheat_reset(void) {
}

void retro_cheat_set(unsigned index, bool enabled, const char *code) {
    (void)index; (void)enabled; (void)code;
}
