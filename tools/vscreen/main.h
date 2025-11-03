#ifndef MAIN_H
#define MAIN_H

#include "library/evdi_lib.h"
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>

extern void go_dpms_handler(int dpms_mode, void *user_data);
extern void go_mode_changed_handler(struct evdi_mode mode, void *user_data);
extern void go_update_ready_handler(int buffer_to_be_updated, void *user_data);
extern void go_crtc_state_handler(int buffer_to_be_updated, void *user_data);
extern void go_cursor_set_handler(struct evdi_cursor_set cursor_set, void *user_data);
extern void go_cursor_move_handler(struct evdi_cursor_move cursor_move, void *user_data);
extern void go_ddcci_data_handler(struct evdi_ddcci_data ddcci_data, void *user_data);
extern void go_log(void *user_data, char *msg);

struct evdi_event_context my_events(void *user_data);
void set_evdi_log(void *user_data);

#endif