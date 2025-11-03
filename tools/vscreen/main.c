#include "main.h"

static void my_dpms_handler(int dpms_mode, void *user_data)
{
    go_dpms_handler(dpms_mode, user_data);
}

static void my_mode_changed_handler(struct evdi_mode mode, void *user_data)
{
    go_mode_changed_handler(mode, user_data);
}

static void my_update_ready_handler(int buffer_to_be_updated, void *user_data)
{
    go_update_ready_handler(buffer_to_be_updated, user_data);
}

static void my_crtc_state_handler(int buffer_to_be_updated, void *user_data)
{
    go_crtc_state_handler(buffer_to_be_updated, user_data);
}

static void my_cursor_set_handler(struct evdi_cursor_set cursor_set, void *user_data)
{
    go_cursor_set_handler(cursor_set, user_data);
}

static void my_cursor_move_handler(struct evdi_cursor_move cursor_move, void *user_data)
{
    go_cursor_move_handler(cursor_move, user_data);
}

static void my_ddcci_data_handler(struct evdi_ddcci_data ddcci_data, void *user_data)
{
    go_ddcci_data_handler(ddcci_data, user_data);
}

static void my_evdi_log(void *user_data, const char *fmt, ...)
{

    char buffer[4096];

    va_list args;
    va_start(args, fmt);
    vsnprintf(buffer, sizeof(buffer), fmt, args);
    va_end(args);

    go_log(user_data, buffer);
}

void set_evdi_log(void *user_data)
{
    struct evdi_logging l = {
        my_evdi_log,
        user_data};
    evdi_set_logging(l);
}

struct evdi_event_context my_events(void *user_data)
{
    struct evdi_event_context e = {
        &my_dpms_handler,
        &my_mode_changed_handler,
        &my_update_ready_handler,
        &my_crtc_state_handler,
        &my_cursor_set_handler,
        &my_cursor_move_handler,
        &my_ddcci_data_handler,
        user_data};
    return e;
}