#include <MagickWand/MagickWand.h>
#include <stdio.h>

int main() {
    MagickWandGenesis();
    printf("ImageMagick Wand initialized successfully.\n");
    MagickWandTerminus();
    return 0;
}
