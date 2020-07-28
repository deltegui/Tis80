#include <SDL2/SDL.h>
#include <SDL2/SDL_ttf.h>
#include <stdbool.h>
#include <string.h>
#include "screen.h"

#define SCREEN_WIDTH 320
#define SCREEN_HEIGHT 200

typedef struct {
	SDL_Window* window;
	SDL_Renderer* renderer;
	TTF_Font* font;
} Screen;

Screen screen;

bool init_screen() {
    if(SDL_Init(SDL_INIT_VIDEO) != 0) {
        printf("SDL could not initialize! SDL_Error: %s\n", SDL_GetError());
        return false;
    }
    
    screen.window = SDL_CreateWindow("Tis80", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, SCREEN_WIDTH, SCREEN_HEIGHT, SDL_WINDOW_SHOWN | SDL_WINDOW_RESIZABLE);
    if(screen.window == NULL){
        printf("Window could not be created! SDL_Error: %s\n", SDL_GetError());
        return false;
    }
    
    screen.renderer = SDL_CreateRenderer(screen.window, -1, SDL_RENDERER_ACCELERATED | SDL_RENDERER_PRESENTVSYNC);
    if(screen.renderer == NULL) {
    	printf("Renderer could not be created! SDL_Error: %s\n", SDL_GetError());
    	SDL_DestroyWindow(screen.window);
    	return false;
    }

    if(TTF_Init() != 0) {
    	printf("TTF system cannot be loaded! TTF_Error: %s\n", TTF_GetError());
    	SDL_DestroyWindow(screen.window);
    	SDL_DestroyRenderer(screen.renderer);
    	return false;
    }

    screen.font = TTF_OpenFont("./tis80.ttf", 25);
    if(screen.font == NULL) {
    	printf("Cannot open font: %s\n", TTF_GetError());
    	SDL_DestroyWindow(screen.window);
    	SDL_DestroyRenderer(screen.renderer);
    	return false;
    }

    return true;
}

void free_screen() {
    TTF_CloseFont(screen.font);
    SDL_DestroyRenderer(screen.renderer);
    SDL_DestroyWindow(screen.window);
    TTF_Quit();
	SDL_Quit();
}

static void print_letter(char letter) {
    static int letter_x = 0;
    static int letter_y = 0;
    const char str[] = {letter, '\0'}; 

    SDL_Color White = {255, 255, 255};
    SDL_Surface* surfaceMessage = TTF_RenderText_Solid(screen.font, str, White);
    SDL_Texture* Message = SDL_CreateTextureFromSurface(screen.renderer, surfaceMessage);
    
    SDL_Rect Message_rect;
    Message_rect.x = letter_x * 8;
    Message_rect.y = letter_y * 8;
    Message_rect.w = 8; // controls the width of the rect
    Message_rect.h = 8; // controls the height of the rect

    SDL_RenderCopy(screen.renderer, Message, NULL, &Message_rect);
    SDL_FreeSurface(surfaceMessage);
    SDL_DestroyTexture(Message);

    letter_x++;
    if(letter_x >= 40) {
        letter_x = 0;
        letter_y++;
    }
    if(letter_y >= 25) {
        letter_y = 0;
    }
}

void print_text_screen(char* text) {
    printf("Print text: %s\n", text);
    int len = strlen(text);
    for(int i = 0; i < len; i++) {
        print_letter(text[i]);
    }
}

void clear_screen() {
    SDL_RenderClear(screen.renderer);
}

void present_screen() {
    SDL_RenderPresent(screen.renderer);
}