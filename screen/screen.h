#ifndef tisscreen_screen_h
#define tisscreen_screen_h

bool init_screen();
void free_screen();
void print_text_screen(char* text);
void clear_screen();
void present_screen();

#endif